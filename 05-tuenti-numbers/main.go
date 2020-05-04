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
		logCase(result, i)
	}
}

func logCase (result result, i int) {
	if !result.isPossible {
		fmt.Printf("Case #%d: IMPOSSIBLE\n", i + 1)
	} else {
		fmt.Printf("Case #%d: %d\n", i + 1, result.result)
	}
}

type result struct {
	isPossible bool
	result uint64
}
func solveCase (c Case) result {
	nTwenties := c.Query / 20
	reminder := c.Query % 20
	if nTwenties * 9 < reminder {
		return result{}
	}
	return result{isPossible: true, result: nTwenties}
}

type Challenge struct {
	NCases int    `index:"0"`
	Cases  []Case `index:"1" indexed:"NCases"`
}
type Case struct {
	Query uint64 `index:"0"`
}


