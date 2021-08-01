// Package topics for handling all topics
package topics

import (
	"fmt"
	"strings"
)

const (
	// SettingsType is indicative of a settings data object
	SettingsType = 1
	// SettingsTopic is the topic for settings data
	SettingsTopic = "settings"
	// ImageType is indicative of a image data object
	ImageType = 2
	// ImageTopic is the topic for image data
	ImageTopic = "image"
	// VideoType is indicative of a video data object
	VideoType = 3
	// VideoTopic is the topic for video data
	VideoTopic = "video"
	// AudioType is indicative of a audio data object
	AudioType = 4
	// AudioTopic is the topic for audio data
	AudioTopic = "audio"

	// NumberTopicPortions is the number of parts of an incoming topic that is expected
	NumberTopicPortions = 4
)

// DeviceData interface that all device data packets implement
type DeviceData interface {
	GetData() []byte
	SetData(incomingData []byte) error
	GetDeviceID() string
	SetDeviceID(string)
	GetType() int
}

// DeviceSettings struct which represent settings that can be pushed to a device
type DeviceSettings struct {
	data     []byte
	deviceID string
}

// GetData implements DeviceData intereface to return the data
func (ds *DeviceSettings) GetData() []byte {
	return ds.data
}

// SetData sets the data for the incoming object
func (ds *DeviceSettings) SetData(incomingData []byte) error {
	ds.data = incomingData

	return nil
}

// GetDeviceID returns the associated device id with the data
func (ds *DeviceSettings) GetDeviceID() string {
	return ds.deviceID
}

// SetDeviceID set the associated device id
func (ds *DeviceSettings) SetDeviceID(identifier string) {
	ds.deviceID = identifier
}

// GetType returns the settings type
func (ds *DeviceSettings) GetType() int {
	return SettingsType
}

// ImageData struct representing image data from a client
type ImageData struct {
	data     []byte
	deviceID string
}

// GetData implements DeviceData intereface to return the data
func (img *ImageData) GetData() []byte {
	return img.data
}

// SetData sets the data for the incoming object
func (img *ImageData) SetData(incomingData []byte) error {
	img.data = incomingData

	return nil
}

// GetDeviceID returns the associated device id with the data
func (img *ImageData) GetDeviceID() string {
	return img.deviceID
}

// SetDeviceID set the associated device id
func (img *ImageData) SetDeviceID(identifier string) {
	img.deviceID = identifier
}

// GetType returns the image type
func (img *ImageData) GetType() int {
	return ImageType
}

// VideoData struct representing video data from a client
type VideoData struct {
	data     []byte
	deviceID string
}

// GetData implements DeviceData intereface to return the data
func (vid *VideoData) GetData() []byte {
	return vid.data
}

// SetData sets the data for the incoming object
func (vid *VideoData) SetData(incomingData []byte) error {
	vid.data = incomingData

	return nil
}

// GetDeviceID returns the associated device id with the data
func (vid *VideoData) GetDeviceID() string {
	return vid.deviceID
}

// SetDeviceID set the associated device id
func (vid *VideoData) SetDeviceID(identifier string) {
	vid.deviceID = identifier
}

// GetType returns the video type
func (vid *VideoData) GetType() int {
	return VideoType
}

// AudioData struct representing audio data from a client
type AudioData struct {
	data     []byte
	deviceID string
}

// GetData implements DeviceData intereface to return the data
func (aud *AudioData) GetData() []byte {
	return aud.data
}

// SetData sets the data for the incoming object
func (aud *AudioData) SetData(incomingData []byte) error {
	aud.data = incomingData

	return nil
}

// GetDeviceID returns the associated device id with the data
func (aud *AudioData) GetDeviceID() string {
	return aud.deviceID
}

// SetDeviceID set the associated device id
func (aud *AudioData) SetDeviceID(identifier string) {
	aud.deviceID = identifier
}

// GetType returns the audio type
func (aud *AudioData) GetType() int {
	return AudioType
}

func processAudioData(deviceID string, data []byte) (*AudioData, error) {
	var audioData AudioData

	audioData.SetDeviceID(deviceID)
	err := audioData.SetData(data)
	if err != nil {
		return nil, err
	}

	return &audioData, nil
}

func processVideoData(deviceID string, data []byte) (*VideoData, error) {
	var videoData VideoData

	videoData.SetDeviceID(deviceID)
	err := videoData.SetData(data)
	if err != nil {
		return nil, err
	}

	return &videoData, nil
}

func processImageData(deviceID string, data []byte) (*ImageData, error) {
	var imageData ImageData

	imageData.SetDeviceID(deviceID)
	err := imageData.SetData(data)
	if err != nil {
		return nil, err
	}

	return &imageData, nil
}

func processDeviceSettings(deviceID string, data []byte) (*DeviceSettings, error) {
	var deviceSettings DeviceSettings

	deviceSettings.SetDeviceID(deviceID)
	err := deviceSettings.SetData(data)
	if err != nil {
		return nil, err
	}

	return &deviceSettings, nil
}

// ProcessIncomingMQTTMessage to process an incoming message
// topics should be afm/v1/<action>/device_id
func ProcessIncomingMQTTMessage(topic, message string) (DeviceData, error) {
	// device id
	topicParts := strings.Split(topic, "/")

	if len(topicParts) >= NumberTopicPortions {
		deviceID := topicParts[NumberTopicPortions-1]
		deviceData := []byte(message)

		// which topic?
		topics := []string{SettingsTopic, ImageTopic, VideoTopic, AudioTopic}
		for _, v := range topics {
			if strings.Contains(topic, v) {
				switch v {
				case SettingsTopic:
					return processDeviceSettings(deviceID, deviceData)
				case ImageTopic:
					return processImageData(deviceID, deviceData)
				case VideoTopic:
					return processVideoData(deviceID, deviceData)
				case AudioTopic:
					return processAudioData(deviceID, deviceData)
				}
			}
		}
	}
	return nil, fmt.Errorf("failed to process message")
}
