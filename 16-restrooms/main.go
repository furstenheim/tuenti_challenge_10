package main

import (
	"bufio"
	"fmt"
	"github.com/furstenheim/challenge_encoding"
	"github.com/yourbasic/graph"
	"log"
	"os"
	"sort"
)

func main () {
	reader := bufio.NewReader(os.Stdin)
	challenge := &Challenge{}
	err := challenge_encoding.Unmarshall(challenge, reader)
	if err != nil {
		log.Fatal(err)
	}
	for i, c := range(challenge.Cases) {
		log.Println("Case ", i)
		result := solveCase(c)
		logCase(result, i)
	}
}
func logCase (result int, i int) {
	fmt.Printf("Case #%d: %d\n", i + 1, result)
}

func solveCase (c Case) int {
	nPeople := int64(0)
	for _, g := range c.Groups {
		nPeople += g.NPeople
	}
	nBathrooms := sort.Search(int(nPeople), func (i int) bool {
		return nPeople == maxPeopleInRestroom(c, int64(i + 1))
	})
	return nBathrooms + 1
}

func findLowerBound (target uint64, method func (i uint64) uint64) LowerBoundResult {
	if method(0) > target {
		return LowerBoundResult{
			result:     0,
			isPossible: false,
		}
	}
	var lowerBound uint64
	for method(lowerBound + 1) <= target {
		lowerBound += 1
		var delta uint64 = 1
		for method(lowerBound + delta) <= target {
			lowerBound += delta
			delta = 2 * delta
		}
	}
	return LowerBoundResult{
		result: lowerBound,
		isPossible: true,
	}
}

type LowerBoundResult struct {
	result uint64
	isPossible bool
}

func maxPeopleInRestroom (c Case, capacity int64) int64	{
	g := graph.New(c.NFloors + c.NGroups + 2)
	for i, group := range (c.Groups) {
		g.AddCost(0, i + 1, group.NPeople)
		for _, v := range (group.FloorsAccessible) {
			g.AddCost(i + 1, 1 + c.NGroups + v, group.NPeople)
		}
	}
	for i := 0; i < c.NFloors; i++ {
		g.AddCost(1 + c.NGroups + i, 1 + c.NGroups + c.NFloors, capacity)
	}
	flow, _ := graph.MaxFlow(g, 0, 1 + c.NGroups + c.NFloors)
	return flow
}

type Challenge struct {
	NCases int    `index:"0"`
	Cases  []Case `index:"1" indexed:"NCases"`
}
type Case struct {
	NFloors   int       `index:"0" delimiter:"space" `
	NGroups int          `index:"1"`
	Groups  []Group `index:"2" indexed:"NGroups"`
}

type Group struct {
	NPeople int64 `index:"0" delimiter:"space" `
	NFloorsAccessible int64 `index:"1"`
	FloorsAccessible []int `index:"2" indexed:"NFloorsAccessible" elem_delimiter:"space"`
}


