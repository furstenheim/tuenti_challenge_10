package main

import (
	"bufio"
	"fmt"
	"github.com/furstenheim/challenge_encoding"
	"log"
	"math"
	"math/big"
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
func logCase (result Result, i int) {
	if !result.isPossible {
		fmt.Printf("Case #%d: IMPOSSIBLE\n", i + 1)
	} else {
		fmt.Printf("Case #%d: %d %d\n", i + 1, result.height + 2, result.nRolls)
	}
}
func solveCase (c Case) Result {
	heightResult := findLowerBound(c.NRolls, func (i uint64) uint64 {
		return genericCastle(1, 1, i + 1)
	})
	if !heightResult.isPossible {
		return Result{}
	}
	height := heightResult.result + 1
	squareCastle := findLowerBound(c.NRolls, func (i uint64) uint64 {
		return genericCastle(i, i, height)
	})
	rectangularCastle := findLowerBound(c.NRolls, func (i uint64) uint64 {
		return genericCastle(i, i + 1, height)
	})
	// log.Println(rectangularCastle)
	// log.Println(squareCastle)
	if !rectangularCastle.isPossible {
		return Result{
			isPossible: true,
			height:     height,
			nRolls:     genericCastle(squareCastle.result, squareCastle.result, height),
		} // possible since 1,1 square is always possible is possible
	}
	squareRolls := genericCastle(squareCastle.result, squareCastle.result, height)
	rectangularRolls := genericCastle(rectangularCastle.result, rectangularCastle.result + 1, height)
	// log.Println(squareRolls, rectangularRolls, height)
	if squareRolls > rectangularRolls {
		return Result{
			isPossible: true,
			height:     height,
			nRolls:     squareRolls,
		}
	} else {
		return Result{
			isPossible: true,
			height:     height,
			nRolls:     rectangularRolls,
		}
	}
}


// 2 \cdot n \cdot m + \sum_{i = 1}^{s - 2} 2(n + 4\cdot i) (m + 4\cdot i) - (n + 4 \cdot i - 2) (m + 4 \cdot i - 2)
// 2 n m + smn + 4sm + 4sn + 4s + 2s^2n + 2 s^2 m + 16 s^2 + \frac{8}{3}(2s^3 + s)
func genericCastle (m, n, s uint64) uint64 {
	result := big.NewInt(0)
	aux := big.NewInt(0)
	aux2 := big.NewInt(0)
	mb := big.NewInt(0).SetUint64(m)
	nb := big.NewInt(0).SetUint64(n)
	sb := big.NewInt(0).SetUint64(s)
	// 2mn
	aux.SetUint64(2)
	aux.Mul(aux, mb)
	aux.Mul(aux, nb)
	result.Add(result, aux)

	// smn
	aux.Set(sb)
	aux.Mul(aux, mb)
	aux.Mul(aux, nb)
	result.Add(result, aux)

	// 4sm
	aux.SetInt64(4)
	aux.Mul(aux, sb)
	aux.Mul(aux, mb)
	result.Add(result, aux)

	// 4sn
	aux.SetInt64(4)
	aux.Mul(aux, sb)
	aux.Mul(aux, nb)
	result.Add(result, aux)

	// 4s
	aux.SetInt64(4)
	aux.Mul(aux, sb)
	result.Add(result, aux)

	// 2s2n
	aux.SetInt64(2)
	aux.Mul(aux, sb)
	aux.Mul(aux, sb)
	aux.Mul(aux, nb)
	result.Add(result, aux)

	// 2s2m
	aux.SetInt64(2)
	aux.Mul(aux, sb)
	aux.Mul(aux, sb)
	aux.Mul(aux, mb)
	result.Add(result, aux)

	// 16s2
	aux.SetInt64(16)
	aux.Mul(aux, sb)
	aux.Mul(aux, sb)
	result.Add(result, aux)

	// 8/3(2s3 + s)
	aux.SetInt64(2)
	aux.Mul(aux, sb)
	aux.Mul(aux, sb)
	aux.Mul(aux, sb)
	aux.Add(aux, sb)
	aux2.SetInt64(8)
	aux.Mul(aux, aux2)
	aux2.SetInt64(3)
	aux.Div(aux, aux2)
	result.Add(result, aux)

	aux.SetUint64(math.MaxUint64)
	cmp := result.Cmp(aux)
	if cmp > 0 {
		return math.MaxUint64
	}
	return result.Uint64()
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
type Result struct {
	isPossible bool
	height     uint64
	nRolls     uint64
}

type LowerBoundResult struct {
	result uint64
	isPossible bool
}


type Challenge struct {
	NCases int    `index:"0"`
	Cases  []Case `index:"1" indexed:"NCases"`
}
type Case struct {
	NRolls uint64 `index:"0"`
}
