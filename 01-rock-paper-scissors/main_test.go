package main

import (
	"bytes"
	"github.com/furstenheim/challenge_encoding"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInput (t *testing.T) {
	input := `8
R S
R P
P P
S P
R R
P R
P P
S R
`
	output := &Challenge{}
	challenge_encoding.Unmarshall(output, bytes.NewReader([]byte(input)))
	assert.Equal(t, output.NCases, 8)
	assert.Equal(t, output.Games[7], Game{
		First:  "S",
		Second: "R",
	})
}
