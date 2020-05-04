package main

import (
	"bufio"
	"fmt"
	"github.com/furstenheim/challenge_encoding"
	"log"
	"os"
)

func main () {
	reader := bufio.NewReader(os.Stdin)
	challenge := &Challenge{}
	challenge_encoding.Unmarshall(challenge, reader)
	log.Println(challenge.Games[1])
	for i, g := range(challenge.Games) {
		result := solveGame(g)
		logResult(i, result)
	}
}

func logResult(i int, result string) {
	fmt.Printf("Case #%d: %s\n", i + 1, result)
}

var results = map[Game]string{
	Game{"R", "R"}: "-",
	Game{"R", "P"}: "P",
	Game{"R", "S"}: "R",
	Game{"P", "R"}: "P",
	Game{"P", "P"}: "-",
	Game{"P", "S"}: "S",
	Game{"S", "R"}: "R",
	Game{"S", "P"}: "S",
	Game{"S", "S"}: "-",
}

func solveGame (game Game) string {
	result, ok := results[game]
	if !ok {
		log.Fatal("Unknown game", game)
	}
	return result
}

type Challenge struct {
	NCases int   `index:"0"`
	Games []Game `index:"1" indexed:"NCases"`
}
type Game struct {
	First string `index:"0" delimiter:"space"`
	Second string `index:"1"`
}
