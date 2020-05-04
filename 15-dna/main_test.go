package main

import (
	"bytes"
	"fmt"
	"github.com/furstenheim/challenge_encoding"
	"github.com/stretchr/testify/assert"
	"hash/crc32"
	"log"
	"strconv"
	"testing"
)
func TestInput (t *testing.T) {
	input := `2
aardwolf0000 2
0 224
0 192
admiral0000 4
1 227
2 232
2 46
0 169
`
	output := &Challenge{}
	err := challenge_encoding.Unmarshall(output, bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatal(err)
	}
	log.Println(output)
	assert.Equal(t,
		Challenge{
			2,
			[]Case{{FileName: "aardwolf0000", NAdditions: 2, Additions: []SparseByte{{0, 224}, {0, 192}}}, {FileName: "admiral0000", NAdditions: 4, Additions: []SparseByte{{1, 227}, {2, 232}, {2, 46}, {0, 169}}}},
		}, *output)
}
/*
func TestParseCase (t *testing.T) {
	tcs := []struct {
		input Case
		expected []SparseFile
	} {
		{
			Case{FileName: "aardwolf0000", NAdditions: 2, Additions: []SparseByte{{0, 224}, {0, 192}}},
				[]SparseFile{
					{size: 1, fileBytes: []SparseByte{{0, 224}}},
					{size: 2, fileBytes: []SparseByte{{0, 192}, {1, 224}}},
				},
		},
		{
			Case{FileName: "admiral0000", NAdditions: 4, Additions: []SparseByte{{1, 227}, {2, 232}, {2, 46}, {0, 169}}},
				[]SparseFile{
					{size: 2, fileBytes: []SparseByte{{1, 227}}},
					{size: 3, fileBytes: []SparseByte{{1, 227}, {2, 232}}},
					{size: 4, fileBytes: []SparseByte{{1, 227}, {2, 46}, {3, 232}}},
					{size: 5, fileBytes: []SparseByte{{0, 169}, {2, 227}, {3, 46}, {4, 232}}},
				},
		},
	}

	for _, tc := range tcs {
		parsed := parseCase(tc.input)
		assert.Equal(t, tc.expected, parsed)
	}

}*/

func TestRudimentary (t *testing.T) {
	tcs := []struct{
		input []byte
	} {
		{[]byte{'a'}},
		{[]byte{0}},
		{[]byte{224}},
		{[]byte{1}},
		{[]byte{1 ^224}},
		{[]byte{224}},
		{[]byte{1}},
		{[]byte{1, 0, 0, 0}},
		{[]byte{0, 244, 0, 0}},
		{[]byte{1, 244, 0, 0}},
		{[]byte{'b','c','d','a'}},

	}

	for _, tc := range tcs {
		checksumIEEE := crc32.ChecksumIEEE(tc.input)
		rudimentary := ieeeRudimentary(tc.input)
		assert.Equal(t, checksumIEEE, rudimentary)
	}
}

func TestRudimentaryGreater4 (t *testing.T) {
	tcs := []struct{
		input []byte
	} {
		{[]byte{1, 0, 0, 0}},
		{[]byte{0, 244, 0, 0}},
		{[]byte{1, 244, 0, 0}},
		{[]byte{'b','c','d','a'}},
		{[]byte{5, 6, 7, 8, 7, 8}},
		{[]byte{5, 6, 7, 8, 7, 8, 5, 6, 7, 8, 7, 8}},

	}

	for _, tc := range tcs {
		checksumIEEE := crc32.ChecksumIEEE(tc.input)
		rudimentary := ieeeRudimentaryGreater4(tc.input)
		assert.Equal(t, checksumIEEE, rudimentary)
	}
}

func TestToSparseFile (t * testing.T) {
	tcs := []struct {
		input SparseFile
		expected []byte
	} {
		{
			SparseFile{[]SparseByte{}, 4},
			[]byte{0,0 ,0, 0},
		},		{
			SparseFile{[]SparseByte{{1, 2}}, 4},
			[]byte{0,2 ,0, 0},
		},{
			SparseFile{[]SparseByte{{1, 2}}, 5},
			[]byte{0,2 , 0, 0, 0},
		},
	}
	for _, tc := range tcs {
		res := tc.input.toBytes()
		// log.Println(res, tc.expected)
		assert.Equal(t, tc.expected, res)
	}
}

func TestCRCSparseFile (t *testing.T) {
	tcs := []struct {
		input    SparseFile
		expected []byte
	}{
		{
			SparseFile{[]SparseByte{}, 4},
			[]byte{0, 0, 0, 0},
		}, {
			SparseFile{[]SparseByte{{1, 2}}, 4},
			[]byte{0, 2, 0, 0},
		}, {
			SparseFile{[]SparseByte{{1, 2}}, 5},
			[]byte{0, 2, 0, 0, 0},
		},
	}
	for _, tc := range tcs {
		res := tc.input.toBytes()
		// log.Println(res, tc.expected)
		assert.Equal(t, tc.expected, res)
		assert.Equal(t, crc32.ChecksumIEEE(tc.expected), ieeeRudimentarySparse(tc.input))
	}
}

