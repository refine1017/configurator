package cmd

import (
	"github.com/urfave/cli"
	"server/modules"
	"server/modules/manager"
	"server/modules/setting"
	"server/routers"
)

var Web = cli.Command{
	Name:        "web",
	Usage:       "Start Web server",
	Description: "Configurator Web",
	Action:      runWeb,
	Flags:       []cli.Flag{},
}

func runWeb(ctx *cli.Context) error {
	if err := modules.Startup(setting.APP_MODE_WEB); err != nil {
		return err
	}

	if err := manager.Startup(); err != nil {
		return err
	}

	routers.Startup()

	return nil
}
