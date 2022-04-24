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
	"strconv"
	"strings"

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

type Turn struct {
	GameId       string                       `json:"gameId"`
	Player       string                       `json:"player"`
	Color        string                       `json:"color"`
	Move         Move                         `json:"move"`
	Point        string                       `json:"point"`
	FinishedTurn bool                         `json:"finishedTurn"`
	BoardTemp    map[string]map[string]string `json:"boardTemp"`
}

type GameState struct {
	Id         string                       `json:"id"`
	Board      map[string]map[string]string `json:"board"`
	LastPlayed string                       `json:"lastPlayed"`
	NextPlayer string                       `json:"nextPlayer"`
	Players    *PlayerMap                   `json:"players"`
	Started    bool                         `json:"started"`
	Ended      bool                         `json:"ended"`
	Winner     string                       `json:"winner"`
}

type PlayerMap struct {
	B string `json:"b"`
	W string `json:"w"`
}

var GET = "GET"
var POST = "POST"
var PUT = "PUT"
var DELETE = "DELETE"

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Turn)
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

func calcGame(tm *Turn) (*GameState, error) {
	prevGame, _ := loadGame(tm.GameId)
	log.Printf("prev game %v", prevGame.Id)
	if prevGame == nil {
		return &GameState{}, nil
	}
	var started bool
	newPlayers := prevGame.Players
	if tm.Move == Switch {
		if prevGame.Started == true {
			return nil, errors.New("invalid switch!")
		}
		newPlayers = &PlayerMap{
			B: prevGame.Players.W,
			W: prevGame.Players.B,
		}
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
		nextPlayer = tm.Player
	} else if tm.Move != Switch {
		if tm.Color == "b" {
			nextPlayer = prevGame.Players.W
		} else {
			nextPlayer = prevGame.Players.B
		}
	}

	valid := true
	if tm.Move == Play {
		valid = false
		pointState := makePointState(tm.Point, tm.BoardTemp)
		liberties := getSurroundingPoints(pointState, tm.BoardTemp, "e")
		if len(liberties) > 0 {
			valid = true
		}
		for _, sp := range getSurroundingPoints(pointState, tm.BoardTemp, oppositeOf(tm.Color)) {
			hasNoLiberties, checkedPoints := wholeThinghasNoLiberties(sp, tm.BoardTemp)
			if hasNoLiberties {
				valid = true
				removeWholeThing(checkedPoints, tm.BoardTemp)
			}
		}
	}

	if !valid {
		return nil, errors.New("not a valid move")
	}

	return &GameState{
		Id:         prevGame.Id,
		Board:      tm.BoardTemp,
		LastPlayed: tm.Point,
		NextPlayer: nextPlayer,
		Players:    newPlayers,
		Started:    started,
		Ended:      ended,
		Winner:     "",
	}, nil
}

type PointState struct {
	key   string
	x     int
	y     int
	state string
}

func makePointState(point string, board map[string]map[string]string) *PointState {
	coordinates := strings.Split(point, ":")
	row, err := strconv.Atoi(coordinates[0])
	if err != nil {
		log.Fatal("done gone wrong")
	}
	col, err := strconv.Atoi(coordinates[1])
	if err != nil {
		log.Fatal("done gone wrong")
	}
	return &PointState{
		key:   point,
		x:     row,
		y:     col,
		state: board[coordinates[0]][coordinates[1]],
	}
}

