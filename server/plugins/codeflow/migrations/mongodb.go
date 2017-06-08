package codeflow_migrations

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/checkr/codeflow/server/plugins"
	"github.com/checkr/codeflow/server/plugins/codeflow"
	"github.com/checkr/codeflow/server/plugins/codeflow/migrations/driver"
	"github.com/mattes/migrate/driver/mongodb/gomethods"
	"github.com/maxwellhealth/bongo"
	"github.com/spf13/viper"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoDbMigrator struct {
}

func (r *MongoDbMigrator) DbName() string {
	return viper.GetString("plugins.codeflow.mongodb.database")
}

func (r *MongoDbMigrator) SSL() bool {
	return viper.GetBool("plugins.codeflow.mongodb.ssl")
}

var _ mongodb_bongo.MethodsReceiver = (*MongoDbMigrator)(nil)

func init() {
	gomethods.RegisterMethodsReceiverForDriver("mongodb", &MongoDbMigrator{})
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func (r *MongoDbMigrator) V001_init_users_up(c *bongo.Connection) error {
	var collection []codeflow.User
	var obj codeflow.User

	obj = codeflow.User{
		Name:     "Codeflow",
		Username: "codeflow@checkr.com",
		Email:    "codeflow",
	}
	obj.SetId(bson.ObjectIdHex("586bf835adfd558a1802772d"))

	collection = append(collection, obj)

	obj = codeflow.User{
		Name:     "Demo",
		Username: "demo@development.com",
		Email:    "demo@development.com",
	}
	obj.SetId(bson.ObjectIdHex("58dbe954df8ab3002a71dc07"))

	collection = append(collection, obj)

	for _, o := range collection {
		if err := c.Collection("users").Save(&o); err != nil {
			log.Printf("Save::Error: %v", err.Error())
			return err
		}
	}

	return nil
}

func (r *MongoDbMigrator) V001_init_users_down(c *bongo.Connection) error {
	return c.Session.DB(r.DbName()).C("users").DropCollection()
}

func (r *MongoDbMigrator) V001_init_projects_up(c *bongo.Connection) error {
	obj := codeflow.Project{
		Name:                  "checkr/codeflow",
		Slug:                  "checkr-codeflow",
		Repository:            "checkr/codeflow",
		Secret:                "9givnw17mqxd9m6um95pho9nkmy9st",
		Pinged:                true,
		GitUrl:                "https://github.com/checkr/codeflow.git",
		GitProtocol:           "HTTPS",
		RsaPrivateKey:         "-----BEGIN RSA PRIVATE KEY-----\nMIIEhgIBAAKB/DAojMGcTaZSDi9oaLLwYuQfFAQca/nK+d3X+EYhIKtEFszuVnk1\nXYXsD+Zutxt+N2awOyvnGjqNPZ6/4wSCvrf7PM+jIBreM+XICGl/0r0h2AldEyfq\ny7LD9KaDp91UOWrUK8bKEU0T0b9oASwbdNOKkXTs1bSbkpOCWSk6d8/V+K66vKnW\nh1H55esXzB+3xW7hdSzJC2f/0Xk2kjNddJAh1gH7mITbeD45CV761HvHGp8A1tSl\nPWRZz6qZrQ34zJJEX+iGMEQIE2GSFIRI33HXlEk5pRaQQ8LDhIAGESETqImZkyN7\nDYuFkoBXYXBj08JYGjwiaW7hM0dPGwIDAQABAoH8IMIEdNJAU2kvcvn/dfBkJC4r\nrFw06lYyPr/ghruUAEuxgraApbQyKJ2ZdzJKZW4mezhXF5b81WUrzCdUYcYZuwYv\nqEGa3gvVm3DEoBatn688x6nDFPz2kGQQr4+QiNH4uH0YRgE/YYGgxCUX3wvSHO79\n4F4VQ+QrASHCSnQVYhTcTuQd7XaKaRgyE2b8bv7Doyi/8P35ygAcSBfQWUUg8fRL\no1Te/YXyXtF91b9Lbvtjg1UpLYz4V9rbOO1iiR48x8vobpRKuxMH64VK5rwrbfns\nWAJwnjUCp5jS28hAfKoMzajxsanCz2rkgmfG3biESTNZWIun5UKRl7w5An53qlDX\nSz0lSRFcDMs2BQkG8gcLsUC8WL4rZigb78XgqVlwEgjU9OpHzGo6OLor8Q0YyAqi\nXaTXiLQzu4FD6P1N1JS9IoG+0TODvD1y+8rQbCSaV2AaH2A0qmtCnjCvDpNJ51RG\n4GYIi18PFwoMIYhT1c7BxYA2FNJuJw4ymR0CfmcGeE+011p9sTn9awmNcsCAZ1vL\nmurq6dxlEbBzuFaZm7mQ8B8wcGZXwzmkzYSuIaXESLDn376hh5sgIQwV4VXLmUF8\nWN5I+nsQjIhPVGmTSte+5l3ASNFxDxeU4ATLEK39YuZ1TV0eBZJgQd6o16ZqXo3f\nLSCGynro+gnLlwJ+TntPRcQ8uAVx8zMY27b1sq5tXIfF80EoiAIZ8CiTWML4u324\neSKfvLMeQE0QHN2dP1GDV/WetRUdSoiBQO6/opn3awwEmAdQh+efTZhB7evfHbKM\nftVxHVlfu3NQbp9aji+/oDRv9s6ha54qosYjSQiC76b+bXm+gSvwLdMpAn5Yvr0C\nQ8/B1kXEoyQBrYNsiO7/pqpCs6pRPAp5yaS/jEAVH+GHrE0WC4FSdUDHisvXI/ZN\n1N7qMfBC0vFEnNBm/CN+wmM2zvxc58t2W4dmDgfJQlrj5Q+UwmPytz4lQtqSVZNM\n2zyR+ptoFFyJNT3Vzwi2Asm3nARszaUcrO8CfltEtBuUxoMTaEj/Y1Nnu8oQMbgG\nuc+mNS6ikh1iTLwdppfz8rCtoTfojmAG8/AdtgSeHEAoVkbAN1wi7J+oIGp32Qn+\nHTOir9RbgoE2uFoenq6ABKP02JAnAJdDteNxJJI0sGUxbmNn4JjLTA7wiXagSOVB\no7Kec8tOPS09RA==\n-----END RSA PRIVATE KEY-----\n",
		RsaPublicKey:          "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAA/DAojMGcTaZSDi9oaLLwYuQfFAQca/nK+d3X+EYhIKtEFszuVnk1XYXsD+Zutxt+N2awOyvnGjqNPZ6/4wSCvrf7PM+jIBreM+XICGl/0r0h2AldEyfqy7LD9KaDp91UOWrUK8bKEU0T0b9oASwbdNOKkXTs1bSbkpOCWSk6d8/V+K66vKnWh1H55esXzB+3xW7hdSzJC2f/0Xk2kjNddJAh1gH7mITbeD45CV761HvHGp8A1tSlPWRZz6qZrQ34zJJEX+iGMEQIE2GSFIRI33HXlEk5pRaQQ8LDhIAGESETqImZkyN7DYuFkoBXYXBj08JYGjwiaW7hM0dPGw==\n",
		ContinuousIntegration: false,
		ContinuousDelivery:    false,
		Workflows:             []string{"build/DockerImage"},
		NotifyChannels:        "",
	}
	obj.SetId(bson.ObjectIdHex("58dbe995df8ab3002a71dc08"))

	if err := c.Collection("projects").Save(&obj); err != nil {
		log.Printf("Save::Error: %v", err.Error())
		return err
	}

	return nil
}

func (r *MongoDbMigrator) V001_init_projects_down(c *bongo.Connection) error {
	return c.Session.DB(r.DbName()).C("projects").DropCollection()
}

func (r *MongoDbMigrator) V001_init_builds_up(c *bongo.Connection) error {
	var collection []codeflow.Build
	var obj codeflow.Build

	obj = codeflow.Build{
		FeatureHash: "41fa95f534d8f57433ed4a48de1620cba35f2c2d",
		Type:        "DockerImage",
		State:       "complete",
		Image:       "docker.io/checkr/checkr-codeflow:latest",
		BuildLog:    "Step 1/21 ...",
		BuildError:  "",
	}
	obj.SetId(bson.ObjectIdHex("58dbeb6cdf8ab3002a71dc0b"))

	collection = append(collection, obj)

	obj = codeflow.Build{
		FeatureHash: "df600016edb26c48e1c999b28bb874257f65d037",
		Type:        "DockerImage",
		State:       "complete",
		Image:       "docker.io/checkr/checkr-codeflow:latest",
		BuildLog:    "Step 1/21 ...",
		BuildError:  "",
	}
	obj.SetId(bson.ObjectIdHex("58dbebacdf8ab3002a71dc0d"))

	collection = append(collection, obj)

	for _, o := range collection {
		if err := c.Collection("builds").Save(&o); err != nil {
			log.Printf("Save::Error: %v", err.Error())
			return err
		}
	}

	return nil
}

func (r *MongoDbMigrator) V001_init_builds_down(c *bongo.Connection) error {
	return c.Session.DB(r.DbName()).C("builds").DropCollection()
}

func (r *MongoDbMigrator) V001_init_bookmarks_up(c *bongo.Connection) error {
	obj := codeflow.Bookmark{
		ProjectId: bson.ObjectIdHex("58dbe995df8ab3002a71dc08"),
		UserId:    bson.ObjectIdHex("58dbe954df8ab3002a71dc07"),
	}
	obj.SetId(bson.ObjectIdHex("58dbe995df8ab3002a71dc09"))

	if err := c.Collection("bookmarks").Save(&obj); err != nil {
		log.Printf("Save::Error: %v", err.Error())
		return err
	}

	return nil
}

func (r *MongoDbMigrator) V001_init_bookmarks_down(c *bongo.Connection) error {
	return c.Session.DB(r.DbName()).C("bookmarks").DropCollection()
}

func (r *MongoDbMigrator) V001_init_extensions_up(c *bongo.Connection) error {
	var collection []codeflow.LoadBalancer
	var obj codeflow.LoadBalancer

	obj = codeflow.LoadBalancer{
		Name:      "codeflow-api",
		ProjectId: bson.ObjectIdHex("58dbe995df8ab3002a71dc08"),
		ServiceId: bson.ObjectIdHex("58dbecefdf8ab3002a71dc0e"),
		Extension: "LoadBalancer",
		DNS:       "unknown",
		Type:      "office",
		ListenerPairs: []codeflow.ListenerPair{
			{
				Source: codeflow.Listener{
					Port:     3001,
					Protocol: "",
				},
				Destination: codeflow.Listener{
					Port:     3001,
					Protocol: "TCP",
				},
			},
			{
				Source: codeflow.Listener{
					Port:     3002,
					Protocol: "",
				},
				Destination: codeflow.Listener{
					Port:     3002,
					Protocol: "TCP",
				},
			},
			{
				Source: codeflow.Listener{
					Port:     3003,
					Protocol: "",
				},
				Destination: codeflow.Listener{
					Port:     3003,
					Protocol: "TCP",
				},
			},
		},
		State:        "complete",
		StateMessage: "",
	}
	obj.SetId(bson.ObjectIdHex("58dbef04df8ab3002a71dc15"))
	collection = append(collection, obj)

	obj = codeflow.LoadBalancer{
		Name:      "codeflow-dashboard",
		ProjectId: bson.ObjectIdHex("58dbe995df8ab3002a71dc08"),
		ServiceId: bson.ObjectIdHex("58dbed10df8ab3002a71dc0f"),
		Extension: "LoadBalancer",
		DNS:       "unknown",
		Type:      "office",
		ListenerPairs: []codeflow.ListenerPair{
			{
				Source: codeflow.Listener{
					Port:     80,
					Protocol: "",
				},
				Destination: codeflow.Listener{
					Port:     9000,
					Protocol: "HTTP",
				},
			},
			{
				Source: codeflow.Listener{
					Port:     443,
					Protocol: "",
				},
				Destination: codeflow.Listener{
					Port:     9000,
					Protocol: "HTTPS",
				},
			},
		},
		State:        "complete",
		StateMessage: "",
	}
	obj.SetId(bson.ObjectIdHex("58dbef1edf8ab3002a71dc16"))
	collection = append(collection, obj)

	for _, o := range collection {
		if err := c.Collection("extensions").Save(&o); err != nil {
			log.Printf("Save::Error: %v", err.Error())
			return err
		}
	}

	return nil
}

func (r *MongoDbMigrator) V001_init_extensions_down(c *bongo.Connection) error {
	return c.Session.DB(r.DbName()).C("extensions").DropCollection()
}

func (r *MongoDbMigrator) V001_init_features_up(c *bongo.Connection) error {
	var collection []codeflow.Feature
	var obj codeflow.Feature

	created, _ := time.Parse(time.RFC3339, "2017-03-29T17:15:27.172Z")
	obj = codeflow.Feature{
		ProjectId:    bson.ObjectIdHex("58dbe995df8ab3002a71dc08"),
		Message:      "docker.io/checkr/codeflow:latest",
		User:         "checkr",
		Hash:         "df600016edb26c48e1c999b28bb874257f65d037",
		ParentHash:   "41fa95f534d8f57433ed4a48de1620cba35f2c2d",
		ExternalLink: "",
		Ref:          "refs/heads/master",
		Created:      created,
	}
	obj.SetId(bson.ObjectIdHex("58dbeb6cdf8ab3002a71dc0a"))
	collection = append(collection, obj)

	for _, o := range collection {
		if err := c.Collection("features").Save(&o); err != nil {
			log.Printf("Save::Error: %v", err.Error())
			return err
		}
	}

	return nil
}

func (r *MongoDbMigrator) V001_init_features_down(c *bongo.Connection) error {
	return c.Session.DB(r.DbName()).C("features").DropCollection()
}

func (r *MongoDbMigrator) V001_init_releases_up(c *bongo.Connection) error {
	var collection []codeflow.Release
	var obj codeflow.Release

	obj = codeflow.Release{
		ProjectId:     bson.ObjectIdHex("58dbe995df8ab3002a71dc08"),
		HeadFeatureId: bson.ObjectIdHex("58dbeb6cdf8ab3002a71dc0a"),
		TailFeatureId: bson.ObjectIdHex("58dbeb6cdf8ab3002a71dc0a"),
		UserId:        bson.ObjectIdHex("58dbe954df8ab3002a71dc07"),
		State:         "complete",
		StateMessage:  "",
		Secrets:       []codeflow.Secret{},
		Services:      []codeflow.Service{},
	}
	obj.SetId(bson.ObjectIdHex("58dbf989df8ab300cb0e4af6"))
	collection = append(collection, obj)

	for _, o := range collection {
		if err := c.Collection("releases").Save(&o); err != nil {
			log.Printf("Save::Error: %v", err.Error())
			return err
		}
	}

	return nil
}

func (r *MongoDbMigrator) V001_init_releases_down(c *bongo.Connection) error {
	return c.Session.DB(r.DbName()).C("releases").DropCollection()
}

func (r *MongoDbMigrator) V001_init_secrets_up(c *bongo.Connection) error {

	var collection []codeflow.Secret
	var obj codeflow.Secret

	// Add the viper configuration as ENVs
	allViperKeys := viper.AllKeys()
	for _, vKey := range allViperKeys {
		if vKey == "run" {
			continue
		}
		// Todo: handle arrays: eg. allowed_origins
		upcaseKey := strings.Replace(strings.ToUpper(vKey), ".", "_", -1)
		upcaseCfKey := strings.Join([]string{"CF_", upcaseKey}, "")
		obj = codeflow.Secret{
			ProjectId: bson.ObjectIdHex("58dbe995df8ab3002a71dc08"),
			Key:       upcaseCfKey,
			Value:     viper.GetString(vKey),
			Type:      "env",
			Deleted:   false,
		}
		collection = append(collection, obj)
	}

	// Add the additional REACT config variables from the ENV (not available in viper)
	obj = codeflow.Secret{
		ProjectId: bson.ObjectIdHex("58dbe995df8ab3002a71dc08"),
		Key:       "REACT_APP_PORT",
		Value:     getenv("REACT_APP_PORT", "9000"),
		Type:      "env",
		Deleted:   false,
	}
	collection = append(collection, obj)

	obj = codeflow.Secret{
		ProjectId: bson.ObjectIdHex("58dbe995df8ab3002a71dc08"),
		Key:       "REACT_APP_API_ROOT",
		Value:     getenv("REACT_APP_API_ROOT", "https://codeflow-api.example.net"),
		Type:      "env",
		Deleted:   false,
	}
	collection = append(collection, obj)

	obj = codeflow.Secret{
		ProjectId: bson.ObjectIdHex("58dbe995df8ab3002a71dc08"),
		Key:       "REACT_APP_ROOT",
		Value:     getenv("REACT_APP_ROOT", "https://codeflow.example.net"),
		Type:      "env",
		Deleted:   false,
	}
	collection = append(collection, obj)

	obj = codeflow.Secret{
		ProjectId: bson.ObjectIdHex("58dbe995df8ab3002a71dc08"),
		Key:       "REACT_APP_WEBHOOKS_ROOT",
		Value:     getenv("REACT_APP_WEBHOOKS_ROOT", "https://codeflow-webhooks.example.net"),
		Type:      "env",
		Deleted:   false,
	}
	collection = append(collection, obj)

	obj = codeflow.Secret{
		ProjectId: bson.ObjectIdHex("58dbe995df8ab3002a71dc08"),
		Key:       "REACT_APP_WS_ROOT",
		Value:     getenv("REACT_APP_WS_ROOT", "wss://codeflow-websockets.example.net"),
		Type:      "env",
		Deleted:   false,
	}
	collection = append(collection, obj)

	obj = codeflow.Secret{
		ProjectId: bson.ObjectIdHex("58dbe995df8ab3002a71dc08"),
		Key:       "REACT_APP_OKTA_CLIENT_ID",
		Value:     getenv("REACT_APP_OKTA_CLIENT_ID", "dummy"),
		Type:      "env",
		Deleted:   false,
	}
	collection = append(collection, obj)

	obj = codeflow.Secret{
		ProjectId: bson.ObjectIdHex("58dbe995df8ab3002a71dc08"),
		Key:       "REACT_APP_OKTA_LOGO",
		Value:     getenv("REACT_APP_OKTA_LOGO", "https://ok4static.oktacdn.com/bc/image/fileStoreRecord?id=dummy"),
		Type:      "env",
		Deleted:   false,
	}
	collection = append(collection, obj)

	obj = codeflow.Secret{
		ProjectId: bson.ObjectIdHex("58dbe995df8ab3002a71dc08"),
		Key:       "REACT_APP_OKTA_URL",
		Value:     getenv("REACT_APP_OKTA_URL", "https://dummy.okta.com"),
		Type:      "env",
		Deleted:   false,
	}
	collection = append(collection, obj)

	for _, o := range collection {
		if err := c.Collection("secrets").Save(&o); err != nil {
			log.Printf("Save::Error: %v", err.Error())
			return err
		}
	}

	return nil
}

func (r *MongoDbMigrator) V001_init_secrets_down(c *bongo.Connection) error {
	return c.Session.DB(r.DbName()).C("secrets").DropCollection()
}

func (r *MongoDbMigrator) V001_init_service_specs_up(c *bongo.Connection) error {
	var collection []codeflow.ServiceSpec
	var obj codeflow.ServiceSpec

	obj = codeflow.ServiceSpec{
		Name:                          "General-purpose",
		CpuRequest:                    "500m",
		CpuLimit:                      "1000m",
		MemoryRequest:                 "512Mi",
		MemoryLimit:                   "1Gi",
		TerminationGracePeriodSeconds: 600,
		Default: true,
	}
	obj.SetId(bson.ObjectIdHex("589bb6d6b158cdb147ef5dd0"))
	collection = append(collection, obj)

	obj = codeflow.ServiceSpec{
		Name:                          "Console",
		CpuRequest:                    "500m",
		CpuLimit:                      "1000m",
		MemoryRequest:                 "512Mi",
		MemoryLimit:                   "1Gi",
		TerminationGracePeriodSeconds: 86400,
	}
	obj.SetId(bson.ObjectIdHex("589cb50eb158cdb147f9cb5c"))
	collection = append(collection, obj)

	for _, o := range collection {
		if err := c.Collection("serviceSpecs").Save(&o); err != nil {
			log.Printf("Save::Error: %v", err.Error())
			return err
		}
	}

	return nil
}

func (r *MongoDbMigrator) V001_init_service_specs_down(c *bongo.Connection) error {
	return c.Session.DB(r.DbName()).C("serviceSpecs").DropCollection()
}

func (r *MongoDbMigrator) V001_init_services_up(c *bongo.Connection) error {
	var collection []codeflow.Service
	var obj codeflow.Service

	obj = codeflow.Service{
		ProjectId:    bson.ObjectIdHex("58dbe995df8ab3002a71dc08"),
		SpecId:       bson.ObjectIdHex("589bb6d6b158cdb147ef5dd0"),
		State:        "running",
		StateMessage: "",
		Name:         "api",
		Count:        1,
		Command:      "WORKDIR=./server /go/bin/codeflow --config ./configs/codeflow.yml server --run=git_sync,kubedeploy,heartbeat,docker_build,slack,route53,webhooks,codeflow,websockets",
		Listeners: []codeflow.Listener{
			{
				Port:     3001,
				Protocol: "TCP",
			},
			{
				Port:     3002,
				Protocol: "TCP",
			},
			{
				Port:     3003,
				Protocol: "TCP",
			},
		},
	}
	obj.SetId(bson.ObjectIdHex("58dbecefdf8ab3002a71dc0e"))
	collection = append(collection, obj)

	obj = codeflow.Service{
		ProjectId:    bson.ObjectIdHex("58dbe995df8ab3002a71dc08"),
		SpecId:       bson.ObjectIdHex("589bb6d6b158cdb147ef5dd0"),
		State:        "running",
		StateMessage: "",
		Name:         "www",
		Count:        1,
		Command:      "WORKDIR=./dashboard node server.js",
		Listeners: []codeflow.Listener{
			{
				Port:     9000,
				Protocol: "TCP",
			},
		},
	}
	obj.SetId(bson.ObjectIdHex("58dbed10df8ab3002a71dc0f"))
	collection = append(collection, obj)

	obj = codeflow.Service{
		ProjectId:    bson.ObjectIdHex("58dbe995df8ab3002a71dc08"),
		SpecId:       bson.ObjectIdHex("589bb6d6b158cdb147ef5dd0"),
		State:        "running",
		StateMessage: "",
		Name:         "docs",
		Count:        0,
		Command:      "WORKDIR=./docs node server.js",
		Listeners: []codeflow.Listener{
			{
				Port:     3000,
				Protocol: "TCP",
			},
		},
	}
	obj.SetId(bson.ObjectIdHex("58dbeeaadf8ab3002a71dc14"))
	collection = append(collection, obj)

	for _, o := range collection {
		if err := c.Collection("services").Save(&o); err != nil {
			log.Printf("Save::Error: %v", err.Error())
			return err
		}
	}

	return nil
}

func (r *MongoDbMigrator) V001_init_services_down(c *bongo.Connection) error {
	return c.Session.DB(r.DbName()).C("services").DropCollection()
}

func (r *MongoDbMigrator) V001_init_workflows_up(c *bongo.Connection) error {
	var collection []codeflow.Flow

	var obj codeflow.Flow

	obj = codeflow.Flow{
		ReleaseId: bson.ObjectIdHex("58dbf938df8ab300cb0e4af4"),
		Type:      "Build",
		Name:      "DockerImage",
		Message:   "",
		State:     "complete",
	}
	obj.SetId(bson.ObjectIdHex("58dbf939df8ab300cb0e4af5"))
	collection = append(collection, obj)

	obj = codeflow.Flow{
		ReleaseId: bson.ObjectIdHex("58dbf989df8ab300cb0e4af6"),
		Type:      "Build",
		Name:      "DockerImage",
		Message:   "",
		State:     "complete",
	}
	obj.SetId(bson.ObjectIdHex("58dbf989df8ab300cb0e4af7"))
	collection = append(collection, obj)

	for _, o := range collection {
		if err := c.Collection("workflows").Save(&o); err != nil {
			log.Printf("Save::Error: %v", err.Error())
			return err
		}
	}

	return nil
}

func (r *MongoDbMigrator) V001_init_workflows_down(c *bongo.Connection) error {
	return c.Session.DB(r.DbName()).C("workflows").DropCollection()
}

func (r *MongoDbMigrator) V002_projects_deleted_up(c *bongo.Connection) error {
	return c.Session.DB(r.DbName()).C("projects").Update(bson.M{}, bson.M{"$set": bson.M{"deleted": false}})
}

func (r *MongoDbMigrator) V002_projects_deleted_down(c *bongo.Connection) error {
	return c.Session.DB(r.DbName()).C("projects").Update(bson.M{}, bson.M{"$unset": bson.M{"deleted": ""}})
}

func (r *MongoDbMigrator) V003_r53_up(c *bongo.Connection) error {
	if _, err := c.Session.DB(r.DbName()).C("extensions").UpdateAll(bson.M{}, bson.M{"$rename": bson.M{"dnsName": "dns"}}); err != nil {
		if err != mgo.ErrNotFound {
			return err
		}
	}

	if viper.GetInt("plugins.route53.workers") <= 0 {
		return nil
	}

	// sync existing dns records
	records := make(map[string]string)

	// Create the client
	sess := session.Must(session.NewSessionWithOptions(
		session.Options{
			Config: aws.Config{
				Credentials: credentials.NewStaticCredentials(viper.GetString("plugins.route53.aws_access_key_id"), viper.GetString("plugins.route53.aws_secret_key"), ""),
			},
		},
	))
	client := route53.New(sess)

	// Look for this dns name
	params := &route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(viper.GetString("plugins.route53.hosted_zone_id")), // Required
	}
	pageNum := 0

	errList := client.ListResourceRecordSetsPages(params,
		func(page *route53.ListResourceRecordSetsOutput, lastPage bool) bool {
			pageNum++
			for _, p := range page.ResourceRecordSets {
				if *p.Type == "CNAME" {
					records[*p.ResourceRecords[0].Value] = strings.Split(*p.Name, ".")[0]
				}
			}
			return false
		})

	if errList != nil {
		log.Printf("Error listing ResourceRecordSets for Route53: %s", errList)
		return errList
	}

	var loadBalancers []codeflow.LoadBalancer
	if err := c.Session.DB(r.DbName()).C("extensions").Find(bson.M{"extension": "LoadBalancer"}).All(&loadBalancers); err != nil {
		if err != mgo.ErrNotFound {
			return err
		}
	}

	for _, lb := range loadBalancers {
		if lb.Type == plugins.Internal {
			continue
		}

		if lb.DNS == "" {
			continue
		}

		for key, value := range records {
			if lb.DNS == key && lb.Subdomain == "" {
				if err := c.Session.DB(r.DbName()).C("extensions").Update(bson.M{"_id": lb.Id}, bson.M{"$set": bson.M{"subdomain": value}}); err != nil {
					return err
				}
				log.Printf("Updating %s: %s linked to %s", lb.Id, key, value)
			}
		}
	}
	return nil
}

func (r *MongoDbMigrator) V003_r53_down(c *bongo.Connection) error {
	if _, err := c.Session.DB(r.DbName()).C("extensions").UpdateAll(bson.M{}, bson.M{"$rename": bson.M{"dnsName": "dns"}}); err != nil {
		if err != mgo.ErrNotFound {
			return err
		}
	}

	if _, err := c.Session.DB(r.DbName()).C("extensions").UpdateAll(bson.M{}, bson.M{"$unset": bson.M{"subdomain": ""}}); err != nil {
		if err != mgo.ErrNotFound {
			return err
		}
	}

	return nil
}
