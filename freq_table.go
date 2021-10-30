package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type FreqTable struct {
	table     map[string]map[string]int
	prefixLen int
}

// NewFreqTable returns a new frequency table
func NewFreqTable(prefixLen int) *FreqTable {
	return &FreqTable{make(map[string]map[string]int), prefixLen}
}

// Build takes a Chain and parses its contents to be stored in a frequency table
func (ft *FreqTable) Build(c Chain) {
	for prefix, suffixes := range c.chain {
		// must initialize the map first, otherwise we are adding entries to nil
		var freqs map[string]int = make(map[string]int)
		for _, s := range suffixes {
			freqs[s]++
		}
		ft.table[prefix] = freqs
	}
}

// TotalFreq gives the number of occurrences of all suffices mapped to a given prefix
// This provides a reference for generating a random number in Generate
func (ft *FreqTable) TotalFreq(prefix string) int {
	var n int = 0
	for _, count := range ft.table[prefix] {
		n += count
	}
	return n
}

// SelectSuffix selects a possible suffix given a prefix
func (ft *FreqTable) SelectSuffix(prefix string) string {
	var next string = prefix
	totalFreqCount := ft.TotalFreq(prefix)
	n := rand.Intn(totalFreqCount)
	var sum int = 0
	// this "random" selection of a suffix works by
	for p, f := range ft.table[prefix] {
		sum += f
		if n <= sum {
			next = p
			break
		}
	}
	return next
}

// Generate returns a string of at most n words, which is what the program tries to ramble based on the frequency table
func (ft *FreqTable) Generate(n int) string {
	var text []string
	p := make(Prefix, ft.prefixLen)
	for i := 0; i < n; i++ {
		pString := p.String()
		choiceMap := ft.table[pString]
		if len(choiceMap) == 0 {
			break
		}
		next := ft.SelectSuffix(pString)
		text = append(text, next)
		p.Shift(next)
	}
	return strings.Join(text, " ")
}

// PrintFreqTable prints all key-value pairs from a FrepTable; a helper function for debugging
func (ft *FreqTable) Print() {
	for prefix, freqs := range ft.table {
		fmt.Print(prefix, " | ")
		for suffix, count := range freqs {
			fmt.Print(suffix, " ", count, ", ")
		}
		fmt.Println("")
	}
}

// PrintMap - prints the suffix map
func (ft *FreqTable) PrintMap(prefix string) {
	for suffix, count := range ft.table[prefix] {
		fmt.Println(prefix, ",", suffix, count)
	}
}

// WriteMapToFile generates a file that stores the frequency map; it takes the output filename as a param
func (ft *FreqTable) WriteMapToFile(output string) {
	// create output file
	outFile, err := os.Create(output)
	if err != nil {
		panic("lol, cannot create file.")
	}

	// write the prefix length
	fmt.Fprintln(outFile, ft.prefixLen)

	for prefix, freqMap := range ft.table {
		wordsInPrefix := strings.Fields(prefix)

		if len(wordsInPrefix) < ft.prefixLen {
			for i := 0; i < (ft.prefixLen - len(wordsInPrefix)); i++ {
				fmt.Fprint(outFile, "\"\" ")
			}
		}

		for _, word := range wordsInPrefix {
			fmt.Fprint(outFile, word, " ")
		}

		for suffix, freq := range freqMap {
			fmt.Fprint(outFile, suffix, " ", freq, " ")
		}
		fmt.Fprintln(outFile)
	}
	outFile.Close()
}

// ConvertFreqsToMap parses the list of suffixes and their frequencies to a map of suffix to the corresponding frequency
func ConvertFreqsToMap(freqs []string) map[string]int {
	freqMap := make(map[string]int)
	if len(freqs)%2 != 0 {
		panic("your freqs list does not have even length; one suffix has no frequency.")
	} else {
		for i := 0; i < len(freqs)-1; i += 2 {
			for j := 0; j < 2; j++ {
				freq, err := strconv.Atoi(freqs[i+1])
				if err != nil {
					panic("cannot convert string to integer")
				}
				freqMap[freqs[i]] = freq
			}
		}
	}
	return freqMap
}

// ReadTableFromFile parses a frequency table file (.txt) into a FreqTable
func ReadTableFromFile(input string) FreqTable {

	file, err := os.Open(input)
	if err != nil {
		panic("can't open file")
	}
	scanner := bufio.NewScanner(file)

	var prefixLen int
	// read the prefixlen from the first line of the table file first
	for scanner.Scan() {
		fl := strings.Fields(scanner.Text())
		if len(fl) == 0 {
			panic("wow, your frequency table file does not specify the prefix length. Can't go any further from here.")
		} else {
			l, err := strconv.Atoi(fl[0])
			if err != nil {
				panic("Error tryign to convert prefix length from file to an int")
			}
			prefixLen = l
		}
		break
	}

	// now that we know the prefix length, we can make a table
	freqTable := *NewFreqTable(prefixLen)

	i := 1
	for scanner.Scan() {
		// skip the fist line
		if i > 1 {
			i++
			continue
		}
		line := strings.Fields(scanner.Text())
		prefix := line[0:prefixLen]
		freqs := line[prefixLen:]

		for i, p := range prefix {
			if p == "\"\"" {
				prefix[i] = ""
			}
		}

		freqTable.table[strings.Join(prefix, " ")] = ConvertFreqsToMap(freqs)
	}

	if scanner.Err() != nil {
		panic("something went wrong when reading the file?")
	}
	file.Close()

	return freqTable
}
