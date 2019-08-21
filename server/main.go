package main

import (
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/urfave/cli"
	"log"
	"net/http"
	"os"
	"server/cmd"
	"server/modules/setting"
	"server/public"
)

func main() {
	fs := assetfs.AssetFS{
		Asset:     public.Asset,
		AssetDir:  public.AssetDir,
		AssetInfo: public.AssetInfo,
	}
	http.Handle("/", http.FileServer(&fs))

	app := cli.NewApp()
	app.Name = "Configurator"
	app.Usage = "A support multi format config tools"
	app.Description = ""
	app.Version = "1.0.0"
	app.Commands = []cli.Command{
		cmd.Web,
		cmd.Install,
	}

	// default configuration flags
	defaultFlags := []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Value:       setting.ConfFile,
			Usage:       "configuration file path",
			Destination: &setting.ConfFile,
		},
		cli.VersionFlag,
	}

	// Set the default to be equivalent to cmdWeb and add the default flags
	app.Flags = append(app.Flags, cmd.Web.Flags...)
	app.Flags = append(app.Flags, defaultFlags...)
	app.Action = cmd.Web.Action

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal("Failed to run app with %s: %v", os.Args, err)
	}
}
