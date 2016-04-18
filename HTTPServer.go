package main

import (
	// "fmt"
	"bomberman-server/gamemanager"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type HTTPServer struct {
	channel     chan string
	mainChannel chan string
	port        string
	game        *gamemanager.Game
}

func NewHTTPServer() *HTTPServer {
	ch := make(chan string)

	return &HTTPServer{
		channel: ch,
		port:    "8080",
	}
}

func (httpServer *HTTPServer) start() {
	go httpServer.handleHTTPChannel()

	// start http server
	address := "localhost:" + httpServer.port

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		viewHandler(w, r, httpServer)
	})

	log.Fatal(http.ListenAndServe(address, nil))
}

// receives all information about http
func (httpServer *HTTPServer) handleHTTPChannel() {
	for {
		var x = <-httpServer.channel
		//fmt.Printf("httpServer: %s\n", x)
		switch x {
		case "":

		}
	}
}

// http stuff
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

// func handler2(w http.ResponseWriter, r *http.Request) {
// 	var err error
//
// 	tpl := template.New("tpl.gohtml")
// 	tpl = tpl.Funcs(template.FuncMap{
// 		"uppercase": func(str string) string {
// 			return strings.ToUpper(str)
// 		},
// 	})
// 	tpl, err = tpl.ParseFiles("tpl.gohtml")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	err = tpl.Execute(os.Stdout, Page{
// 		Title: "My Title 2",
// 		Body:  `hello world <script>alert("Yow!");</script>`,
// 	})
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// }

// ##############################

// Page struct, which will contain template
type Page struct {
	Title   string
	Body    template.HTML
	Players []*gamemanager.Player
}

// Loads a page for use
func loadPage(title string, r *http.Request, game *gamemanager.Game) (*Page, error) {
	body, err := ioutil.ReadFile(title)
	if err != nil {
		return nil, err
	}

	//testplayers := []string{"player 1", "player 2", "player 3"}
	return &Page{Title: title, Body: template.HTML(body), Players: game.GetPlayersArray()}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, httpServer *HTTPServer) {
	httpServer.mainChannel <- r.URL.Path
	// Parses URL to obtain title of file to add to .body
	title := r.URL.Path[len("/"):]

	// Load templatized page, given title
	page, _ := loadPage(title, r, httpServer.game)

	// Generate template t
	t, _ := template.ParseFiles("index.html")

	// Write the template attributes of page (from load page) to t
	t.Execute(w, page)
}
