package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/furstenheim/challenge_encoding"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func main () {
	file, err := ioutil.ReadFile("pg17013.txt")
	if err != nil {
		log.Fatal(err)
	}
	book := parseBook(string(file))
	frequencyMap := parseFrequencyMap(book)
	dictionary := parseFrequencyDictionary(frequencyMap)

	reader := bufio.NewReader(os.Stdin)
	challenge := &Challenge{}
	err = challenge_encoding.Unmarshall(challenge, reader)
	if err != nil {
		log.Fatal(err)
	}
	for i, c := range(challenge.Cases) {
		log.Println("Case ", i, c)
		solveCase(i, c, dictionary)
	}
}

var numberRegex = regexp.MustCompile(`\d+`)
func solveCase (i int, c Case, dictionary FrequencyDictionary) {
	if (numberRegex.MatchString(c.Query)) {
		integer, err := strconv.Atoi(c.Query)
		if err != nil {
			log.Fatal(err)
		}
		result := dictionary.ByPosition[integer-1]
		fmt.Printf("Case #%d: %s %d\n", i +1, result.Word, result.Count)
	} else {
		result, ok := dictionary.ByWord[c.Query]
		if !ok {
			log.Fatal("Unknown word", c.Query)
		}
		fmt.Printf("Case #%d: %d #%d\n", i + 1, result.Count, result.Position + 1)
	}
}



type Challenge struct {
	NCases int    `index:"0"`
	Cases  []Case `index:"1" indexed:"NCases"`
}
type Case struct {
	Query string `index:"0"`
}
func parseBook (in string) string {
	lowerString := strings.ToLower(in)
	re := regexp.MustCompile(`[^abcdefghijklmnñopqrstuvwxyzáéíóúü]`)
	spacedString := re.ReplaceAll([]byte(lowerString), []byte(" "))
	return string(spacedString)
}

func parseFrequencyMap (in string) map[string]int {
	inWithSpace := in + " "
	reader := bufio.NewReader(bytes.NewReader([]byte(inWithSpace)))
	frequencyMap := map[string]int{}
	for word, err := reader.ReadString(' '); err != io.EOF; word, err = reader.ReadString(' ') {
		if err != nil {
			log.Fatal(err)
		}
		word = strings.TrimSuffix(word, " ")

		if len([]rune(word)) < 3 {
			continue
		}
		frequencyMap[word]++
	}
	return frequencyMap
}

func parseFrequencyDictionary (frequencyMap map[string]int) FrequencyDictionary {
	positionResult := make([]PositionResult, 0, len(frequencyMap))
	for word, count := range(frequencyMap) {
		positionResult = append(positionResult, PositionResult{Word: word, Count: count})
	}
	sort.Slice(positionResult, func (i, j int) bool {
		if positionResult[i].Count > positionResult[j].Count {
			return true
		}
		if positionResult[i].Count < positionResult[j].Count {
			return false
		}
		return lessSpanish(positionResult[i].Word, positionResult[j].Word)
	})
	wordResult := map[string]WordResult{}
	for i, byPosition := range(positionResult) {
		wordResult[byPosition.Word] = WordResult{Count: byPosition.Count, Position: i}
	}
	return FrequencyDictionary{ByWord: wordResult, ByPosition: positionResult}
}

func lessSpanish (a, b string) bool {
	aRunes := []rune(a)
	bRunes := []rune(b)
	for i, v := range(aRunes) {
		if i >= len(bRunes) {
			return false
		}
		w := bRunes[i]
		if v == w {
			continue
		}
		return v < w
	}
	return true
}

type FrequencyDictionary struct {
	ByWord map[string]WordResult
	ByPosition []PositionResult
}

type WordResult struct {
	Position int
	Count int
}

type PositionResult struct {
	Word string
	Count int
}
