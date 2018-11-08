package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	markdown "gopkg.in/russross/blackfriday.v2"
)

// Post : struct for a blog post
type Post struct {
	URL      string
	HTML     template.HTML
	Title    string
	Markdown string
}

func (post Post) parse() error {

	return nil
}

var posts []Post

func main() {
	mux := httprouter.New()

	// Fileserver
	mux.ServeFiles("/static/*filepath", http.Dir(config.Static))

	// Handlers
	mux.GET("/", index)
	mux.GET("/blog/:id", getPosts)
	// Server Config
	server := &http.Server{
		Addr:           config.Address,
		Handler:        mux,
		ReadTimeout:    time.Duration(config.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}

func index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	templates := template.Must(template.ParseFiles("templates/index.html"))
	err := templates.ExecuteTemplate(w, "index", posts)
	if err != nil {
		fmt.Println(err)
	}
}

func parseMarkdown(post *Post) {
	post.HTML = template.HTML(string(markdown.Run([]byte(post.Markdown))))

}

func getPosts(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	templates := template.Must(template.ParseFiles("templates/layout.html"))
	index, err := strconv.ParseInt(p.ByName("id"), 10, 32)
	parseMarkdown(&posts[index])
	err = templates.ExecuteTemplate(w, "layout", posts[index])
	if err != nil {
		fmt.Println(err)
	}
}
