package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/furstenheim/challenge_encoding"
	"log"
	"math/bits"
	"os"
	"sort"
)

const (
	// IEEE is by far and away the most common CRC-32 polynomial.
	// Used by ethernet (IEEE 802.3), v.42, fddi, gzip, zip, png, ...
	IEEE = 0xedb88320
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
		solveCase(c)
	}
}
func solveCase (c Case) {
	files := parseCase(c)
	for i, file := range(files) {
		crc := ieeeRudimentarySparse(file)
		fmt.Printf("%s %d: %08x\n", c.FileName, i, crc)
	}
}

func parseCase (c Case) []SparseFile {
	fi, err := os.Stat("animals/" + c.FileName);
	if err != nil {
		log.Fatal(err)
	}
	size := fi.Size()
	subCases := []SparseFile{
		{
			fileBytes: []SparseByte{},
			size:      uint64(size),
		},
	}
	for i, _ := range (c.Additions) {
		additions := []SparseByte{}
		subCase := SparseFile{
			fileBytes: nil,
			size:      uint64(size) + uint64(i + 1) ,
		}
		for _, v := range(c.Additions[0: i + 1]) {
			modIndex := sort.Search(len(additions), func(i int) bool {
				return additions[i].Position >= v.Position
			})
			temp := append([]SparseByte{v}, additions[modIndex: len(additions)]...)
			for i := 1; i < len(temp); i++ {
				temp[i].Position += 1
			}
			additions = append(additions[0: modIndex], temp...)
			sort.Slice(additions, func (i, j int) bool {
				return additions[i].Position < additions[j].Position
			})
		}
		subCase.fileBytes = additions
		subCases = append(subCases, subCase)
	}
	return subCases
}


func ieeeRudimentary (array []byte) uint32 {
	for i, v := range(array) {
		array[i] = bits.Reverse8(v)
	}
	newReminder := reminderIEEE(array)
	return newReminder ^ 0xFFFFFFFF
}

func ieeeRudimentaryGreater4 (array []byte) uint32 {
	for i, v := range(array) {
		array[i] = bits.Reverse8(v)
	}
	for i, v := range(array[0: 4]) {
		array[i] = v ^ 0xFF
	}
	newReminder := bits.Reverse32(reminderIEEEGreater4Padded(array))

	return newReminder ^ 0xFFFFFFFF
}

func ieeeRudimentarySparse (file SparseFile) uint32 {
	position2Index := map[uint64]int{}
	for i , v := range(file.fileBytes) {
		file.fileBytes[i].Value = bits.Reverse8(v.Value)
		position2Index[v.Position] = i
	}
	for i := 0; i < 4; i++ {
		index, ok := position2Index[uint64(i)]
		if !ok {
			file.fileBytes = append(file.fileBytes, SparseByte{Position: uint64(i), Value: 0xFF})
		} else {
			file.fileBytes[index].Value ^= 0xFF
		}
	}

	newReminder := bits.Reverse32(reminderIEEEGreater4SparsePadded(file))

	return newReminder ^ 0xFFFFFFFF
}
type SparseFile struct {
	fileBytes []SparseByte
	size uint64
}
type SparseByte struct {
	Position uint64 `index:"0" delimiter:"space"`
	Value    byte `index:"1"`
}

func (f SparseFile) toBytes () []byte {
	slice := make([]byte, f.size)
	for _, v := range f.fileBytes {
		slice[v.Position] = v.Value
	}
	return slice
}




func reminderIEEE(array []byte) uint32{
	poly := bits.Reverse32(uint32(IEEE))
	auxArray := make([]byte, 4)
	for i := 0; i < minInt(4, len(array)); i++ {
		auxArray[i] = array[i]
	}

	current := 0xFFFFFFFF ^ (uint32(auxArray[0]) << 24 + uint32(auxArray[1]) << 16 + uint32(auxArray[2]) << 8 + uint32(auxArray[3]))
	if current & (1 << 31) == 1 {
		current ^= poly
	}
	for i := 0; i < len(array); i++ {
		nextIndex := i + 4;
		var nextValue byte
		if nextIndex >= len(array) {
			nextValue = 0
		} else {
			nextValue = array[nextIndex]
		}
		for next := byte(1 << 7); next > 0; next >>= 1 {
			// log.Println(fmt.Sprintf("%8b", next))
			// log.Println(fmt.Sprintf("%8b", nextValue))
			nextBit := nextValue & next
			// log.Println(fmt.Sprintf("%8b", nextBit))
			if current >> 31 & 1 == 1 {
				current = ((current << 1) ^ uint32(nextBit)) ^ poly
			} else {
				current = (current << 1) ^ uint32(nextBit)
			}
		}
	}
	return bits.Reverse32(current)
}

