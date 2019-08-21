package logger

import "github.com/sirupsen/logrus"

func Default() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{})
}
