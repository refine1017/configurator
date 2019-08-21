package setting

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"server/modules/util"
)

var Router struct {
	Host string `ini:"host"`
	Port int    `ini:"Port"`
}

func parseRouter() error {
	sec := Cfg.Section("router")

	Router.Host = sec.Key("Url").MustString(util.GetIntranetIp())
	Router.Port = sec.Key("Port").MustInt(8000)

	AppURL = fmt.Sprintf("http://%s:%d", Router.Host, Router.Port)

	logrus.Infof("Http Listen %v:%v", Router.Host, Router.Port)

	return nil
}

func init() {
	registerParser(parseRouter)
}
