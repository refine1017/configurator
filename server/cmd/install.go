package cmd

import (
	"github.com/urfave/cli"
	"server/models"
	"server/modules"
	"server/modules/manager"
	"server/modules/setting"
)

var Install = cli.Command{
	Name:   "install",
	Usage:  "Install Web server",
	Action: runInstall,
	Flags:  []cli.Flag{},
}

func runInstall(ctx *cli.Context) error {
	if err := modules.Startup(setting.APP_MODE_WEB); err != nil {
		return err
	}

	if err := manager.Startup(); err != nil {
		return err
	}
	return models.Install()
}
