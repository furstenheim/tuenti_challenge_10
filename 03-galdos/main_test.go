package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestParseBook (t *testing.T) {
	book := parseBook(`Release Date: November 5, 2005 [EBook #17013]
[Last updated on December 21, 2019]
Mi primer hijo--decía--nació cuando vino la tropa carlista
«D. Plácido, ¿tiene usted pana azul?».--«¡Pana azul!
--¿Vive aquí--le preguntó--el Sr. de Estupiñá?
--¡Qué poca vergüenza!
Pues apechuguemos con las _novedades_» dijo Isabel a su marido
Barbarita declaraba riendo que con estos teje-manejes se había vuelto, sin saberlo, una doña Beatriz Galindo
¡Pero las niñas!... ¡Y con estas modas de ahora y este suponer!... ¿Viste la pieza de merino azul?, pues no fue bastante`)
	assert.Equal(t, `release date  november          ebook          last updated on december           mi primer hijo  decía  nació cuando vino la tropa carlista  d  plácido   tiene usted pana azul       pana azul     vive aquí  le preguntó  el sr  de estupiñá     qué poca vergüenza  pues apechuguemos con las  novedades   dijo isabel a su marido barbarita declaraba riendo que con estos teje manejes se había vuelto  sin saberlo  una doña beatriz galindo  pero las niñas      y con estas modas de ahora y este suponer      viste la pieza de merino azul   pues no fue bastante`, book)
}

func TestParseFrequencyMap (t *testing.T) {
	input := `release date  november          ebook          last updated on december           mi primer hijo  decía  nació cuando vino la tropa carlista  d  plácido   tiene usted pana azul       pana azul     vive aquí  le preguntó  el sr  de estupiñá     qué poca vergüenza  pues apechuguemos con las  novedades   dijo isabel a su marido barbarita declaraba riendo que con estos teje manejes se había vuelto  sin saberlo  una doña beatriz galindo  pero las niñas      y con estas modas de ahora y este suponer      viste la pieza de merino azul   pues no fue bastante`
	frequencyMap := parseFrequencyMap(input)

	i, ok := frequencyMap["ahora"]
	log.Println(i, ok)
	log.Println(frequencyMap)
	assert.Equal(t, 2, frequencyMap["pana"])
	assert.Equal(t, 0, frequencyMap["y"])
	assert.Equal(t, 1, frequencyMap["estupiñá"])
	assert.Equal(t, 0, frequencyMap["sr"])
	assert.Equal(t, 1, frequencyMap["sin"])
}

func TestParseFrequencyMap_twoLettersWithAccent (t *testing.T) {
	input := `release él`
	frequencyMap := parseFrequencyMap(input)

	assert.Equal(t, 0, frequencyMap["él"])
	assert.Equal(t, 1, frequencyMap["release"])
}

func TestLessSpanish (t * testing.T) {
	ts := []struct{
		a,b string
		expected bool
	}{

		{"a", "b", true},
		{"b", "a", false},
		{"b", "ab", false},
		{"ab", "abc", true},
		{"abc", "ab", false},
		{"abc", "abá", true},
		{"á", "a", false},
		{"a", "b", true},
		{"b", "c", true},
		{"c", "d", true},
		{"d", "e", true},
		{"e", "f", true},
		{"f", "g", true},
		{"g", "h", true},
		{"h", "i", true},
		{"i", "j", true},
		{"j", "k", true},
		{"k", "l", true},
		{"l", "m", true},
		{"m", "n", true},
		{"n", "o", true},
		{"o", "p", true},
		{"p", "q", true},
		{"q", "r", true},
		{"r", "s", true},
		{"s", "t", true},
		{"t", "u", true},
		{"u", "v", true},
		{"v", "w", true},
		{"w", "x", true},
		{"x", "y", true},
		{"y", "z", true},
		{"z", "á", true},
		{"á", "é", true},
		{"é", "í", true},
		{"í", "ñ", true},
		{"ñ", "ó", true},
		{"ó", "ú", true},
		{"ú", "ü", true},
	}
		for _, tc := range(ts) {
		spanishLess := lessSpanish(tc.a, tc.b)
		assert.Equal(t, tc.expected, spanishLess, fmt.Sprintf("%s %s", tc.a, tc.b))
	}
}

func TestParseDictionary1 (t *testing.T) {
	input := `zumo ámbar doctor déjeme ejército ómnibus úlceras volveré ñoños baldomero ropa abrazo`
	frequencyMap := parseFrequencyMap(input)
	dictionary := parseFrequencyDictionary(frequencyMap)
	assert.Equal(t, []PositionResult{
		{Count: 1, Word:  "abrazo"},
		{Count: 1, Word:  "baldomero"},
		{Count: 1, Word:  "doctor"},
		{Count: 1, Word:  "déjeme"},
		{Count: 1, Word:  "ejército"},
		{Count: 1, Word:  "ropa"},
		{Count: 1, Word:  "volveré"},
		{Count: 1, Word:  "zumo"},
		{Count: 1, Word:  "ámbar"},
		{Count: 1, Word:  "ñoños"},
		{Count: 1, Word:  "ómnibus"},
		{Count: 1, Word:  "úlceras"},
	}, dictionary.ByPosition)
	assert.Equal(t, map[string]WordResult{
		"abrazo": {Count: 1, Position: 0},
		"baldomero": {Count: 1, Position: 1},
		"doctor": {Count: 1, Position: 2},
		"déjeme": {Count: 1, Position: 3},
		"ejército": {Count: 1, Position: 4},
		"ropa": {Count: 1, Position: 5},
		"volveré": {Count: 1, Position: 6},
		"zumo": {Count: 1, Position: 7},
		"ámbar": {Count: 1, Position: 8},
		"ñoños": {Count: 1, Position: 9},
		"ómnibus": {Count: 1, Position: 10},
		"úlceras": {Count: 1, Position: 11},
	}, dictionary.ByWord)

	log.Println(dictionary)
}


