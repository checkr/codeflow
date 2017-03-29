package mongodb_bongo

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/mattes/migrate/driver"
	"github.com/mattes/migrate/driver/mongodb/gomethods"
	"github.com/mattes/migrate/file"
	"github.com/mattes/migrate/migrate/direction"
	"github.com/maxwellhealth/bongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UnregisteredMethodsReceiverError string

func (e UnregisteredMethodsReceiverError) Error() string {
	return "Unregistered methods receiver for driver: " + string(e)
}

type WrongMethodsReceiverTypeError string

func (e WrongMethodsReceiverTypeError) Error() string {
	return "Wrong methods receiver type for driver: " + string(e)
}

const MIGRATE_C = "migrations"
const DRIVER_NAME = "gomethods.mongodb"

type Driver struct {
	Connection *bongo.Connection

	methodsReceiver MethodsReceiver
	migrator        gomethods.Migrator
}

var _ gomethods.GoMethodsDriver = (*Driver)(nil)

type MethodsReceiver interface {
	DbName() string
	SSL() bool
}

func (d *Driver) MethodsReceiver() interface{} {
	return d.methodsReceiver
}

func (d *Driver) SetMethodsReceiver(r interface{}) error {
	r1, ok := r.(MethodsReceiver)
	if !ok {
		return WrongMethodsReceiverTypeError(DRIVER_NAME)
	}

	d.methodsReceiver = r1
	return nil
}

func init() {
	driver.RegisterDriver("mongodb", &Driver{})
}

type DbMigration struct {
	Id      bson.ObjectId `bson:"_id"`
	Version uint64        `bson:"version"`
}

func (driver *Driver) Initialize(url string) error {
	var err error
	if driver.methodsReceiver == nil {
		return UnregisteredMethodsReceiverError(DRIVER_NAME)
	}

	urlWithoutScheme := strings.SplitN(url, "mongodb://", 2)
	if len(urlWithoutScheme) != 2 {
		return errors.New("invalid mongodb:// scheme")
	}

	config := &bongo.Config{
		ConnectionString: url,
		Database:         driver.methodsReceiver.DbName(),
	}

	if driver.methodsReceiver.SSL() {
		if config.DialInfo, err = mgo.ParseURL(config.ConnectionString); err != nil {
			panic(fmt.Sprintf("cannot parse given URI %s due to error: %s", config.ConnectionString, err.Error()))
		}

		tlsConfig := &tls.Config{}
		config.DialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			return conn, err
		}

		config.DialInfo.Timeout = time.Second * 3
	}

	connection, err := bongo.Connect(config)
	if err != nil {
		log.Fatal(err)
	}

	c := connection.Session.DB(driver.methodsReceiver.DbName()).C(MIGRATE_C)
	err = c.EnsureIndex(mgo.Index{
		Key:    []string{"version"},
		Unique: true,
	})
	if err != nil {
		return err
	}

	driver.Connection = connection
	driver.migrator = gomethods.Migrator{MethodInvoker: driver}

	return nil
}

func (driver *Driver) Close() error {
	if driver.Connection != nil {
		driver.Connection.Session.Close()
	}
	return nil
}

func (driver *Driver) FilenameExtension() string {
	return "bongo"
}

func (driver *Driver) Version() (uint64, error) {
	var latestMigration DbMigration
	c := driver.Connection.Session.DB(driver.methodsReceiver.DbName()).C(MIGRATE_C)

	err := c.Find(bson.M{}).Sort("-version").One(&latestMigration)

	switch {
	case err == mgo.ErrNotFound:
		return 0, nil
	case err != nil:
		return 0, err
	default:
		return latestMigration.Version, nil
	}
}
func (driver *Driver) Migrate(f file.File, pipe chan interface{}) {
	defer close(pipe)
	pipe <- f

	err := driver.migrator.Migrate(f, pipe)
	if err != nil {
		return
	}

	migrate_c := driver.Connection.Session.DB(driver.methodsReceiver.DbName()).C(MIGRATE_C)

	if f.Direction == direction.Up {
		id := bson.NewObjectId()
		dbMigration := DbMigration{Id: id, Version: f.Version}

		err := migrate_c.Insert(dbMigration)
		if err != nil {
			pipe <- err
			return
		}

	} else if f.Direction == direction.Down {
		err := migrate_c.Remove(bson.M{"version": f.Version})
		if err != nil {
			pipe <- err
			return
		}
	}
}

func (driver *Driver) Validate(methodName string) error {
	methodWithReceiver, ok := reflect.TypeOf(driver.methodsReceiver).MethodByName(methodName)
	if !ok {
		return gomethods.MethodNotFoundError(methodName)
	}
	if methodWithReceiver.PkgPath != "" {
		return gomethods.MethodNotFoundError(methodName)
	}

	methodFunc := reflect.ValueOf(driver.methodsReceiver).MethodByName(methodName)
	methodTemplate := func(*bongo.Connection) error { return nil }

	if methodFunc.Type() != reflect.TypeOf(methodTemplate) {
		return gomethods.WrongMethodSignatureError(methodName)
	}

	return nil
}

func (driver *Driver) Invoke(methodName string) error {
	name := methodName
	migrateMethod := reflect.ValueOf(driver.methodsReceiver).MethodByName(name)
	if !migrateMethod.IsValid() {
		return gomethods.MethodNotFoundError(methodName)
	}

	retValues := migrateMethod.Call([]reflect.Value{reflect.ValueOf(driver.Connection)})
	if len(retValues) != 1 {
		return gomethods.WrongMethodSignatureError(name)
	}

	if !retValues[0].IsNil() {
		err, ok := retValues[0].Interface().(error)
		if !ok {
			return gomethods.WrongMethodSignatureError(name)
		}
		return &gomethods.MethodInvocationFailedError{MethodName: name, Err: err}
	}

	return nil
}
