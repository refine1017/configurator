package setting

import "github.com/sirupsen/logrus"

// Mongo
var Mongo struct {
	Url    string `ini:"MONGO_URL"`
	DBName string `ini:"MONGO_DB_NAME"`
}

func parseMongo() error {
	if err := Cfg.Section("Mongo").MapTo(&Mongo); err != nil {
		return err
	}

	logrus.Infof("MONGO_URL=%v", Mongo.Url)

	return nil
}

func init() {
	registerParser(parseMongo)
}
