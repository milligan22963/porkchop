// Package cmd is for any command line arguments this application utilizes
package cmd

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"site/config"
	"site/pkg/database"
	"site/pkg/server"
	"site/pkg/topics"
	"strconv"
	"sync"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"

	// Pulling in mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

const (
	mqttWait = 250
	httpWait = 5 * time.Second
)

// RunCommand is a struct to enclose all run related sub commands if any
type RunCommand struct {
	ConfigurationFile string `short:"c" help:"Defines the non-default configuration file to use."`
}

func (cmd *RunCommand) processSettingsObject(db *sqlx.DB, device topics.DeviceData) error {
	// look up device by id (serial)
	deviceObj := database.DeviceObject{}

	// test out a database seting
	err := deviceObj.LoadByField(db, device.GetDeviceID())
	if err == nil {
		logrus.Infof("retrieved setting: %v", deviceObj)
	} else {
		// do we need to add this one
		deviceObj.Serial = device.GetDeviceID()
		deviceObj.Active = 1
		if json.Unmarshal(device.GetData(), &deviceObj) != nil {
			logrus.Error("unable to unmarshal json device data")
		}

		err = deviceObj.Create(db)
		if err != nil {
			logrus.Errorf("failed to add device to database: %v", err)
		}
	}

	return err
}

func (cmd *RunCommand) processImageObject(db *sqlx.DB, device topics.DeviceData) error {
	// look up device by id (serial)
	deviceObj := database.DeviceObject{}

	// look up device to determine user id
	err := deviceObj.LoadByField(db, device.GetDeviceID())
	if err != nil {
		return err
	}

	logrus.Infof("retrieved device: %v", deviceObj)
	deviceUserMap := database.DeviceUserMappingObject{}
	err = deviceUserMap.LoadByField(db, strconv.Itoa(deviceObj.ID))
	if err != nil {
		return err
	}

	userObj := database.UserObject{ID: deviceUserMap.UserID}
	err = userObj.Load(db)
	if err != nil {
		return err
	}

	// Get cache file location
	cacheStorage := viper.GetString(config.WebServerCache)

	rawData := device.GetData()
	fileNameLength := rawData[0]
	fileName := cacheStorage + userObj.UserName + "/" + string(rawData[1:fileNameLength])
	imageSizeTemp := rawData[fileNameLength+1 : 4] // copy to byte align
	imageSize := binary.LittleEndian.Uint32(imageSizeTemp)

	imageFile, err := os.Create(fileName)
	if err != nil {
		return err
	}

	bytesWritten, err := imageFile.Write(rawData[fileNameLength+5:])

	if err != nil {
		return err
	}

	if uint32(bytesWritten) != imageSize {
		return fmt.Errorf("failed to write complete image file....%d of size %d", bytesWritten, imageSize)
	}

	// Create new image for this user
	imageObj := database.ImageObject{UserID: deviceUserMap.UserID, DeviceID: deviceUserMap.DeviceID}
	imageObj.Path = fileName
	err = imageObj.Create(db)
	if err != nil {
		return err
	}
	// Notify of new image
	return nil
}

func (cmd *RunCommand) processMQTTRequest(siteConfig *config.SiteConfiguration, topic, message string) error {
	deviceData, mqttErr := topics.ProcessIncomingMQTTMessage(topic, message)
	if mqttErr == nil {
		logrus.Infof("RECEIVED type: %d device: %s", deviceData.GetType(), deviceData.GetDeviceID())

		switch deviceData.GetType() {
		case topics.SettingsType:
			err := cmd.processSettingsObject(siteConfig.Database, deviceData)
			if err != nil {
				return err
			}

		case topics.ImageType:
			// Store image data for this user in the database
			err := cmd.processImageObject(siteConfig.Database, deviceData)
			if err != nil {
				return err
			}
		case topics.VideoType:
		case topics.AudioType:
		default:
		}
	}
	return mqttErr
}

