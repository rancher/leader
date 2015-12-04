package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/rancher/go-rancher-metadata/metadata"
	"github.com/rancher/leader/election"
)

var (
	port        = "proxy-tcp-port"
	check       = "check"
	metadataUrl = "http://rancher-metadata/2015-07-25"
)

func main() {
	app := cli.NewApp()
	app.Author = "Rancher Labs, Inc."
	app.EnableBashCompletion = true
	app.Usage = "Simple leader election with Rancher"
	app.Action = appAction

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  port,
			Usage: "Port to proxy to the leader",
		},
		cli.BoolFlag{
			Name:  check,
			Usage: "Check if we are the leader and exit",
		},
	}

	app.Run(os.Args)
}

func appAction(cli *cli.Context) {
	client, err := metadata.NewClientAndWait(metadataUrl)
	if err != nil {
		logrus.Fatal(err)
	}

	w := election.New(client, cli.Int(port), cli.Args())

	if cli.Bool(check) {
		if w.IsLeader() {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	if err := w.Watch(); err != nil {
		logrus.Fatal(err)
	}
}
