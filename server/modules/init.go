package modules

import (
	"github.com/sirupsen/logrus"
	"server/modules/logger"
	"server/modules/setting"
)

func Startup(mode string) error {
	logger.Default()

	setting.AppMode = mode

	if err := setting.Load(); err != nil {
		logrus.Fatal("setting Load with err: %v", err)
		return err
	}

	logrus.SetLevel(setting.LogLevel)

	return nil
}
