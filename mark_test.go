package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestBuildFreqTable(t *testing.T) {
	// Register command-line flags.
	prefixLen := flag.Int("prefix", 2, "prefix length in words")

	flag.Parse()                     // Parse command-line flags.
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator.

	c := NewChain(*prefixLen) // Initialize a new Chain.
	c.Build(os.Stdin)         // Build chains from standard input.

	table := NewFreqTable(*prefixLen)
	table.Build(*c)
	// table.Print()
}

func TestReadingFromFile(t *testing.T) {
	// Register command-line flags.
	// numWords := 100
	prefixLen := 2

	rand.Seed(time.Now().UnixNano()) // Seed the random number generator.

	c := NewChain(prefixLen)   // Initialize a new Chain.
	c.BuildFromFile("poe.txt") // Build chains from file.

	table := NewFreqTable(prefixLen)
	table.Build(*c)
	// text := table.Generate(numWords)
	// fmt.Print(text)
}

func TestWriteToFile(t *testing.T) {
	// Register command-line flags.
	prefixLen := 2

	rand.Seed(time.Now().UnixNano()) // Seed the random number generator.

	c := NewChain(prefixLen)        // Initialize a new Chain.
	c.BuildFromFile("shortPoe.txt") // Build chains from file.

	table := NewFreqTable(prefixLen)
	table.Build(*c)

	table.WriteMapToFile("freq_table.txt")
}

func TestReadTableFromFile(t *testing.T) {
	// sen := "you brought flowers for Geralt and Regis"
	// str := strings.Fields(sen)
	// str1 := str[0:2]
	// str2 := str[2:]
	// fmt.Println(len(str1))
	// fmt.Println(len(str2))
	// for _, word := range str1 {
	// 	fmt.Print(word, " ")
	// }
	// fmt.Println()
	// for _, word := range str2 {
	// 	fmt.Print(word, " ")
	// }
}

func TestConvertFreqsToMap(t *testing.T) {
	// var input []string
	// input = append(input, "Regis")
	// input = append(input, "1")
	// input = append(input, "Geralt")
	// input = append(input, "2")

	// m := make(map[string]int)
	// m = ConvertFreqsToMap(input)
	// fmt.Print(m["Regis"], m["Geralt"])
}

func TestReadFreqTableFromFile(t *testing.T) {
	// table := ReadTableFromFile("freq_table.txt", 2)
	// table.Print()
}

func TestReadFromFreqTableFile(t *testing.T) {
	table := ReadTableFromFile("freq_table.txt")
	text := table.Generate(50)
	fmt.Println(text)
}

func TestReadFromMultipleFiles(t *testing.T) {
	// var files []string
	// files = append(files, "shortPoe.txt")
	// files = append(files, "shortPoe2.txt")
	// files = append(files, "shortPoe3.txt")
	// prefixLen := 3
	// c := NewChain(prefixLen)
	// for _, file := range files {
	// 	c.BuildFromFile(file)
	// }
	// table := NewFreqTable(prefixLen)
	// table.Build(*c)
	// table.WriteMapToFile("outputTable.txt")
}