func reminderIEEEGreater4Padded(array []byte) uint32{
	array = append(array, 0, 0, 0, 0)
	var i int
	carryOver := uint32(1)
	result := uint64(0)
	increase := reminderIEEEGreater4([]byte{1, 0, 0, 0, 0})
	for i = len(array); i > 8; i -= 4 {
		result ^= xorMultiply32(reminderIEEEGreater4(array[i - 4: i]), carryOver)
		carryOver = reminderIEEEGreater4(uint64ToArray(xorMultiply32(carryOver, increase)))
		// log.Println("a", fmt.Sprintf("%8b", uint64ToArray(xorMultiply32(carryOver, increase))))
		// log.Println("a", fmt.Sprintf("%8b", uint64ToArray(xorMultiply32(carryOver, increase))))
	}
	result ^= xorMultiply32(reminderIEEEGreater4(array[0: i]) , carryOver)

	return reminderIEEEGreater4(uint64ToArray(result))
}

func reminderIEEEGreater4SparsePadded (file SparseFile) uint32{
	file.size += 4
	return reminderIEEEGreater4Sparse(file)
}
func reminderIEEEGreater4Sparse(file SparseFile) uint32{
	size := file.size
	result := uint32(0)
	for _, sparseByte := range file.fileBytes {
		power := size - 1 - sparseByte.Position
		valueCoeff := reminderIEEEGreater4([]byte{0, 0, 0, sparseByte.Value})
		carry := uint32(1)
		power2Value := reminderIEEEGreater4([]byte{0, 0, 1, 0})
		for power > 0 {
			if power % 2 == 1 {
				carry = multiplyOnModule(carry, power2Value)
			}
			power = power >> 1
			power2Value = multiplyOnModule(power2Value, power2Value)
		}
		result ^= multiplyOnModule(carry, valueCoeff)
	}

	return result
}

func multiplyOnModule (a, b uint32) uint32 {
	return reminderIEEEGreater4(uint64ToArray(xorMultiply32(a, b)))
}

func reminderIEEEGreater4(array []byte) uint32{
	poly := bits.Reverse32(uint32(IEEE))
	current := (uint32(array[0]) << 24 + uint32(array[1]) << 16 + uint32(array[2]) << 8 + uint32(array[3]))
	if current & (1 << 31) == 1 {
		current ^= poly
	}
	for i := 0; i < len(array) - 4; i++ {
		nextIndex := i + 4;
		var nextValue byte
		nextValue = array[nextIndex]

		for i := 7; i >= 0; i-- {
			nextBit := (nextValue >> i) & 1
			if current >> 31 & 1 == 1 {
				current = ((current << 1) ^ uint32(nextBit)) ^ poly
			} else {
				current = (current << 1) ^ uint32(nextBit)
			}
		}
	}
	return current
}




func maxInt (a, b int) int {
	if a < b {
		return b
	}
	return a
}
func minInt (a, b int) int {
	if a > b {
		return b
	}
	return a
}


func xorMultiply32 (a, b uint32) uint64 {
	return xorMultiply(uint64(a), uint64(b))
}
func xorMultiply (a, b uint64) uint64{
	if a<1{return 0};return a%2*b^xorMultiply(a/2,b*2)
}
func uint32ToArray(a uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(a))
	return b
}
func arrayToUint32(input []byte) uint32 {
	return  binary.BigEndian.Uint32(input)
}

func arrayToUint64(input []byte) uint64 {
	return  binary.BigEndian.Uint64(input)
}

func uint64ToArray (a uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(a))
	return b
}



type Challenge struct {
	NCases int    `index:"0"`
	Cases  []Case `index:"1" indexed:"NCases"`
}
type Case struct {
	FileName   string       `index:"0" delimiter:"space" `
	NAdditions int          `index:"1"`
	Additions  []SparseByte `index:"2" indexed:"NAdditions"`
}
