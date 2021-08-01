// Package cmd is for any command line arguments this application utilizes
package cmd

import (
	"io/ioutil"
	"path/filepath"
	"site/config"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// InitializeCommand is a struct to enclose all initialization related sub commands if any
type InitializeCommand struct {
	ConfigurationFile string `short:"c" help:"Defines the non-default configuration file to use."`
	DatabaseFile      string `short:"d" help:"Defines the database file to initialize with."`
}

func createDatabase(sqlFile string, database *sqlx.DB) error {
	file, err := ioutil.ReadFile(filepath.Clean(sqlFile))
	if err != nil {
		return err
	}

	// Convert to strings
	sqlStatements := strings.Split(string(file), "\n")
	fullStatement := ""
	for _, statement := range sqlStatements {
		if len(statement) > 0 && !strings.Contains(statement, "--") && !strings.Contains(statement, "/*") {
			if statement[len(statement)-1] == '\n' {
				statement = statement[:len(statement)-1]
			}
			fullStatement += statement
			if strings.Contains(fullStatement, ";") {
				logrus.Infof("query: %s", fullStatement)
				_, err = database.Exec(fullStatement)
				if err != nil {
					logrus.Errorf("query failed on initialize: %v", err)
					return err
				}
				fullStatement = ""
			}
		}
	}
	return nil
}

// Run is the method that is executed when the version command is selected
func (cmd *InitializeCommand) Run() error {
	logrus.Info("Initializing system for first usage")

	// Setup configuration and logging
	// Will use passed in configuration file if any
	siteConfig := config.NewSiteConfiguration(cmd.ConfigurationFile, false)

	targetDatabaseName := viper.GetString(config.DatabaseName)

	_, err := siteConfig.Database.Query("drop database " + targetDatabaseName)
	if err != nil {
		logrus.Warnf("failed to drop database: %v", err)
	}

	// Create the new database
	_, err = siteConfig.Database.Query("create database " + targetDatabaseName)
	if err != nil {
		logrus.Errorf("failed to create new database: %v", err)
		return err
	}

	// Select the database
	err = siteConfig.Database.Close()
	if err != nil {
		logrus.Warnf("unable to close the database connection: %v", err)
	}

	connectionString := viper.GetString(config.DatabaseUser) + ":" + viper.GetString(config.DatabasePassword)
	connectionString += "@tcp(" + viper.GetString(config.DatabaseHost) + ":" + strconv.Itoa(viper.GetInt(config.DatabasePort)) + ")/"

	siteConfig.Database, err = sqlx.Open(viper.GetString(config.DatabaseType), connectionString+targetDatabaseName)

	if err != nil {
		panic(err.Error())
	}

	// Once database is dropped we can load the new one
	err = createDatabase(cmd.DatabaseFile, siteConfig.Database)

	_ = siteConfig.Database.Close()

	return err
}
