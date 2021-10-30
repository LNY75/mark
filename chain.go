package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
)

// Prefix is a Markov chain prefix of one or more words.
type Prefix []string

// String returns the Prefix as a string (for use as a map key).
func (p Prefix) String() string {
	return strings.Join(p, " ")
}

// Shift removes the first word from the Prefix and appends the given word.
func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

// Chain contains a map ("chain") of prefixes to a list of suffixes.
// A prefix is a string of prefixLen words joined with spaces.
// A suffix is a single word. A prefix can have multiple suffixes.
// We map the prefix (a string) to a list of its suffices (a list of strings)
type Chain struct {
	chain     map[string][]string
	prefixLen int
}

// NewChain returns a new Chain with prefixes of prefixLen words.
func NewChain(prefixLen int) *Chain {
	return &Chain{make(map[string][]string), prefixLen}
}

// Build reads text from the provided Reader and
// parses it into prefixes and suffixes that are stored in Chain.
func (c *Chain) Build(r io.Reader) {
	br := bufio.NewReader(r)
	p := make(Prefix, c.prefixLen)
	for {
		var s string
		if _, err := fmt.Fscan(br, &s); err != nil {
			break
		}
		// join the words in Prefix, separated by space
		key := p.String()
		// map the suffix to the prefix. A prefix is mapped to a list of suffices.
		c.chain[key] = append(c.chain[key], s)
		p.Shift(s)
	}
}

// Build reads text from a file; parse them into prefixes and suffixies to be stored in Chain
func (c *Chain) BuildFromFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		panic("can't open file")
	}
	scanner := bufio.NewScanner(file)
	var words []string
	p := make(Prefix, c.prefixLen)

	for scanner.Scan() {
		// strings.Split doesn't split quite properly here when there is more than one consecutive white space
		line := strings.Fields(scanner.Text())
		for _, word := range line {
			words = append(words, word)
		}
	}
	if scanner.Err() != nil {
		panic("something went wrong when reading the file?")
	}
	file.Close()

	for _, word := range words {
		key := p.String()
		c.chain[key] = append(c.chain[key], word)
		p.Shift(word)
	}
}

// Generate returns a string of at most n words generated from Chain.
func (c *Chain) Generate(n int) string {
	p := make(Prefix, c.prefixLen)
	var words []string
	for i := 0; i < n; i++ {
		choices := c.chain[p.String()]
		if len(choices) == 0 {
			break
		}
		next := choices[rand.Intn(len(choices))]
		words = append(words, next)
		p.Shift(next)
	}
	return strings.Join(words, " ")
}
