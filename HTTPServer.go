package main

import (
  // "fmt"
  "net/http"
  "io/ioutil"
  "html/template"
)

type HTTPServer struct {
  channel chan string
}

func NewHTTPServer() *HTTPServer {
  ch := make(chan string)

  return &HTTPServer{
    channel: ch,
  }
}

func (httpServer *HTTPServer) start() {
  // start http server
  http.HandleFunc("/", viewHandler)
  http.ListenAndServe("localhost:8000", nil)
}

// http Stuff
func handler(w http.ResponseWriter, r *http.Request) {
  // html := "hallo"
  //fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
  // fmt.Fprintf(w, html)
  path := r.URL.Path[1:]
  data, err := ioutil.ReadFile(string(path))

  if err == nil {
    w.Write(data)
  } else {
    w.WriteHeader(404)
    w.Write([]byte("404 Pate not found"))
  }
}

// ##############################

// Page struct, which will contain template
type Page struct {
	Title   string
	Body    template.HTML
  Number  int
}

// Loads a page for use
func loadPage(title string, r *http.Request) (*Page, error) {
	body, err := ioutil.ReadFile(title)
	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: template.HTML(body), Number: 23}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
  httpServer.channel <- r.URL.Path
        // Parses URL to obtain title of file to add to .body
	title := r.URL.Path[len("/"):]

        // Load templatized page, given title
	page, _ := loadPage(title, r)

        // Generate template t
	t, _ := template.ParseFiles("index.html")

        // Write the template attributes of page (from load page) to t
	t.Execute(w, page)
}
