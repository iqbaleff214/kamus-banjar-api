package dictionary

import (
	"errors"
	"sort"
	"strings"

	"github.com/agnivade/levenshtein"
)

type Entries struct {
	index map[string]int
	words []Word
}

type Dictionary struct {
	entries map[string]*Entries
}

func NewDictionary() Dictionary {
	return Dictionary{
		entries: make(map[string]*Entries),
	}
}

func (d *Dictionary) AddEntry(index int, word Word) {
	// it is possible that Word have empty Alphabet
	// ignore alphabet from entry
	alphabet := strings.ToLower(string(word.Word[0]))

	if _, ok := d.entries[alphabet]; !ok {
		d.entries[alphabet] = &Entries{
			index: map[string]int{},
			words: []Word{},
		}
	}

	e := d.entries[alphabet]
	e.index[word.Word] = index
	e.words = append(e.words, word)
}

func (d *Dictionary) AddEntries(words []Word) {
	for idx, word := range words {
		d.AddEntry(idx, word)
	}
}

func (d *Dictionary) GetWord(word string) (Word, bool) {
	ch := strings.ToLower(string(word[0]))
	e, ok := d.entries[ch]
	if !ok {
		return Word{}, false
	}

	if idx, ok := e.index[word]; !ok {
		return Word{}, false
	} else {
		return e.words[idx], true
	}
}

func (d *Dictionary) GetEntries(alphabet string) []Word {
	e, ok := d.entries[alphabet]
	if !ok {
		return nil
	}

	return e.words
}

func (d *Dictionary) Search(keyword string) (SearchResult, error) {
	matches := []string{}

	if w, ok := d.GetWord(keyword); ok {
		return SearchResult{
			Search: keyword,
			Words:  []string{w.Word},
			Total:  1,
		}, nil
	}

	// Scan all words across all buckets. A front-insertion/deletion can shift
	// the first letter, so limiting to one bucket would miss valid matches.
	// With ~250 words this is fast; the main saving vs. the old code is that
	// we no longer maintain a redundant flat d.words []string copy.
	for _, e := range d.entries {
		for _, w := range e.words {
			if levenshtein.ComputeDistance(w.Word, keyword) == 1 {
				matches = append(matches, w.Word)
			}
		}
	}

	sort.Strings(matches)
	result := SearchResult{
		Search: keyword,
		Words:  matches,
		Total:  len(matches),
	}

	if result.Total > 0 {
		return result, nil
	}

	return result, errors.New("no matching word found")
}
