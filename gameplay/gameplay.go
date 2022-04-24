package gameplay

import (
	"log"
	"strconv"
	"strings"
)

func HandleStonePlay(point string, color string, board map[string]map[string]string) (bool, map[string]map[string]string) {
	valid := false
	pointState := makePointState(point, board)
	liberties := getSurroundingPoints(pointState, board, "e")
	if len(liberties) > 0 {
		valid = true
	}
	for _, sp := range getSurroundingPoints(pointState, board, oppositeOf(color)) {
		hasNoLiberties, checkedPoints := wholeThinghasNoLiberties(sp, board)
		if hasNoLiberties {
			valid = true
			removeWholeThing(checkedPoints, board)
		}
	}
	return valid, board
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