func TestReminderSparse (t *testing.T) {
	tcs := []struct {
		input    SparseFile
	}{
		{
			SparseFile{[]SparseByte{}, 4},
		}, {
			SparseFile{[]SparseByte{{3, 1}}, 4},
		}, {
			SparseFile{[]SparseByte{{1, 2}}, 5},
		}, {
			SparseFile{[]SparseByte{{1, 2}, {2, 3}}, 10},
		},
	}
	for _, tc := range tcs {
		res := tc.input.toBytes()
		// log.Println(res, tc.expected)
		assert.Equal(t, reminderIEEEGreater4(res),reminderIEEEGreater4Sparse(tc.input))
	}
}


func TestArrayToUint32 (t *testing.T) {
	tcs := []struct{
		input []byte
		expected uint32
	} {
		{[]byte{0,0, 0, 0}, 0},
		{[]byte{0,0, 1, 4}, 1 << 8 + 1 << 2},
		{[]byte{0, 0, 1, 0}, 1 << 8},

	}

	for _, tc := range tcs {
		res := arrayToUint32(tc.input)
		// log.Println(res, tc.expected)
		assert.Equal(t, tc.expected, res)
		assert.Equal(t, tc.input, uint32ToArray(res))
	}
}

func TestArrayToUint64 (t *testing.T) {
	tcs := []struct{
		input []byte
		expected uint64
	} {
		{[]byte{0,0, 0, 0, 0,0, 0, 0}, 0},
		{[]byte{0,0, 0, 0, 0,0, 1, 4}, 1 << 8 + 1 << 2},
		{[]byte{0,0, 0, 0, 0, 0, 1, 0}, 1 << 8},
	}

	for _, tc := range tcs {
		res := arrayToUint64(tc.input)
		// log.Println(res, tc.expected)
		assert.Equal(t, tc.expected, res)
		assert.Equal(t, tc.input, uint64ToArray(res))
	}
}

func TestMultiplicative (t *testing.T) {
	increase1 := reminderIEEEGreater4([]byte{1, 0, 0, 0, 0})
	increase2 := reminderIEEEGreater4([]byte{1, 0, 0, 0, 0})
	increase3 := reminderIEEEGreater4([]byte{1, 0, 0, 0, 0, 0, 0, 0, 0})
	assert.Equal(t, increase3, reminderIEEEGreater4(uint64ToArray(xorMultiply32(increase1, increase2))))
	log.Println(fmt.Sprintf("%8b", reminderIEEEGreater4(uint64ToArray(xorMultiply32(increase1, increase2)))))
	log.Println(fmt.Sprintf("%8b", increase3))
	log.Println(fmt.Sprintf("%8b", uint64ToArray(xorMultiply32(increase1, increase2))))
	log.Println(fmt.Sprintf("%8b", xorMultiply32(increase1, increase2)))
	log.Println(fmt.Sprintf("%8b", increase1))
}

func TestPadded (t *testing.T) {
	increase1 := reminderIEEEGreater4([]byte{0, 1, 0, 0, 0, 0})
	increase2 := reminderIEEEGreater4([]byte{0,0 , 1, 0, 0, 0, 0})
	increase3 := reminderIEEEGreater4([]byte{0, 0, 0, 1, 0, 0, 0, 0})
	assert.Equal(t, increase2, increase3)
	assert.Equal(t, increase2, increase1)
}
func TestUint64ToArray (t * testing.T) {
/*// 1010001 1010000   10100 1010001 1000001 1000100 1000101   10101
	log.Println("ddd", fmt.Sprintf("%8b", AsByteSlice("101000101010000000101000101000101000001010001000100010100010101")))
	log.Println("ddd", AsByteSlice("101000101010000000101000101000101000001010001000100010100010101"))*/
	x, _ := strconv.ParseUint("101000101010000000101000101000101000001010001000100010100010101", 2, 64)
	tcs := []struct{
		expected []byte
		input uint64
	} {
		{[]byte{0,0, 0, 0, 0,0, 0, 0}, 0},
		{[]byte{0,0, 0, 0, 0,0, 1, 4}, 1 << 8 + 1 << 2},
		{[]byte{81, 80, 20, 81, 65, 68, 69, 21}, x},
	}

	for _, tc := range tcs {
		res := uint64ToArray(tc.input)
		// log.Println(res, tc.expected)
		assert.Equal(t, tc.expected, res)
	}
}

func AsByteSlice(b string) []byte {
	var out []byte
	var str string

	for i := len(b); i > 0; i -= 8 {
		if i-8 < 0 {
			str = string(b[0:i])
		} else {
			str = string(b[i-8 : i])
		}
		v, err := strconv.ParseUint(str, 2, 8)
		if err != nil {
			panic(err)
		}
		out = append([]byte{byte(v)}, out...)
	}
	return out
}


func TestXORMultiply (t *testing.T) {
	tcs := []struct{
		input1, input2 uint32
		expected       uint64
	} {
		{1 << 3 + 1 << 2 + 1<<1, 1<< 3 + 1 << 2 + 1, 1 << 6 + 1<<2 + 1 << 1},
		{1 << 5, 1<<7, 1 << 12},
	}

	for _, tc := range tcs {
		res := xorMultiply32(tc.input1, tc.input2)
		assert.Equal(t, tc.expected, res)
	}
}


