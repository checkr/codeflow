package codeflow_db

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	mgo "gopkg.in/mgo.v2"
)

type Config struct {
	URI   string
	SSL   bool
	Creds mgo.Credential
}

func NewConnection(config Config) (*mgo.Session, error) {
	var dialInfo *mgo.DialInfo
	var err error

	if dialInfo, err = mgo.ParseURL(config.URI); err != nil {
		panic(fmt.Sprintf("cannot parse given URI %s due to error: %s", config.URI, err.Error()))
	}

	if config.SSL {
		tlsConfig := &tls.Config{}
		tlsConfig.InsecureSkipVerify = true
		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			return conn, err
		}
	}

	dialInfo.Timeout = time.Second * 3
	dialInfo.Username = config.Creds.Username
	dialInfo.Password = config.Creds.Password
	dialInfo.Mechanism = config.Creds.Mechanism

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, err
	}

	session.Login(&config.Creds)

	return session, nil
}
