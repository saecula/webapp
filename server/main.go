package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Move string

const (
	// switch colors, only valid before first turn
	Switch Move = "switch"
	// play a stone
	Play Move = "play"
	// pass your turn
	Pass Move = "pass"
	// resign the game
	Resign Move = "resign"
)

type TurnMessage struct {
	GameId       string      `json:"gameId"`
	PlayerId     string      `json:"playerId"`
	Move         Move        `json:"move"`
	Point        string      `json:"point"`
	FinishedTurn bool        `json:"finishedTurn"`
	BoardTemp    interface{} `json:"boardtemp"`
}

type GameStateMessage struct {
	Id         string      `json:"id"`
	Board      interface{} `json:"board"`
	NextPlayer string      `json:"nextPlayer"`
}

type GameState struct {
	Id         string      `json:"id"`
	Board      interface{} `json:"board"`
	NextPlayer string      `json:"nextPlayer"`
	Players    *PlayerMap  `json:"playerMap"`
	Started    bool        `json:"started"`
	Ended      bool        `json:"ended"`
	Winner     string      `json:"winner"`
}

type PlayerMap struct {
	B string `json:"b"`
	W string `json:"w"`
}

type Player struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

var GET = "GET"
var POST = "POST"
var PUT = "PUT"
var DELETE = "DELETE"

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan TurnMessage)
var upgrader = websocket.Upgrader{}

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

func calcGame(tm *TurnMessage) *GameState {
	prevGame, _ := loadGame(tm.GameId)
	log.Printf("prev game %v", prevGame.Id)
	if prevGame == nil {
		return &GameState{}
	}
	var started bool
	if tm.Move == Switch {
		started = false
	} else {
		started = true
	}
	var ended bool
	if tm.Move == Resign {
		ended = true
	} else {
		ended = false
	}

	var nextPlayer string
	if tm.FinishedTurn == false {
		nextPlayer = tm.PlayerId
	} else if tm.PlayerId == "b" {
		nextPlayer = "w"
	} else {
		nextPlayer = "b"
	}

	return &GameState{
		Id:         prevGame.Id,
		Board:      tm.BoardTemp,
		NextPlayer: nextPlayer,
		Players:    prevGame.Players,
		Started:    started,
		Ended:      ended,
		Winner:     "",
	}
}

func (g *GameState) save() error {
	filename := "games/" + g.Id + ".json"
	v, err := json.Marshal(g)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return ioutil.WriteFile(filename, v, 0600)
}

func loadGame(id string) (*GameState, error) {
	filename := "games/" + id + ".json"
	body, err := ioutil.ReadFile(filename)
	var new bool
	if err != nil {
		log.Println("defaulting back to new game")
		new = true
		body, _ = ioutil.ReadFile("games/newgame.json")
	}
	var g *GameState
	err = json.Unmarshal(body, &g)
	if new {
		g.Id = uuid.NewString()
	}
	return g, nil
}

func enableCors(w *http.ResponseWriter) {
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}
	fmt.Println("host: " + host)
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
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
		var turn TurnMessage
		err := ws.ReadJSON(&turn)
		if err != nil {
			log.Printf("error in handleConnections: %#v", err)
			delete(clients, ws)
			break
		}

		broadcast <- turn
	}
}

func handleMessages() {
	for {
		turnmsg := <-broadcast
		log.Printf("received message: %v", turnmsg)
		gamemsg := calcGame(&turnmsg)
		log.Printf("sending message %v", gamemsg)
		for client := range clients {
			err := client.WriteJSON(gamemsg)
			if err != nil {
				log.Printf("error in handleMessages: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
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

func mainHandler(w http.ResponseWriter, r *http.Request, game string) {
	fmt.Println("httping!")
	id := strings.TrimPrefix(r.URL.Path, "/")
	log.Printf("hello here is url %v and id %s", r.URL.RawPath, id)
	switch id {
	case "":
		serveNewGame(w)
	default:
		serveGame(w, id)
	}
}

func serveNewGame(w http.ResponseWriter) {
	gb, err := ioutil.ReadFile("games/newgame.json")
	if err != nil {
		log.Fatal(err)
	}
	var g *GameState
	if err := json.Unmarshal(gb, &g); err != nil {
		log.Printf("goddammit")
		log.Fatal(err)
	}
	g.Id = uuid.NewString()
	log.Printf("saving new game %s", g.Id)
	g.save()
	gameMsg := &GameStateMessage{
		Id:         g.Id,
		Board:      g.Board,
		NextPlayer: g.NextPlayer,
	}
	msg, _ := json.Marshal(&gameMsg)
	w.WriteHeader(201)
	w.Write(msg)
	return
}

func serveGame(w http.ResponseWriter, id string) {
	sg, err := ioutil.ReadFile(fmt.Sprintf("./games/%s", id))
	if err != nil {
		serveNewGame(w)
	} else {
		var g *GameState
		if err := json.Unmarshal(sg, &g); err != nil {
			log.Printf("goddammit id game")
			log.Fatal(err)
		}
		g.Id = uuid.NewString()
		log.Printf("saving new game %s", g.Id)
		gameMsg := &GameStateMessage{
			Id:         g.Id,
			Board:      g.Board,
			NextPlayer: g.NextPlayer,
		}
		msg, _ := json.Marshal(&gameMsg)
		w.WriteHeader(201)
		w.Write(msg)
		return
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
