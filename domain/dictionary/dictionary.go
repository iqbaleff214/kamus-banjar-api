package dictionary

import "strings"

type Entries struct {
	index map[string]int
	words []Word
}

type Dictionary struct {
	entries map[string]Entries
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
		entries := Entries{
			index: map[string]int{},
			words: []Word{},
		}
		d.entries[alphabet] = entries
	}

	entries := d.entries[alphabet]
	entries.index[word.Word] = index
	entries.words = append(entries.words, word)
	d.entries[alphabet] = entries
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
