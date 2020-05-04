package main

import (
	"bytes"
	"github.com/furstenheim/challenge_encoding"
	"github.com/stretchr/testify/assert"
	"log"
	"math"
	"testing"
)

func TestInput (t *testing.T) {
	input := `3
20
76
160
`
	output := &Challenge{}
	err := challenge_encoding.Unmarshall(output, bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatal(err)
	}
	log.Println(output)
	assert.Equal(t,
		Challenge{
			3,
			[]Case{{NRolls: 20}, {NRolls: 76}, {NRolls: 160}},
		}, *output)
}

func TestFindLowerBound (t *testing.T) {
	tcs := []struct {
		target   uint64
		method   func (uint64 uint64) uint64
		expected LowerBoundResult
	} {
		{
			51,
			func (i uint64) uint64 {return 2 * i},
			LowerBoundResult{
				result:     25,
				isPossible: true,
			},
		},
		{
			1 << 62,
			func (i uint64) uint64 {return 2 * i},
			LowerBoundResult{
				result:  1<<61,
				isPossible: true,
			},
		},
		{
			40,
			func (i uint64) uint64 {return 2 * i},
			LowerBoundResult{
				result:  20,
				isPossible: true,
			},
		},
	}
	for _, tc := range tcs {
		bound := findLowerBound(tc.target, tc.method)
		assert.Equal(t, tc.expected, bound)
	}
}

func TestGenericFormula (t *testing.T) {
	tcs := []struct{
		m,n,s uint64
		expected uint64
	} {
		{3, 2, 1, 76},
		{1, 1, 1, 43},
		{1, 1, 2, 156},
		{1, 1, 1 << 31, math.MaxUint64},
	}
	for _, tc := range tcs {
		size := genericCastle(tc.m, tc.n, tc.s)
		log.Println(int(size))
		assert.Equal(t, tc.expected, size)
	}
}