func (cmd *RunCommand) process(siteConfig *config.SiteConfiguration) error {

	onQuit := make(chan os.Signal, 1)
	signal.Notify(onQuit, syscall.SIGINT, syscall.SIGTERM)

	var quitReason os.Signal
	for quitReason == nil {
		select {
		case incomingMQTT := <-siteConfig.IncomingMQTT:
			logrus.Infof("Received mqtt request: %s, %s", incomingMQTT[0], incomingMQTT[1])
			if cmd.processMQTTRequest(siteConfig, incomingMQTT[0], incomingMQTT[1]) != nil {
				logrus.Warnf("failed to process mqtt request: %s, %s", incomingMQTT[0], incomingMQTT[1])
			}
		case quitReason = <-onQuit:
			logrus.Info("applcation is now exiting on signal")
		}
	}
	if quitReason != syscall.SIGINT {
		return fmt.Errorf("shutting down due to signal: %+v", quitReason)
	}
	return nil
}

// Run is the method that is executed when the run command is selected
func (cmd *RunCommand) Run() error {

	logrus.Info("Starting up application")

	siteConfig := config.NewSiteConfiguration(cmd.ConfigurationFile, true)

	// connect to mqtt
	go setupMQTTMessages(siteConfig)

	// server up the world
	go setupWebserver(siteConfig)

	quitReason := cmd.process(siteConfig)

	siteConfig.AppActive <- struct{}{}

	_ = siteConfig.Database.Close()

	return quitReason
}

func setupMQTTMessages(siteConfig *config.SiteConfiguration) {
	prefix := "ssl://"
	if !viper.GetBool(config.BrokerSSL) {
		// configure for tcp
		prefix = "tcp://"
	}

	logrus.Infof("ClientID: %s", siteConfig.ClientID)
	broker := prefix + viper.GetString(config.BrokerAddress)
	broker += ":" + strconv.Itoa(viper.GetInt(config.BrokerPort))

	opts := MQTT.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(siteConfig.ClientID)
	opts.SetCleanSession(true)
	// opts.SetStore(MQTT.NewFileStore("path to store")) // default is memory

	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		siteConfig.IncomingMQTT <- [2]string{msg.Topic(), string(msg.Payload())}
	})

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := client.Subscribe("afm/v1/#", 1, nil); token.Wait() && token.Error() != nil {
		logrus.Errorf(token.Error().Error())
	}

	logrus.Infof("Broker: %s", broker)

	quit := false
	for !quit {
		select {
		case outgoingMessage := <-siteConfig.OutgoingMQTT:
			qos, err := strconv.Atoi(outgoingMessage[2])
			if err != nil {
				qos = 0
			}

			token := client.Publish(outgoingMessage[0], byte(qos), false, outgoingMessage[1])
			if token.Error() != nil {
				logrus.Errorf("failed publishing message: %v", token.Error())
			}
		// wait for the app to go down
		case <-siteConfig.AppActive:
			quit = true
		}
	}
	client.Disconnect(mqttWait)
}

func setupWebserver(siteConfig *config.SiteConfiguration) {
	httpServerDone := &sync.WaitGroup{}
	router := mux.NewRouter().StrictSlash(true)

	serverPort := viper.GetInt(config.WebServerPort)
	serverAddress := viper.GetString(config.WebServerAddress)
	router.HandleFunc("/", server.GenerateHomePage)
	server := &http.Server{Addr: serverAddress + ":" + strconv.Itoa(serverPort), Handler: router}

	logrus.Infof("http server: %v", server.Addr)

	httpServerDone.Add(1) // Add before our go routine
	go func() {
		defer httpServerDone.Done()

		if err := server.ListenAndServe(); err != nil {
			logrus.Errorf("http listen error: %v", err)
		}
	}()

	<-siteConfig.AppActive

	ctx, cancel := context.WithTimeout(context.Background(), httpWait)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logrus.Errorf("server shutdown error: %v", err)
	}

	// wait for the server func to finish
	httpServerDone.Wait()
}
