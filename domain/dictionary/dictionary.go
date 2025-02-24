package dictionary

import (
	"errors"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type Entries struct {
	index map[string]int
	words []Word
}

type Dictionary struct {
	entries map[string]Entries
	words   []string
}

func NewDictionary() Dictionary {
	d := Dictionary{
		entries: make(map[string]Entries),
	}

	return d
}

func (d *Dictionary) AddEntry(index int, word Word) {
	// it is possible that Word have empty Alphabet
	// ignore alphabet from entry
	alphabet := strings.ToLower(string(word.Word[0]))

	_, ok := d.entries[alphabet]
	if !ok {
		d.entries[alphabet] = Entries{
			index: map[string]int{},
			words: []Word{},
		}
	}

	entries := d.entries[alphabet]
	entries.index[word.Word] = index
	entries.words = append(entries.words, word)
	d.entries[alphabet] = entries

	d.words = append(d.words, word.Word)
}

func (d *Dictionary) AddEntries(words []Word) {
	for idx, word := range words {
		d.AddEntry(idx, word)
	}
}

func (d *Dictionary) GetWord(word string) (Word, bool) {
	ch := string(word[0])
	entries, ok := d.entries[ch]
	if !ok {
		return Word{}, false
	}

	if idx, ok := entries.index[word]; !ok {
		return Word{}, false
	} else {
		return entries.words[idx], true
	}
}

func (d *Dictionary) GetEntries(alphabet string) []Word {
	entries, ok := d.entries[alphabet]
	if !ok {
		return nil
	}

	return entries.words
}

func (d *Dictionary) Search(keyword string) (SearchResult, error) {
	matches := fuzzy.Find(keyword, d.words)

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
