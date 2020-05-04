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
	err := challenge_encoding.Unmarshall(challenge, reader)
	if err != nil {
		log.Fatal(err)
	}
	for i, c := range(challenge.Cases) {
		log.Println("Case ", i, c)
		result := solveCase(c)
		logResult(i, result)
	}
}

func logResult (nCase, result int) {
	fmt.Printf("Case #%d: %d\n", nCase + 1, result)
}

func solveCase (c Case) int {
	losers := map[int]bool{}
	maxPlayer := 1
	for _, g := range(c.Games) {
		p1, p2 := g.First, g.Second
		if g.Result == 1 {
			losers[p2] = true
		} else {
			losers[p1] = true
		}
		maxPlayer = maxInt(p1, maxPlayer)
		maxPlayer = maxInt(p2, maxPlayer)
	}
	winners := []int{}
	for i := 1; i <= maxPlayer; i++ {
		_, ok := losers[i]
		if !ok {
			winners = append(winners, i)
		}
	}
	if len(winners) != 1 {
		log.Fatal("Winners does not have expected length", winners)
	}
	return winners[0]
}


func maxInt (a, b int) int {
	if a > b {
		return a
	}
	return b
}

type Challenge struct {
	NCases int    `index:"0"`
	Cases  []Case `index:"1" indexed:"NCases"`
}
type Case struct {
	NGames int `index:"0"`
	Games []Game `index:"1" indexed:"NGames"`
}
type Result int

type Game struct {
	First int `index:"0" delimiter:"space"`
	Second int `index:"1" delimiter:"space"`
	Result Result `index:"2"`
}

