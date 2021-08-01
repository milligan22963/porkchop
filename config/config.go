// Package config is used to define any configuration that isn't passed in from the command line
// or is default options that can be overridden
package config

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	identificationPath = "/var/cache/afm/identifier.id"
	defaultSQLPort     = 3306
	defaultSQLType     = "mysql"
	defaultMQTTPort    = 1883
	defaultFileOptions = 0600
	defaultWebPort     = 8080
)

// ConfigurationDetails stores the configuration that will be used
var ConfigurationDetails = map[string]interface{}{
	BrokerAddress:        "127.0.0.1",
	BrokerPort:           defaultMQTTPort,
	BrokerSSL:            false,
	BrokerPrivateKeyPath: "/etc/afm/ssl",
	BrokerPublicKeyPath:  "/etc/afm/ssl",
	BrokerCAPath:         "/etc/afm/ssl",

	DatabaseName: "afmcamera",
	DatabaseHost: "localhost",
	DatabasePort: defaultSQLPort,
	DatabaseType: defaultSQLType,

	LoggingUseFile: true,
	LoggingFile:    "/var/log/afm/camera.log",
	LoggingLevel:   "error",
	LoggingFormat:  "text",

	WebServerAddress: "localhost",
	WebServerFiles:   "/var/www/html",
	WebServerPort:    defaultWebPort,
}

// DefaultConfigPath to our default config
const DefaultConfigPath = "/etc/afm/camera.yaml"

// DatabaseQueryResponse struct representing the response to a given query
type DatabaseQueryResponse struct {
	Query    string
	Response *sqlx.Rows
	Err      error
}

// SiteConfiguration is configuration
type SiteConfiguration struct {
	AppActive    chan struct{}
	IncomingMQTT chan [2]string
	OutgoingMQTT chan [3]string
	ClientID     string
	Database     *sqlx.DB
}

func determineCurrentIP() (string, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		// = GET LOCAL IP ADDRESS
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("unable to find ip address")
}

func determineCurrentNetworkHardwareInterface(currentIP string) (string, error) {
	// get all the system's or local machine's network interfaces
	interfaces, interfaceerr := net.Interfaces()
	for _, interf := range interfaces {

		if addrs, err := interf.Addrs(); err == nil {
			for _, addr := range addrs {
				// only interested in the name with current IP address
				if strings.Contains(addr.String(), currentIP) {
					return interf.Name, nil
				}
			}
		}
	}
	return "", interfaceerr
}

func determineDeviceMACAddress() (string, error) {

	currentIP, err := determineCurrentIP()
	if err != nil {
		return "", err
	}

	hardwareInterfaceName, err := determineCurrentNetworkHardwareInterface(currentIP)

	if err != nil {
		return "", err
	}

	// extract the hardware information base on the interface name
	// capture above
	netInterface, err := net.InterfaceByName(hardwareInterfaceName)

	if err != nil {
		return "", err
	}

	macAddress := netInterface.HardwareAddr

	// verify if the MAC address can be parsed properly
	hwAddr, err := net.ParseMAC(macAddress.String())

	if err != nil {
		return "", err
	}

	return hwAddr.String(), nil
}

func determineDeviceSerialNumber() (string, error) {
	out, err := exec.Command("/usr/sbin/dmidecode", "-t", "system").Output()
	if err == nil {
		return "", err
	}
	for _, l := range strings.Split(string(out), "\n") {
		if strings.Contains(l, "Serial Number") {
			s := strings.Split(l, ":")
			return s[len(s)-1], nil
		}
	}
	return "", fmt.Errorf("unable to find serial number")
}

func determineDeviceClientID() string {
	// Try reading well known file location
	identifier, err := ioutil.ReadFile(identificationPath)

	if err == nil {
		cleansedIdentifier := strings.TrimRight(string(identifier), "\n")
		return cleansedIdentifier
	}

	// Try reading hardware assigned serial number
	serialNumber, err := determineDeviceSerialNumber()
	if err == nil {
		return serialNumber
	}

	// Default back to mac address
	macAddress, err := determineDeviceMACAddress()
	if err == nil {
		return macAddress
	}

	return "id_failure"
}

func initializeConfigurationOptions(configFilePath string) {
	for k, v := range ConfigurationDetails {
		viper.SetDefault(k, v)
	}

	viper.SetConfigFile(configFilePath)

	err := viper.ReadInConfig()
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Warnf("configuration file: %s does not exist.", configFilePath)
		} else {
			logrus.Errorf("configuration file is invalid: %v", err)
		}
	}
}

func initializeLogging() {
	desiredLevel := logrus.ErrorLevel

	var err error
	if viper.IsSet(LoggingLevel) {
		desiredLevel, err = logrus.ParseLevel(viper.GetString(LoggingLevel))
		if err != nil {
			desiredLevel = logrus.ErrorLevel // if failure then set to error
		}
	}

	logrus.SetLevel(desiredLevel)

	loggingFormat := viper.GetString(LoggingFormat)
	if loggingFormat == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{})
	}

	if viper.IsSet(LoggingUseFile) && viper.IsSet(LoggingFile) {
		file, err := os.OpenFile(viper.GetString(LoggingFile), os.O_CREATE|os.O_WRONLY|os.O_APPEND, defaultFileOptions)
		// if no failure output to file, otherwise default to stderr/stdouts
		if err == nil {
			logrus.SetOutput(file)
		}
	}
}

func setupDatabase(initialDBNameConnect bool) *sqlx.DB {
	connectionString := viper.GetString(DatabaseUser) + ":" + viper.GetString(DatabasePassword)
	connectionString += "@tcp(" + viper.GetString(DatabaseHost) + ":" + strconv.Itoa(viper.GetInt(DatabasePort)) + ")/"

	if initialDBNameConnect {
		connectionString += viper.GetString(DatabaseName)
	}

	database, err := sqlx.Open("mysql", connectionString)

	if err != nil {
		panic(err.Error())
	}

	return database
}

// NewSiteConfiguration creates an instance of the site configuration struct
func NewSiteConfiguration(cfgFileOverride string, initialDBNameConnect bool) *SiteConfiguration {

	viper.SetEnvPrefix("camera")
	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Use default, then try env defined override
	// finally look at passed in command line option
	configPath := DefaultConfigPath
	if viper.IsSet(OverrideConfigFile) {
		configPath = viper.GetString(OverrideConfigFile)
	}

	if len(cfgFileOverride) > 0 {
		configPath = cfgFileOverride
	}

	initializeConfigurationOptions(configPath)

	initializeLogging()

	siteConfig := &SiteConfiguration{
		AppActive:    make(chan struct{}),
		IncomingMQTT: make(chan [2]string),
		OutgoingMQTT: make(chan [3]string),
		ClientID:     determineDeviceClientID(),
		Database:     setupDatabase(initialDBNameConnect),
	}

	return siteConfig
}
