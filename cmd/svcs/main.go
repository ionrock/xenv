package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/ionrock/xenv/config"
	"github.com/ionrock/xenv/manager"

	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

func getServerCreds() (credentials.TransportCredentials, error) {
	// Create the TLS credentials
	creds, err := credentials.NewServerTLSFromFile("cert.pem", "key.pem")
	if err == nil {
		return creds, nil
	}

	log.Errorf("could not load TLS keys: %s", err)
	manager.GenerateCerts()
	return credentials.NewServerTLSFromFile("cert.pem", "key.pem")

}

func getClientCreds() (credentials.TransportCredentials, error) {
	return credentials.NewClientTLSFromFile("cert.pem", "")
}

func initialize(c *cli.Context) error {
	mgr := manager.New()
	cfg, err := loadConfig(c.String("config"))
	if err != nil {
		return err
	}

	fmt.Println(cfg)

	lis, err := net.Listen("tcp", "127.0.0.1:9909")
	if err != nil {
		log.WithError(err).Fatal("failed to bind to tcp socket")
	}

	creds, err := getServerCreds()
	if err != nil {
		log.WithError(err).Fatal("failed to load certs")
	}

	opts := []grpc.ServerOption{grpc.Creds(creds)}

	grpcServer := grpc.NewServer(opts...)
	defer grpcServer.Stop()

	manager.RegisterSvcsManagerServer(grpcServer, &manager.ManagerServer{mgr})

	go func() {
		grpcServer.Serve(lis)
	}()

	for _, svc := range cfg {
		log.Infof("Starting: %s %s", svc.Name, svc.Cmd)
		mgr.Start(svc.Name, svc.Cmd, svc.Dir, os.Environ())
	}

	return mgr.Watch()
}

func getClient() (manager.SvcsManagerClient, error) {
	creds, err := getClientCreds()
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial("127.0.0.1:9909", grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}

	cc := manager.NewSvcsManagerClient(conn)
	return cc, nil
}

func restart(c *cli.Context) error {
	name := c.Args().Get(0)
	if name == "" {
		return errors.New("need a process name")
	}

	fmt.Println("getting client")
	client, err := getClient()
	if err != nil {
		fmt.Println("error getting client")
		fmt.Println(err)
		return err
	}

	fmt.Println("making request to restart " + name)
	result, err := client.Restart(context.Background(), &manager.SvcProcess{name})
	if err != nil {
		log.WithError(err).Fatal("restart failed")
	}

	fmt.Println("output")
	fmt.Println(result.Output)
	return nil
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

	app.Commands = []cli.Command{
		{
			Name:   "restart",
			Usage:  "restart NAME",
			Action: restart,
		},
	}

	app.Run(os.Args)
}
