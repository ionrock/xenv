package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ionrock/xenv/config"
	"github.com/ionrock/xenv/manager"

	"github.com/ghodss/yaml"
	"github.com/urfave/cli"
)

var builddate = ""
var gitref = ""

// // XeAction runs the main command.
// func XeAction(c *cli.Context) error {

// 	if c.Bool("debug") {
// 		log.SetLevel(log.DebugLevel)
// 	}

// 	logCtx := log.WithFields(log.Fields{
// 		"config": c.String("config"),
// 	})
// 	logCtx.Debug("Loading config")
// 	env, err := config.NewEnvironmentFromConfig(c.String("config"))
// 	if err != nil {
// 		logCtx.WithFields(log.Fields{"error": err}).Error("error loading config")
// 		return err
// 	}
// 	env.DataOnly = c.Bool("data")

// 	if c.Bool("data") {
// 		err = env.Pre()
// 		if err != nil {
// 			return err
// 		}
// 		for _, pair := range env.Config.ToEnv() {
// 			fmt.Println(pair)
// 		}
// 		return nil
// 	}

// 	return env.Main(c.Args())
// }

func loadConfig(path string) ([]*config.Service, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := make([]*config.Service, 0)

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func initialize(c *cli.Context) error {
	mgr := manager.New()
	cfg, err := loadConfig(c.String("config"))
	if err != nil {
		return err
	}

	fmt.Println(cfg)

	for _, svc := range cfg {
		mgr.Start(svc.Name, svc.Cmd, svc.Dir, os.Environ())
	}

	return mgr.Wait()
}

func stop(c *cli.Context) error {
	return nil
}

func start(c *cli.Context) error {
	return nil
}

func restart(c *cli.Context) error {
	return nil
}

func main() {
	app := cli.NewApp()
	app.Version = fmt.Sprintf("%s-%s", gitref, builddate)
	app.Name = "svcs"
	app.Usage = "A simple process manager."

	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "start up processes defined in the config",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config, c",
					Usage: "Path to the svcs config file",
					Value: "svcs.yml",
				},

				cli.BoolFlag{
					Name:  "debug, D",
					Usage: "Print debugging output.",
				},
			},

			Action: initialize,
		},
		{
			Name:   "start",
			Usage:  "start up a single process",
			Action: start,
		},

		{
			Name:   "stop",
			Usage:  "stop a managed process",
			Action: stop,
		},
		{
			Name:   "restart",
			Usage:  "restart a managed process",
			Action: restart,
		},
	}

	app.Run(os.Args)
}
