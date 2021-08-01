// Package server for all server related items
package server

import (
	"fmt"
	"io"
	"net/http"
)

// Page is a typical page
type Page struct {
	title       string
	StyleSheets []string
	JavaScript  []string
	Body        []string
}

// Render will consider the input request and render the appropriate response
func (page *Page) Render(w io.Writer, r *http.Request) error {
	fmt.Fprintln(w, "<html><header>")
	fmt.Fprintf(w, "<title>%s</title>\n", page.title)
	for _, sheet := range page.StyleSheets {
		fmt.Fprintf(w, "<link rel=\"stylesheet\" type=\"text/css\" href=\"%s/>\n", sheet)
	}
	for _, script := range page.JavaScript {
		fmt.Fprintf(w, "<script type=\"application/javascript\" src=\"%s/>\n", script)
	}
	fmt.Fprintf(w, "</header>")
	fmt.Fprintf(w, "<body>Welcome to the HomePage!</body></html>")

	return nil
}

// SetTitle sets the title for this page
func (page *Page) SetTitle(text string) {
	page.title = text
}

// AddStyleSheet adds a style sheet to this page
func (page *Page) AddStyleSheet(styleSheet string) {
	// see if its there already, if not add it otherwise skip it
}
