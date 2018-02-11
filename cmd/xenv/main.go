package main

import (
	"fmt"
	"os"
	"os/exec"

	log "github.com/Sirupsen/logrus"
	"github.com/ionrock/xenv/config"
	"github.com/ionrock/xenv/util"
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

	cfgs, err := config.NewXeConfig(c.String("config"))
	if err != nil {
		logCtx.WithFields(log.Fields{"error": err}).Error("error loading config")
		return err
	}

	configDir, err := util.AbsDir(c.String("config"))
	if err != nil {
		return err
	}

	env := config.NewEnvironment(configDir)
	env.DataOnly = c.Bool("data")

	err = env.Pre(cfgs)
	if err != nil {
		return err
	}

	if c.Bool("data") {
		for _, pair := range env.Config.ToEnv() {
			fmt.Println(pair)
		}
		return nil
	}

	parts := c.Args()
	if len(parts) > 0 {
		fmt.Printf("Going to start: %s\n", parts)
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Env = env.Config.ToEnv()

		err = cmd.Run()
	}

	fmt.Println("running post now")
	postErr := env.Post()
	if postErr != nil {
		fmt.Printf("Error running post: %s\n", postErr)
	}

	if err != nil {
		return err
	}

	return nil
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
