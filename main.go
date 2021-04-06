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

	"github.com/gorilla/websocket"
)

type Page struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Title    string `json:"title"`
	Body     []byte `json:"body"`
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

var templates = template.Must(template.ParseFiles("edit.html", "view.html", "main.html", "notfound.html"))
var validPath = regexp.MustCompile("^([a-zA-Z0-0-9-]+)$")

func getTitle(r *http.Request) (string, error) {
	title := strings.Split(r.URL.Path, "/")[1]
	if len(title) == 0 {
		return "", nil
	}
	if isValid := validPath.MatchString(title); !isValid {
		return "", errors.New("Invalid Page Title")
	}
	return title, nil
}

func (p *Page) save() error {
	filename := "pages/" + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func (p *Page) delete() error {
	filename := "pages/" + p.Title + ".txt"
	return os.Remove(filename)
}

func loadPage(title string) (*Page, error) {
	filename := "pages/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func loadPageAsMessage(title string) (*Message, error) {
	filename := "pages/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	strBody := string(body)
	return &Message{Title: title, Body: strBody}, nil
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:5000")
}

func mainHandler(w http.ResponseWriter, r *http.Request, title string) {
	fmt.Println("ping!")

	switch r.Method {
	case GET:
		if title == "" {
			serveMessages(w, r)
		} else {
			servePage(w, r, title)
		}
	case POST:
		savePage(w, r, title)
	case PUT:
		savePage(w, r, title)
	case DELETE:
		deletePage(w, r, title)
	}
}

func serveMessages(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir("./pages")
	if err != nil {
		log.Fatal(err)
	}
	var messages []Message
	for _, f := range files {
		rawTitle := f.Name()
		title := strings.TrimSuffix(rawTitle, filepath.Ext(rawTitle))
		m, err := loadPageAsMessage(title)
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

func servePage(w http.ResponseWriter, r *http.Request, t string) {
	p, err := loadPage(t)
	if err != nil {
		p = &Page{Title: "new"}
	}
	json, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(201)
	w.Write(json)
}

func savePage(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/"+title, http.StatusCreated)
}

func deletePage(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
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
	fmt.Println("wut")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("nuuu")
		log.Fatal(err)
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var data Message

		err := ws.ReadJSON(&data)
		log.Printf("omg maybe %v", data)

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
