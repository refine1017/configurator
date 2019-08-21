package setting

import (
	"github.com/Unknwon/com"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

var (
	Cfg      *ini.File
	ConfFile = "conf/app.ini"
)

type parser func() error

var parsers = []parser{parseApp}

func registerParser(p parser) {
	parsers = append(parsers, p)
}

func Load() error {
	Cfg = ini.Empty()

	if com.IsFile(ConfFile) {
		if err := Cfg.Append(ConfFile); err != nil {
			logrus.Fatalf("Failed to load custom conf '%s': %v", ConfFile, err)
		}
	} else {
		logrus.Warnf("Config '%s' not found, ignore this if you're running first time", ConfFile)
	}
	Cfg.NameMapper = ini.AllCapsUnderscore

	for _, parser := range parsers {
		if err := parser(); err != nil {
			return err
		}
	}

	return nil
}
