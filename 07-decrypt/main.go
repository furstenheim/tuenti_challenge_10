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
	// fmt.Println(challenge.NCases)
	for i, c := range(challenge.Cases) {
		log.Println("Case ", i, c)
		result := solveCase(c)
		fmt.Printf("Case #%d: %s\n", i + 1, result)
	}
}

var translation =  map[byte]byte {
'a': 'a',
'b': 'n',
'c': 'i',
'd': 'h',
'e': 'd',
'f': 'y',
'g': 'u',
'h': 'j',
'i': 'g',
'j': 'c',
'k': 'v',
'l': 'p',
'm': 'm',
'n': 'l',
'o': 's',
'p': 'r',
'q': 'x',
'r': 'o',
's': ';',
't': 'k',
'u': 'f',
'v': '.',
'w': ',',
'x': 'b',
'y': 't',
'z': '/',
'.': 'e',
',': 'w',
';': 'z',
'\'': 'q',
'-': '\'',
'/': '-',
}
func solveCase (c Case) string {
	text := []byte(c.Query)
	for i, v := range(text) {
		if newV, ok := translation[v]; ok {
			text[i] = newV
		}
	}
	return string(text)
}

type Challenge struct {
	NCases int    `index:"0"`
	Cases  []Case `index:"1" indexed:"NCases"`
}
type Case struct {
	Query string `index:"0"`
}
