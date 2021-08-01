package main

import (
	"site/cmd"

	"github.com/alecthomas/kong"
	"github.com/sirupsen/logrus"
)

// cli is an internal command structure to pass into kong
var cli struct {
	Initialize cmd.InitializeCommand `cmd:"" help:"Initialize the system"`
	Run        cmd.RunCommand        `cmd:"" help:"Run this application"`
	Version    cmd.VersionCommand    `cmd:"" help:"version: Print version and exit"`
}

// main is our primary application starting point
func main() {
	context := kong.Parse(&cli)

	err := context.Run()

	if err != nil {
		logrus.Fatal(err)
	}
}
