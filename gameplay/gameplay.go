package gp

import (
	"log"
	"strconv"
	"strings"
)

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
