package main

import (
	"bufio"
	"fmt"
	"github.com/furstenheim/challenge_encoding"
	"log"
	"os"
	"strconv"
	"strings"
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
		parsedCase := parseCase(c)
		result := solveCase(parsedCase)
		logCase(result, i)
	}
}

func logCase (result, i int) {
	fmt.Printf("Case #%d: %d\n", i + 1, result)

}
type Challenge struct {
	NCases int    `index:"0"`
	Cases  []Case `index:"1" indexed:"NCases"`
}
type Case struct {
	Problem string `index:"0"`
}

type ParsedCase struct {
	target int
	forbidden map[int]bool
}

func solveCase (c ParsedCase) int {
	combinations := make([]int, c.target)
	for i := 1; i < c.target; i++ {
		if c.forbidden[i] {
			continue
		}
		aux := make([]int, c.target)
		for j := 1; j <= c.target; j++ {
			if j % i == 0 {
				aux[j - 1]++
				for k, v := range(combinations) {
					if k + j >= len(combinations) {
						break
					}
					aux[k + j]+= v
				}
			}
		}
		for j, v := range(aux) {
			combinations[j] += v
		}
	}
	log.Println(combinations)
	return combinations[len(combinations) - 1]
}

func parseCase (c Case) ParsedCase {
	fields := strings.Fields(c.Problem)
	if len(fields) == 0 {
		log.Fatal("Wrong length")
	}
	target, err := strconv.Atoi(fields[0])
	if err != nil {
		log.Fatal(err)
	}
	forbidden := map[int]bool{}

	for i := 1; i < len(fields); i++ {
		forbid, err := strconv.Atoi(fields[i])
		if err != nil {
			log.Fatal(err)
		}
		forbidden[forbid] = true
	}
	return ParsedCase{target: target, forbidden: forbidden}
}

