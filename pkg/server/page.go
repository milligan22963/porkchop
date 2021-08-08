// Package server for all server related items
package server

import (
	"fmt"
	"io"
	"net/http"
	"site/pkg/tools"
)

// Meta is data added to the page
type Meta struct {
	Name    string
	Content string
}

// Page is a typical page
type Page struct {
	title       string
	styleSheets []string
	javaScript  []string
	meta        []Meta
	body        []string
}

type HtmlPage interface {
	Render(w io.Writer, r *http.Request) error
	SetTitle(title string)
	AddStyleSheet(styleSheet string)
	AddJavaScript(scriptFile string)
	AddMetaData(name, content string)
}

// Render will consider the input request and render the appropriate response
func (page *Page) Render(w io.Writer, r *http.Request) error {
	fmt.Fprintln(w, "<html><head>")
	if len(page.title) > 0 {
		fmt.Fprintf(w, "<title>%s</title>\n", page.title)
	}

	// include style sheets
	for _, sheet := range page.styleSheets {
		// fmt.Printf("<link rel=\"stylesheet\" type=\"text/css\" href=\"%s\"/>\n", sheet)
		fmt.Fprintf(w, "<link type=\"text/css\" rel=\"stylesheet\" href=\"%s\"/>\n", sheet)
	}

	// js files
	for _, script := range page.javaScript {
		fmt.Fprintf(w, "<script type=\"application/javascript\" src=\"%s/>\n", script)
	}

	// meta data
	for _, meta := range page.meta {
		fmt.Fprintf(w, "<meta name=\"%s\" content=\"%s\"/>\n", meta.Name, meta.Content)
	}

	fmt.Fprintf(w, "</head>")
	fmt.Fprintf(w, "<body>Welcome to the HomePage!</body></html>")

	return nil
}

// SetTitle sets the title for this page
func (page *Page) SetTitle(title string) {
	page.title = title
}

// AddStyleSheet adds a style sheet to this page
func (page *Page) AddStyleSheet(styleSheet string) {
	// see if its there already, if not add it otherwise skip it
	scriptIndex := tools.Find(page.styleSheets, styleSheet)
	// if not found then add it
	if scriptIndex == -1 {
		page.styleSheets = append(page.styleSheets, styleSheet)
	}
}

func (page *Page) AddJavaScript(scriptFile string) {
	// see if its there already, if not add it otherwise skip it
	scriptIndex := tools.Find(page.javaScript, scriptFile)
	// if not found then add it
	if scriptIndex == -1 {
		page.javaScript = append(page.styleSheets, scriptFile)
	}
}

func (page *Page) AddMetaData(name, content string) {
	found := false
	for _, metaData := range page.meta {
		if metaData.Name == name {
			metaData.Content = content
			found = true
		}
	}

	if !found {
		newMetaData := Meta{Name: name, Content: content}
		page.meta = append(page.meta, newMetaData)
	}
}
