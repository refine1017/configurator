package mongo

import (
	"gopkg.in/mgo.v2"
	"server/modules/manager"
	"server/modules/setting"
)

var Session *mgo.Session

func DB() *mgo.Database {
	return Session.DB(setting.Mongo.DBName)
}

func connect() error {
	var err error

	Session, err = mgo.Dial(setting.Mongo.Url)
	if err != nil {
		return err
	}

	// Optional. Switch the session to a monotonic behavior.
	Session.SetMode(mgo.Monotonic, true)

	return nil
}

func disconnect() {
	if Session != nil {
		Session.Close()
	}
}

func init() {
	manager.OnStartup(connect)
	manager.OnShutdown(disconnect)
}
