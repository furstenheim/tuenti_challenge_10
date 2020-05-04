package main

import (
	"bytes"
	"github.com/furstenheim/challenge_encoding"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestInput (t *testing.T) {
	input := `1
3
3 2 1
1 3 0
3 2 1
`
	output := &Challenge{}
	err := challenge_encoding.Unmarshall(output, bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatal(err)
	}
	log.Println(output)
	assert.Equal(t, 1, output.NCases)
	assert.Equal(t, Case{
		NGames: 3,
		Games: []Game{
			{3, 2, 1},
			{1, 3, 0},
			{3, 2, 1},
		},
	}, output.Cases[0])
}

