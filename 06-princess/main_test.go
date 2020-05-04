package main

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestParse (t *testing.T) {
	currentView := map[Position]PositionType{}
	parseLine("..#.#\n",currentView, -2)
	parseLine(".#..#\n",currentView, -1)
	parseLine("..KP#\n",currentView, 0)
	parseLine(".#..#\n",currentView, 1)
	parseLine("..###\n",currentView, 2)
	log.Println(currentView)
	assert.Equal(t, currentView[Position{
		X: 0,
		Y: 0,
	}], Knight)
	assert.Equal(t, currentView[Position{
		X: 1,
		Y: 0,
	}], Princess)
	assert.Equal(t, currentView[Position{
		X: 2,
		Y: 1,
	}], Invalid)
	assert.Equal(t, currentView[Position{
		X: 1,
		Y: -2,
	}], Invalid)

}
