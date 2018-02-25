package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ionrock/xenv/config"
	"github.com/ionrock/xenv/manager"

	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
	"github.com/urfave/cli"
)

var builddate = ""
var gitref = ""

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
		log.Infof("Starting: %s %s", svc.Name, svc.Cmd)
		mgr.Start(svc.Name, svc.Cmd, svc.Dir, os.Environ())
	}

	return mgr.Watch()
}

func main() {
	app := cli.NewApp()
	app.Version = fmt.Sprintf("%s-%s", gitref, builddate)
	app.Name = "svcs"
	app.Usage = "A simple process manager."
	app.Action = initialize
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Path to the svcs config file",
			Value: "svcs.yml",
		},

		cli.BoolFlag{
			Name:  "debug, D",
			Usage: "Print debugging output.",
		},
	}

	app.Run(os.Args)
}
