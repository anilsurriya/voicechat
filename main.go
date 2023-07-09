package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	vcserver "github.com/anilsurriya/voicechat/server"
)

// templ represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	srv := vcserver.NewServer(":8080", "/room/")

	r := srv.NewRoom()
	if r == nil {
		fmt.Println()
	}
	http.Handle("/room/", srv)
	http.Handle("/", &templateHandler{filename: "chat.html"})

	log.Println("Starting web server on", srv.Addr)
	if err := http.ListenAndServe(srv.Addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
