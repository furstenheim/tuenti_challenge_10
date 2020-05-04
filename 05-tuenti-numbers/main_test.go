package main

import (
	"bytes"
	"github.com/furstenheim/challenge_encoding"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestInput (t *testing.T) {
	input := `4
20
80
35
4611686018427388000
`
	output := &Challenge{}
	err := challenge_encoding.Unmarshall(output, bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatal(err)
	}
	log.Println(output)
	assert.Equal(t,
	Challenge{
		4,
				[]Case{{20}, {80}, {35}, {4611686018427388000}},
	}, *output)
}

