package spelling

import (
	"os"
	"regexp"
)

const (
	BigTxtFilePath = "./internal/spelling/big.txt"
	Letters        = "abcdefghijklmnopqrstuvwxyz"
)

type Corrector struct {
	// internalState represents if the spelling corrector was initialized
	internalState bool
	// appearanceMap represents a map with each word from the file and the number of appearances
	appearanceMap map[string]int
	totalWords    int
}

func New() *Corrector {
	return &Corrector{
		internalState: false,
		appearanceMap: make(map[string]int),
	}
}

func (c *Corrector) Initialize() error {
	data, err := os.ReadFile(BigTxtFilePath)
	if err != nil {
		return err
	}

	content := string(data)
	wordRegex := regexp.MustCompile(`\w+`)
	words := wordRegex.FindAllString(content, -1)

	c.totalWords = len(words)
	for _, word := range words {
		c.appearanceMap[word] += 1
	}

	c.internalState = true
	return nil
}

func (c *Corrector) Initialized() bool {
	return c.internalState
}

// Correction most probable spelling correction for word.
func (c *Corrector) Correction(word string) string {
	possibleWords := c.known(c.edits2(word))
	maxProbability := 0.0
	bestWord := ""

	for _, possibleWord := range possibleWords {
		probability := c.wordProbability(possibleWord)

		if probability > maxProbability {
			maxProbability = probability
			bestWord = possibleWord
		}
	}

	return bestWord
}

// wordProbability probability of a 'word'.
func (c *Corrector) wordProbability(word string) float64 {
	return float64(c.appearanceMap[word]) / float64(c.totalWords)
}

// known the subset of 'words' that appear in the Corrector.appearanceMap
func (c *Corrector) known(words []string) []string {
	knownWords := make([]string, 0)

	for _, word := range words {
		_, exists := c.appearanceMap[word]
		if exists {
			knownWords = append(knownWords, word)
		}
	}

	return knownWords
}

// edits1 all edits that are one edit away from 'word'
func (c *Corrector) edits1(word string) map[string]struct{} {
	splits := make(map[string]string)
	for i := 0; i < len(word)+1; i++ {
		L := word[:i]
		R := word[i:]

		splits[L] = R
	}

	edits := make(map[string]struct{})

	// deletes
	for L := range splits {
		R := splits[L]

		if R != "" {
			edits[L+R[1:]] = struct{}{}
		}
	}

	// transposes
	for L := range splits {
		R := splits[L]

		if len(R) > 1 {
			edits[L+string(R[1])+string(R[0])+R[2:]] = struct{}{}
		}
	}

	// replaces
	for L := range splits {
		R := splits[L]

		if R != "" {
			for i := range Letters {
				c := string(Letters[i])
				edits[L+c+R[1:]] = struct{}{}
			}
		}
	}

	// inserts
	for L := range splits {
		R := splits[L]

		for i := range Letters {
			c := string(Letters[i])
			edits[L+c+R] = struct{}{}
		}
	}

	return edits
}

// edits2 all edits that are two edits away from 'word'
func (c *Corrector) edits2(word string) []string {
	edits1 := c.edits1(word)
	edits2 := make([]string, 0)

	for edit := range edits1 {
		newEdits := c.edits1(edit)

		for newEdit := range newEdits {
			edits2 = append(edits2, newEdit)
		}
	}

	return edits2
}
