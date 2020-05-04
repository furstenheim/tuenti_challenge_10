package main

import (
	"io/ioutil"
	"log"
	"math/big"
)

func main() {
	m1 := getBig("plaintexts/test1.txt")
	m2 := getBig("plaintexts/test2.txt")
	c1 := getBig("ciphered/test1.txt")
	c2 := getBig("ciphered/test2.txt")
	e := big.NewInt(65537)
	m1.Exp(m1, e, nil)
	m2.Exp(m2, e, nil)
	m1.Sub(m1, c1)
	m2.Sub(m2, c2)
	m1.GCD(nil, nil, m1, m2)
	log.Println(m1)
}

func getBig (file string) *big.Int {
	m1, err := ioutil.ReadFile(file)
	handleError(err)
	m1b := big.NewInt(0)
	m1b.SetBytes(m1)
	return m1b
}

func handleError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
func reverse (s []byte) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

