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
	"runtime"
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

func calcGame(tm *Turn) (*GameState, bool) {
	prevGame, err := loadGame(tm.GameId)
	if err != nil {
		log.Fatal("couldnt load game")
	}
	var started bool
	newPlayers := prevGame.Players
	if tm.Move == Switch {
		if prevGame.Started == true {
			log.Println("invalid switch attempt")
			return prevGame, false
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
	if !tm.FinishedTurn {
		nextPlayer = tm.Player
	} else if tm.Move != Switch {
		if tm.Color == "b" {
			nextPlayer = prevGame.Players.W
		} else {
			nextPlayer = prevGame.Players.B
		}
	}

	valid := true
	newBoard := tm.BoardTemp
	if tm.Move == Play {
		valid, newBoard = HandleStonePlay(tm.Point, tm.Color, tm.BoardTemp)
	}

	if !valid {
		fmt.Println("not valid alas")
		return prevGame, false
	} else {
		return &GameState{
			Id:         prevGame.Id,
			Board:      newBoard,
			LastPlayed: tm.Point,
			NextPlayer: nextPlayer,
			Players:    newPlayers,
			Started:    started,
			Ended:      ended,
			Winner:     "",
		}, true
	}
}

func (g *GameState) save() error {
	log.Printf("ummmmm %v", g)
	if g.Id != "" {
		filename := "db/" + g.Id + ".json"
		v, err := json.Marshal(g)
		if err != nil {
			log.Fatal(err)
			return err
		}
		return ioutil.WriteFile(filename, v, 0600)
	}
	return nil
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
	if g == nil { // not sure this is right way to check
		return nil, errors.New("couldn't find game")
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
	origin := r.Header.Get("Origin")
	log.Printf("hmmmmmm %v", origin)	
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
		
		game, turnWasValid := calcGame(&turnmsg) // screw game msg for now
		if !turnWasValid {
			log.Println("invalid move submitted")
			// make sure front end state comes into alignment with prev saved game state
		} else {
			game.save()
		}
		
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
	log.Printf("httping!")
	origin := r.Header.Get("Origin")
	log.Printf("hmmmmmm %v", origin)
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
	gb, err := ioutil.ReadFile("db/newgame.json")
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
	// http.HandleFunc("/", makeHandler(mainHandler))

	http.HandleFunc("/ws", handleConnections)
	go handleMessages()

	isLocal := runtime.GOOS == "darwin"

	if (isLocal) {
		// skip connection dialog
		log.Fatal(http.ListenAndServe("localhost:4000", nil))
	} else {
		log.Fatal(http.ListenAndServe(":4000", nil))
	}
}

/***/

type PointState struct {
	key   string
	x     int
	y     int
	state string
}

func HandleStonePlay(point string, color string, board map[string]map[string]string) (bool, map[string]map[string]string) {
	valid := false
	if point == "" {
		log.Println("no point given!")
		return valid, board
	}
	pointState := makePointState(point, board)
	liberties := getSurroundingPoints(pointState, board, "e")
	ownColors := getSurroundingPoints(pointState, board, color)

	var loggablelibs []string
	for _, x := range liberties {
		loggablelibs = append(loggablelibs, x.key)
	}

	if len(ownColors) == 4 {
		valid = true
		return valid, board
	}

	log.Printf("found liberties %v", loggablelibs)
	if len(liberties) > 0 {
		valid = true
	}
	for _, sp := range getSurroundingPoints(pointState, board, oppositeOf(color)) {
		log.Println("wholeThinghasNoLiberties: iterating surrounding points of last played")
		hasNoLiberties, checkedPoints := wholeThinghasNoLiberties(sp, board)
		if hasNoLiberties {
			valid = true
			log.Println("Removing no-liberty blob")
			removeWholeThing(checkedPoints, board)
		}
	}
	if !valid {
		valid = placementIsValid(pointState, board)
	}
	log.Println("finished analyzing move")
	return valid, board
}

func makePointState(point string, board map[string]map[string]string) *PointState {
	coordinates := strings.Split(point, ":")
	row, err := strconv.Atoi(coordinates[0])
	log.Printf("point: %s, coords: %v", point, coordinates)
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

func placementIsValid(pointState *PointState, board map[string]map[string]string) bool {
	hasNoLiberties, _ := wholeThinghasNoLiberties(pointState, board)
	return !hasNoLiberties
}

func getSurroundingPoints(pointState *PointState, board map[string]map[string]string, state string) []*PointState {
	upRow := pointState.x - 1
	leftCol := pointState.y - 1
	downRow := pointState.x + 1
	rightCol := pointState.y + 1

	surroundingPoints := []*PointState{}

	rowstring := strconv.Itoa(pointState.x)
	colstring := strconv.Itoa(leftCol)
	if leftCol >= 0 && board[rowstring][colstring] == state {
		pointLeft := &PointState{
			key:   rowstring + ":" + colstring,
			x:     pointState.x,
			y:     leftCol,
			state: board[rowstring][colstring],
		}
		surroundingPoints = append(surroundingPoints, pointLeft)
	}
	rowstring = strconv.Itoa(upRow)
	colstring = strconv.Itoa(pointState.y)
	if upRow >= 0 && board[rowstring][colstring] == state {
		pointUp := &PointState{
			key:   rowstring + ":" + colstring,
			x:     upRow,
			y:     pointState.y,
			state: board[rowstring][colstring],
		}
		surroundingPoints = append(surroundingPoints, pointUp)
	}
	rowstring = strconv.Itoa(pointState.x)
	colstring = strconv.Itoa(rightCol)
	if rightCol <= 18 && board[rowstring][colstring] == state {
		pointRight := &PointState{
			key:   rowstring + ":" + colstring,
			x:     pointState.x,
			y:     rightCol,
			state: board[rowstring][colstring],
		}
		surroundingPoints = append(surroundingPoints, pointRight)
	}
	rowstring = strconv.Itoa(downRow)
	colstring = strconv.Itoa(pointState.y)
	if downRow <= 18 && board[rowstring][colstring] == state {
		pointDown := &PointState{
			key:   rowstring + ":" + colstring,
			x:     downRow,
			y:     pointState.y,
			state: board[rowstring][colstring],
		}
		surroundingPoints = append(surroundingPoints, pointDown)
	}

	var loggablePoints []string
	for _, x := range surroundingPoints {
		loggablePoints = append(loggablePoints, x.key)
	}
	log.Printf("surrounding points for state %v %v", state, loggablePoints)
	return surroundingPoints
}

func wholeThinghasNoLiberties(sp *PointState, board map[string]map[string]string) (bool, map[string]bool) {
	checkedPoints := map[string]bool{}
	hasNoLiberties := true

	pointsOfThisColor := []*PointState{}
	pointsOfThisColor = append(pointsOfThisColor, sp)
	pointsOfThisColor = append(pointsOfThisColor, getSurroundingPoints(sp, board, sp.state)...)

	log.Printf("starting recursive check with ")

	var loggablePoints []string
	for _, x := range pointsOfThisColor {
		loggablePoints = append(loggablePoints, x.key)
	}
	log.Printf("starting iterative search with %v", pointsOfThisColor)

	for len(pointsOfThisColor) > 0 {
		poc := pointsOfThisColor[0] // pop off front of queue

		var loggablePoints []string
		for _, x := range pointsOfThisColor {
			loggablePoints = append(loggablePoints, x.key)
		}
		log.Printf("before attempting pop %v", loggablePoints)

		pointsOfThisColor = pointsOfThisColor[1:]

		var loggablePoints2 []string
		for _, x := range pointsOfThisColor {
			loggablePoints2 = append(loggablePoints2, x.key)
		}
		log.Printf("afte rattempting pop, now have poc %v and remaining points %v", poc.key, loggablePoints2)

		if !checkedPoints[poc.key] {
			emptyPoints := getSurroundingPoints(poc, board, "e")
			var loggableEPoints []string
			for _, x := range emptyPoints {
				loggableEPoints = append(loggableEPoints, x.key)
			}
			log.Printf("found empty points %v", loggableEPoints)

			if len(emptyPoints) > 0 {
				hasNoLiberties = false
				return hasNoLiberties, checkedPoints
			}
			checkedPoints[poc.key] = true
			pointsOfThisColor = append(pointsOfThisColor, getSurroundingPoints(poc, board, sp.state)...)
		}

	}
	return hasNoLiberties, checkedPoints
}

func removeWholeThing(pointsMap map[string]bool, board map[string]map[string]string) map[string]map[string]string {
	for point := range pointsMap {
		coordinates := strings.Split(point, ":")
		board[coordinates[0]][coordinates[1]] = "e"
	}
	return board
}

func oppositeOf(playerColor string) string {
	if playerColor == "b" {
		return "w"
	} else {
		return "b"
	}
}
