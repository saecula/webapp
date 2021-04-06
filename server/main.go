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
	"path/filepath"
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
var broadcast = make(chan Message)
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
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func mainHandler(w http.ResponseWriter, r *http.Request, title string) {
	fmt.Println("ping!")

	switch r.Method {
	case GET:
		if title == "" {
			serveMessages(w, r)
		} else {
			servePost(w, r, title)
		}
	case POST:
		savePost(w, r, title)
	case PUT:
		savePost(w, r, title)
	case DELETE:
		deletePost(w, r, title)
	}
}

func serveMessages(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir("./posts")
	if err != nil {
		log.Fatal(err)
	}
	var messages []Message
	for _, f := range files {
		rawTitle := f.Name()
		title := strings.TrimSuffix(rawTitle, filepath.Ext(rawTitle))
		m, err := loadPostAsMessage(title)
		if err != nil {
			panic(err)
		}
		messages = append(messages, *m)
	}
	json, err := json.Marshal(messages)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(201)
	w.Write(json)
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

func savePost(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Post{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/"+title, http.StatusCreated)
}

func deletePost(w http.ResponseWriter, r *http.Request, title string) {
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

func handleConnections(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var data Message
		err := ws.ReadJSON(&data)
		if err != nil {
			log.Printf("error in handleConnections: %#v", err)
			delete(clients, ws)
			break
		}

		broadcast <- data
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		log.Printf("received message: %v", msg)

		saveMessageAsPost(msg)

		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error in handleMessages: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	http.HandleFunc("/", makeHandler(mainHandler))

	http.HandleFunc("/ws", handleConnections)
	go handleMessages()

	log.Fatal(http.ListenAndServe(":4000", nil))
}
