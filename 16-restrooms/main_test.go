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
2 2
1 2
0 1
1 1
1
1 2
3 1
0
2 1
0
10 2
2 1
4
3 8
8 4 5 2 0 1 7 6
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
			[]Case{{
				NFloors: 2,
				NGroups: 2,
				Groups:  []Group{{
					NPeople:           1,
					NFloorsAccessible: 2,
					FloorsAccessible:  []int{0, 1},
				},{
					NPeople:           1,
					NFloorsAccessible: 1,
					FloorsAccessible:  []int{1},
				}},
			},{
				NFloors: 1,
				NGroups: 2,
				Groups:  []Group{{
					NPeople:           3,
					NFloorsAccessible: 1,
					FloorsAccessible:  []int{0},
				},{
					NPeople:           2,
					NFloorsAccessible: 1,
					FloorsAccessible:  []int{0},
				}},
			},{
				NFloors: 10,
				NGroups: 2,
				Groups:  []Group{{
					NPeople:           2,
					NFloorsAccessible: 1,
					FloorsAccessible:  []int{4},
				},{
					NPeople:           3,
					NFloorsAccessible: 8,
					FloorsAccessible:  []int{8, 4, 5, 2, 0, 1, 7 , 6},
				}},
			}},
		}, *output)
}
