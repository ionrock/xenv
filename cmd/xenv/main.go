package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/ionrock/xenv/config"
	"github.com/urfave/cli"
)

var builddate = ""
var gitref = ""

// XeAction runs the main command.
func XeAction(c *cli.Context) error {

	if c.Bool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	logCtx := log.WithFields(log.Fields{
		"config": c.String("config"),
	})
	logCtx.Debug("Loading config")
	env, err := config.NewEnvironmentFromConfig(c.String("config"))
	if err != nil {
		logCtx.WithFields(log.Fields{"error": err}).Error("error loading config")
		return err
	}
	env.DataOnly = c.Bool("data")

	if c.Bool("data") {
		err = env.Pre()
		if err != nil {
			return err
		}
		for _, pair := range env.Config.ToEnv() {
			fmt.Println(pair)
		}
		return nil
	}

	return env.Main(c.Args())
}

func main() {
	app := cli.NewApp()

	app.Version = fmt.Sprintf("%s-%s", gitref, builddate)

	app.Name = "xenv"
	app.Usage = "Start and monitor processes creating an executable environment."
	app.ArgsUsage = "[COMMAND]"
	app.Action = XeAction

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Path to the xe config file, default is ./xe.yml",
			Value: "xe.yml",
		},

		cli.BoolFlag{
			Name:  "data, d",
			Usage: "Only compute the data and print it out.",
		},

		cli.BoolFlag{
			Name:  "debug, D",
			Usage: "Print debugging output.",
		},
	}

	app.Run(os.Args)
}
