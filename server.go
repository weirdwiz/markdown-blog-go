package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"

	markdown "gopkg.in/russross/blackfriday.v2"
)

// Post : struct for a blog post
type Post struct {
	html     string
	date     time.Time
	title    string
	markdown string
}

func (post Post) parse() error {

	return nil
}

func main() {
	parseMarkdown()
	mux := http.NewServeMux()

	// Fileserver
	files := http.FileServer(http.Dir(config.Static))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	// Handlers
	mux.HandleFunc("/", index)

	server := &http.Server{
		Addr:           config.Address,
		Handler:        mux,
		ReadTimeout:    time.Duration(config.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}

func index(w http.ResponseWriter, r *http.Request) {
	templates := template.Must(template.ParseFiles("templates/layout.html"))
	templates.ExecuteTemplate(w, "layout", template.HTML(parseMarkdown()))
}

func parseMarkdown() string {
	file, err := ioutil.ReadFile("test.md")
	if err != nil {
		panic(err)
	}
	parsed := markdown.Run(file)
	return string(parsed)
}

func generateHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", data)
}
