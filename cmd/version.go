// Package cmd is for any command line arguments this application utilizes
package cmd

import "github.com/sirupsen/logrus"

// VersionCommand is a struct to enclose all version related sub commands if any
type VersionCommand struct {
}

// Version is the version string for this application assigned during CI/CD
var Version string = "Development"

// Run is the method that is executed when the version command is selected
func (cmd *VersionCommand) Run() error {
	logrus.Infof("Version: %s", Version)
	return nil
}
