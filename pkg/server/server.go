// Package server is made up of modules related to the web server
package server

import (
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

// HTTPResponse is a structure defining what a response should look like
type HTTPResponse struct {
	Code    int    `json:"-"`
	Message string `json:"message,omitempty"`
}

// GenerateHomePage generates the home page for this site
func GenerateHomePage(w http.ResponseWriter, r *http.Request) {
	homePage := Page{title: "AFM"}

	homePage.AddStyleSheet("polaroid.css")

	err := homePage.Render(w, r)
	if err != nil {
		logrus.Errorf("failed to render: %v", err)
	}
}

// RetrieveLiveImage returns the live image from the camera
func RetrieveLiveImage(w http.ResponseWriter, r *http.Request) {
	logrus.Info("serving up an image")
	file, err := os.Open("web/images/avatar.webp")
	if err != nil {
		logrus.Errorf("unable to open image")
		return
	}

	defer file.Close()

	w.Header().Set("Content-Type", "image/webp")
	fileInfo, err := file.Stat()
	if err != nil {
		logrus.Error("failed to stat file")
		return
	}
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
	bytesWritten, err := io.Copy(w, file)
	if err != nil {
		logrus.Error("failed to serve file")
		return
	}
	if bytesWritten != fileInfo.Size() {
		logrus.Error("incorrect amount of data copied")
	}
}

// HandleFavoriteIcon will serve up a favorite as desired
func HandleFavoriteIcon(w http.ResponseWriter, r *http.Request) {
	logrus.Error("favicon served w/ chocolate")
	http.ServeFile(w, r, "favicon.ico")
}
