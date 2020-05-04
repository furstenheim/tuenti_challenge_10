package main

import (
	"bufio"
	"log"
	"net"
)

func main () {
	conn, err := net.Dial("tcp", "52.49.91.111:2003")
	if err != nil {
		log.Fatal(err)
	}
	currentView := map[Position]PositionType{}
	globalMap := map[Position]int{}
	globalPosition := Position{X: 0, Y: 0}

	reader := bufio.NewReader(conn)
	previousEmpty := false
	currentLine := -2
	currentBound := 1000
	currentMovements := []Movement{}
	var lastMovement Movement

	counter := 0
	for true {
		counter++
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		if line == "--- Oops! Invalid command ---\n" {
			previousEmpty = true // when invalid there is only one line
			continue
		}
		if line[0] == 'P' || line[0] == '.' || line[0] == '#' {
			parseLine(line, currentView, currentLine)
			currentLine++
			continue
		}

		if line == "\n" && !previousEmpty {
			previousEmpty = true
			continue
		}
		if (line == "--- Quick there is no time to lose! The princess is in danger ---\n" || line == "\n") {
			previousEmpty = false

			var nextCommand string
			if len(currentMovements) == currentBound {
				log.Println("reached limit. current position", globalPosition)
				currentMovements, lastMovement = currentMovements[0: len(currentMovements) - 1], currentMovements[len(currentMovements) - 1]
				nextCommand = lastMovement.reverse

				globalPosition.X -= lastMovement.dX
				globalPosition.Y -= lastMovement.dY

			}

			// A* here. We have different movements for each quadrant because the princess is at the origin
			var POSSIBLE_MOVEMENTS []Movement
			if globalPosition.X >= 0 && globalPosition.Y >= 0 {
				POSSIBLE_MOVEMENTS = POSSIBLE_MOVEMENTS_1
			} else if globalPosition.X >= 0 && globalPosition.Y < 0 {
				POSSIBLE_MOVEMENTS = POSSIBLE_MOVEMENTS_2
			} else if globalPosition.Y < 0 {
				POSSIBLE_MOVEMENTS = POSSIBLE_MOVEMENTS_3
			} else {
				POSSIBLE_MOVEMENTS = POSSIBLE_MOVEMENTS_4
			}
			if nextCommand == "" {
				for i := 0; i < len(POSSIBLE_MOVEMENTS); i++ {
					possibleMovement := POSSIBLE_MOVEMENTS[i]
					possiblePosition := Position{
						X: possibleMovement.dX,
						Y: possibleMovement.dY,
					}
					nextGlobalPosition := Position{X: globalPosition.X + possibleMovement.dX, Y: globalPosition.Y + possibleMovement.dY}

					if currentView[possiblePosition] == 0 {
						log.Fatal("Unknown position", possiblePosition)
					}
					if currentView[possiblePosition] != Invalid && (globalMap[nextGlobalPosition] == 0 || globalMap[nextGlobalPosition] > len(currentMovements) + 1){
						if (counter % 300 == 0) {
							log.Println("Moving forward", globalPosition, currentBound, counter, len(currentMovements))
						}
						nextCommand = possibleMovement.name
						currentMovements = append(currentMovements, possibleMovement)
						// startMovementIndex = 0
						globalPosition = nextGlobalPosition
						globalMap[globalPosition] = len(currentMovements)
						break
					}
				}
			}

			if nextCommand == ""  {
				if len(currentMovements) == 0 {
					currentBound++
					nextCommand = "ODOL" // Force refresh
					log.Println("global map", globalMap)
					if globalMap[Position{1, 0}] != 0 {
						log.Fatal("Princess was found", currentBound)
					}
					globalMap = map[Position]int{}

				} else {
					currentMovements, lastMovement = currentMovements[0: len(currentMovements) - 1], currentMovements[len(currentMovements) - 1]
					nextCommand = lastMovement.reverse
					globalPosition.X -= lastMovement.dX
					globalPosition.Y -= lastMovement.dY
				}
			}


			_, err := conn.Write([]byte(nextCommand + "\n"))
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal("Unknown line", line, currentBound)
		}
		currentView = map[Position]PositionType{}
		currentLine = -2
	}
}



