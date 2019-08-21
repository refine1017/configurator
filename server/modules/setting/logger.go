package setting

import "github.com/sirupsen/logrus"

var (
	LogLevel logrus.Level
)

func parseLogger() error {
	logLevel := Cfg.Section("").Key("LOG_LEVEL").MustUint(uint(logrus.InfoLevel))
	LogLevel = logrus.Level(logLevel)

	return nil
}

func init() {
	registerParser(parseLogger)
}
