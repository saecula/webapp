package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Post struct {
	Email    string
	Username string
	Title    string
	Body     []byte
}

type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Title    string `json:"title"`
	Body     string `json:"body"`
}

var GET = "GET"
var POST = "POST"
var PUT = "PUT"
var DELETE = "DELETE"

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Key)
var upgrader = websocket.Upgrader{}

var templates = template.Must(template.ParseFiles("edit.html", "view.html", "notfound.html"))
var validPath = regexp.MustCompile("^([a-zA-Z0-0-9-]+)$")

func getTitle(r *http.Request) (string, error) {
	title := strings.Split(r.URL.Path, "/")[1]
	if len(title) == 0 {
		return "", nil
	}
	if isValid := validPath.MatchString(title); !isValid {
		return "", errors.New("Invalid Post Title")
	}
	return title, nil
}

func (p *Post) save() error {
	filename := "posts/" + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func (p *Post) delete() error {
	filename := "posts/" + p.Title + ".txt"
	return os.Remove(filename)
}

func loadPost(title string) (*Post, error) {
	filename := "posts/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Post{Title: title, Body: body}, nil
}

func loadPostAsMessage(title string) (*Message, error) {
	filename := "posts/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	strBody := string(body)
	return &Message{Title: title, Body: strBody}, nil
}

func saveMessageAsPost(msg Message) {
	pageBody := []byte(msg.Body)

	var pageTitle string
	if msg.Title == "" {
		pageTitle = uuid.NewString()
	} else {
		pageTitle = msg.Title
	}

	p := &Post{Title: pageTitle, Body: pageBody}
	p.save()
}

func enableCors(w *http.ResponseWriter) {
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}
	fmt.Println("host: " + host)
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func mainHandler(w http.ResponseWriter, r *http.Request, game string) {
	fmt.Println("ping!")

	switch r.Method {
	case GET:
		serveGames(w, r)
	case POST:
		saveGame(w, r, game)
	case PUT:
		saveGame(w, r, game)
	case DELETE:
		deleteGame(w, r, game)
	}
}

type Game struct {
	Game string `json:"game"`
	Next string `json:"next"`
}

func save(g *Game) error {
	filename := "games/19.json"
	gb, _ := json.Marshal(g)
	return ioutil.WriteFile(filename, gb, 0600)
}

func serveGames(w http.ResponseWriter, r *http.Request) {
	gamefiles, err := ioutil.ReadDir("./games")
	if err != nil {
		log.Fatal(err)
	}
	var gb []byte
	for _, f := range gamefiles {
		gb, err = ioutil.ReadFile("games/" + f.Name())
		if err != nil {
			log.Fatal(err)
		}
	}
	w.WriteHeader(201)
	w.Write(gb)
	return
}

func servePost(w http.ResponseWriter, r *http.Request, t string) {
	p, err := loadPost(t)
	if err != nil {
		p = &Post{Title: "new"}
	}
	json, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(201)
	w.Write(json)
}

func saveGame(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	b := []byte(body)
	var g *Game
	json.Unmarshal(b, g)

	err := save(g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/"+title, http.StatusCreated)
}

// todo: make for game
func deleteGame(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Post{Title: title, Body: []byte(body)}
	err := p.delete()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/", http.StatusAccepted)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		title, err := getTitle(r)
		if err != nil {
			w.WriteHeader(400)
			nf, err := json.Marshal("bad request")
			if err != nil {
				panic(err)
			}
			w.Write(nf)
			return
		}
		fn(w, r, title)
	}
}

type Key struct {
	Key string `json:"key"`
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var key Key
		err := ws.ReadJSON(&key)
		if err != nil {
			log.Printf("error in handleConnections: %#v", err)
			delete(clients, ws)
			break
		}

		broadcast <- key
	}
}

func handleMessages() {
	for {
		key := <-broadcast
		log.Printf("received message: %v", key)

		// for client := range clients {
		// 	err := client.WriteJSON(msg)
		// 	if err != nil {
		// 		log.Printf("error in handleMessages: %v", err)
		// 		client.Close()
		// 		delete(clients, client)
		// 	}
		// }
	}
}

func main() {
	log.SetFlags(0)
	log.Printf("hello am running :3")
	http.HandleFunc("/", makeHandler(mainHandler))

	http.HandleFunc("/ws", handleConnections)
	go handleMessages()

	log.Fatal(http.ListenAndServe(":4000", nil))
}
