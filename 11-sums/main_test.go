package main

import (
	"bytes"
	"github.com/furstenheim/challenge_encoding"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestInput (t *testing.T) {
	input := `3
2 1
5 4 2
4
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
			[]Case{{Problem:"2 1"}, {Problem:"5 4 2"}, {Problem:"4"}},
		}, *output)
}