func getSurroundingPoints(pointState *PointState, board map[string]map[string]string, state string) []*PointState {
	leftCol := pointState.x - 1
	upRow := pointState.y - 1
	rightCol := pointState.x + 1
	downRow := pointState.y + 1

	surroundingPoints := []*PointState{}

	rstr := strconv.Itoa(pointState.y)
	cstr := strconv.Itoa(leftCol)
	if leftCol >= 0 && board[rstr][cstr] == state {
		pointLeft := &PointState{
			key:   rstr + ":" + cstr,
			x:     pointState.y,
			y:     leftCol,
			state: board[rstr][cstr],
		}
		surroundingPoints = append(surroundingPoints, pointLeft)
	}
	rstr = strconv.Itoa(upRow)
	cstr = strconv.Itoa(pointState.x)
	if upRow >= 0 && board[rstr][cstr] == state {
		pointUp := &PointState{
			key:   rstr + ":" + cstr,
			x:     upRow,
			y:     pointState.x,
			state: board[rstr][cstr],
		}
		surroundingPoints = append(surroundingPoints, pointUp)
	}
	rstr = strconv.Itoa(pointState.y)
	cstr = strconv.Itoa(rightCol)
	if rightCol <= 18 && board[rstr][cstr] == state {
		pointRight := &PointState{
			key:   rstr + ":" + cstr,
			x:     pointState.y,
			y:     rightCol,
			state: board[rstr][cstr],
		}
		surroundingPoints = append(surroundingPoints, pointRight)
	}
	rstr = strconv.Itoa(downRow)
	cstr = strconv.Itoa(pointState.x)
	if downRow <= 18 && board[rstr][cstr] == state {
		pointDown := &PointState{
			key:   rstr + ":" + cstr,
			x:     downRow,
			y:     pointState.x,
			state: board[rstr][cstr],
		}
		surroundingPoints = append(surroundingPoints, pointDown)
	}

	return surroundingPoints
}

func removeWholeThing(pointsMap map[string]bool, board map[string]map[string]string) map[string]map[string]string {
	for point, _ := range pointsMap {
		coordinates := strings.Split(point, ":")
		board[coordinates[0]][coordinates[1]] = "e"
	}
	return board
}

func wholeThinghasNoLiberties(sp *PointState, board map[string]map[string]string) (bool, map[string]bool) {
	checkedPoints := map[string]bool{}
	hasNoLiberties := true

	pointsOfThisColor := append([]*PointState{sp}, getSurroundingPoints(sp, board, sp.state)...)

	for len(pointsOfThisColor) > 0 {
		poc, pointsOfThisColor := pointsOfThisColor[0], pointsOfThisColor[:1]
		if !checkedPoints[poc.key] {
			emptyPoints := getSurroundingPoints(poc, board, "e")
			if len(emptyPoints) > 0 {
				hasNoLiberties = false
				break
			}
			checkedPoints[poc.key] = true
			pointsOfThisColor = append(pointsOfThisColor, getSurroundingPoints(poc, board, sp.state)...)
		}

	}
	return hasNoLiberties, checkedPoints
}

func oppositeOf(playerColor string) string {
	if playerColor == "b" {
		return "w"
	} else {
		return "b"
	}
}

func (g *GameState) save() error {
	filename := "db/" + g.Id + ".json"
	v, err := json.Marshal(g)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return ioutil.WriteFile(filename, v, 0600)
}

func loadGame(id string) (*GameState, error) {
	filename := "db/" + id + ".json"
	body, err := ioutil.ReadFile(filename)
	// var new bool
	if err != nil {
		log.Println("defaulting to only game")
		body, _ = ioutil.ReadFile("db/theonlygame.json")
	}
	var g *GameState
	err = json.Unmarshal(body, &g)
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
		var turn Turn
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
		game, err := calcGame(&turnmsg) // screw game msg for now
		game.save()
		if err != nil {
			log.Fatal("invalid move submitted")
		}
		log.Printf("sending message %v", game)
		for client := range clients {
			err := client.WriteJSON(game)
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
	log.Printf("hello here is url %v and id %s", r.URL, id)
	switch id {
	case "":
		serveNewGame(w)
	default:
		serveGame(w, id)
	}
}

func serveNewGame(w http.ResponseWriter) {
	gb, err := ioutil.ReadFile("db/theonlygame.json")
	if err != nil {
		log.Fatal(err)
	}
	var g *GameState
	if err := json.Unmarshal(gb, &g); err != nil {
		log.Printf("goddammit")
		log.Fatal(err)
	}
	log.Printf("loading new game %v", g)
	gRes, _ := json.Marshal(&g)
	w.WriteHeader(201)
	w.Write(gRes)
	return
}

func serveGame(w http.ResponseWriter, id string) {
	log.Print("not serving id game for now")
}

func main() {
	log.SetFlags(0)
	log.Printf("hello am running :3")
	http.HandleFunc("/", makeHandler(mainHandler))

	http.HandleFunc("/ws", handleConnections)
	go handleMessages()

	log.Fatal(http.ListenAndServe(":4000", nil))
}