func parseLine (line string, currentView map[Position]PositionType, currentLine int) {
	if len(line) != 6 {
		log.Fatal(line, len(line))
	}
	for i, v := range line {
		position := Position{
			X: i - 2,
			Y: -currentLine,
		}
		if v == '#' {
			currentView[position] = Invalid
		}
		if v == '.' {
			currentView[position] = Valid
		}
		if v == 'P' {
			currentView[position] = Princess
		}
		if v == 'K' {
			currentView[position] = Knight
		}
	}
}
// Simple node script to generate all the possible movements
// > a.forEach(function (b, i) {b.reverse = Math.abs(b.dY) + (b.dY > 0 ? 'D' : 'U') + Math.abs(b.dX) + (b.dX > 0 ? 'L' : 'R'); b.index = i })
//undefined
//> a.forEach(function (b) {b.name = Math.abs(b.dY) + (b.dY < 0 ? 'D' : 'U') + Math.abs(b.dX) + (b.dX < 0 ? 'L' : 'R') })
/*var POSSIBLE_MOVEMENTS = []Movement{
	{ dX: -2, dY: 1, name: "1U2L", reverse: "1D2R", index: 0 },
	{ dX: -2, dY: -1, name: "1D2L", reverse: "1U2R", index: 1 },
	{ dX: -1, dY: -2, name: "2D1L", reverse: "2U1R", index: 2 },
	{ dX: -1, dY: 2, name: "2U1L", reverse: "2D1R", index: 3 },
	{ dX: 1, dY: 2, name: "2U1R", reverse: "2D1L", index: 4 },
	{ dX: 1, dY: -2, name: "2D1R", reverse: "2U1L", index: 5 },
	{ dX: 2, dY: -1, name: "1D2R", reverse: "1U2L", index: 6 },
	{ dX: 2, dY: 1, name: "1U2R", reverse: "1D2L", index: 7 },
}*/

var POSSIBLE_MOVEMENTS_1 = []Movement{
	{ dX: -2, dY: -1, name: "1D2L", reverse: "1U2R", index: 1 },
	{ dX: -1, dY: -2, name: "2D1L", reverse: "2U1R", index: 2 },
	{ dX: -2, dY: 1, name: "1U2L", reverse: "1D2R", index: 0 },
	{ dX: -1, dY: 2, name: "2U1L", reverse: "2D1R", index: 3 },
	{ dX: 1, dY: -2, name: "2D1R", reverse: "2U1L", index: 5 },
	{ dX: 2, dY: -1, name: "1D2R", reverse: "1U2L", index: 6 },
	{ dX: 1, dY: 2, name: "2U1R", reverse: "2D1L", index: 4 },
	{ dX: 2, dY: 1, name: "1U2R", reverse: "1D2L", index: 7 },
}

var POSSIBLE_MOVEMENTS_2 = []Movement{
	{ dX: -1, dY: 2, name: "2U1L", reverse: "2D1R", index: 3 },
	{ dX: -2, dY: 1, name: "1U2L", reverse: "1D2R", index: 0 },
	{ dX: 2, dY: 1, name: "1U2R", reverse: "1D2L", index: 7 },
	{ dX: 1, dY: 2, name: "2U1R", reverse: "2D1L", index: 4 },
	{ dX: -1, dY: -2, name: "2D1L", reverse: "2U1R", index: 2 },
	{ dX: -2, dY: -1, name: "1D2L", reverse: "1U2R", index: 1 },
	{ dX: 2, dY: -1, name: "1D2R", reverse: "1U2L", index: 6 },
	{ dX: 1, dY: -2, name: "2D1R", reverse: "2U1L", index: 5 },
}
var POSSIBLE_MOVEMENTS_3 = []Movement{
	{ dX: 1, dY: 2, name: "2U1R", reverse: "2D1L", index: 4 },
	{ dX: 2, dY: 1, name: "1U2R", reverse: "1D2L", index: 7 },
	{ dX: -1, dY: 2, name: "2U1L", reverse: "2D1R", index: 3 },
	{ dX: 2, dY: -1, name: "1D2R", reverse: "1U2L", index: 6 },
	{ dX: -2, dY: 1, name: "1U2L", reverse: "1D2R", index: 0 },
	{ dX: 1, dY: -2, name: "2D1R", reverse: "2U1L", index: 5 },
	{ dX: -2, dY: -1, name: "1D2L", reverse: "1U2R", index: 1 },
	{ dX: -1, dY: -2, name: "2D1L", reverse: "2U1R", index: 2 },
}

var POSSIBLE_MOVEMENTS_4 = []Movement{
	{ dX: 2, dY: -1, name: "1D2R", reverse: "1U2L", index: 6 },
	{ dX: 1, dY: -2, name: "2D1R", reverse: "2U1L", index: 5 },
	{ dX: 1, dY: 2, name: "2U1R", reverse: "2D1L", index: 4 },
	{ dX: 2, dY: 1, name: "1U2R", reverse: "1D2L", index: 7 },
	{ dX: -2, dY: -1, name: "1D2L", reverse: "1U2R", index: 1 },
	{ dX: -1, dY: -2, name: "2D1L", reverse: "2U1R", index: 2 },
	{ dX: -2, dY: 1, name: "1U2L", reverse: "1D2R", index: 0 },
	{ dX: -1, dY: 2, name: "2U1L", reverse: "2D1R", index: 3 },
}

type Movement struct {
	dX, dY int
	name string
	reverse string
	index int
}

type PositionType int

const (
	Princess PositionType = iota + 1
	Knight
	Valid
	Invalid
)

type Position struct {
	X, Y int
}
