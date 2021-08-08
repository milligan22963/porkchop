// Package server is made up of modules related to the web server
package server

import (
	"net/http"

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

// HandleFavoriteIcon will serve up a favorite as desired
func HandleFavoriteIcon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.ico")
}
