package dictionary_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/iqbaleff214/kamus-banjar-api/domain/dictionary"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testLetters = []string{
	"a", "b", "c", "d", "g", "h",
	"i", "j", "k", "l", "m", "n",
	"p", "r", "s", "t", "u", "w",
	"y",
}

// loadWords reads all words from the data directory relative to the test file.
func loadWords(t *testing.T) []dictionary.Word {
	t.Helper()
	curDir, _ := os.Getwd()
	root := filepath.Join(curDir, "..", "..")

	var all []dictionary.Word
	for _, l := range testLetters {
		b, err := os.ReadFile(filepath.Join(root, "data", l+".json"))
		if err != nil {
			t.Fatalf("loadWords: read %s.json: %v", l, err)
		}
		var words []dictionary.Word
		if err = json.Unmarshal(b, &words); err != nil {
			t.Fatalf("loadWords: unmarshal %s.json: %v", l, err)
		}
		all = append(all, words...)
	}
	return all
}

// loadAlphabets builds the Alphabet list from JSON data files.
func loadAlphabets(t *testing.T) []dictionary.Alphabet {
	t.Helper()
	curDir, _ := os.Getwd()
	root := filepath.Join(curDir, "..", "..")

	alphabets := make([]dictionary.Alphabet, len(testLetters))
	for i, l := range testLetters {
		b, err := os.ReadFile(filepath.Join(root, "data", l+".json"))
		if err != nil {
			t.Fatalf("loadAlphabets: read %s.json: %v", l, err)
		}
		var words []any
		if err = json.Unmarshal(b, &words); err != nil {
			t.Fatalf("loadAlphabets: unmarshal %s.json: %v", l, err)
		}
		alphabets[i] = dictionary.Alphabet{Letter: l, Total: len(words)}
	}
	return alphabets
}

func TestNewService(t *testing.T) {
	serv := dictionary.NewService(nil)

	assert.NotNil(t, serv)
	assert.Implements(t, (*dictionary.Service)(nil), serv)
}

func Test_service_GetAlphabets(t *testing.T) {
	repo := new(dictionary.MockRepository)
	serv := dictionary.NewService(repo)

	t.Run("should return alphabets when data source is available", func(t *testing.T) {
		repo.
			On("GetAlphabets", mock.Anything).
			Return([]dictionary.Alphabet{}, nil).
			Once()

		act, err := serv.GetAlphabets()

		assert.Nil(t, err)
		assert.Equal(t, []dictionary.Alphabet{}, act)
		assert.NotNil(t, act)
		repo.AssertExpectations(t)
	})

	t.Run("should return error when data source is unavailable", func(t *testing.T) {
		repo.
			On("GetAlphabets", mock.Anything).
			Return([]dictionary.Alphabet{}, fmt.Errorf("error during getting alphabets")).
			Once()

		act, err := serv.GetAlphabets()

		assert.NotNil(t, err)
		assert.EqualError(t, err, "error during getting alphabets")
		assert.Equal(t, []dictionary.Alphabet{}, act)
		repo.AssertExpectations(t)
	})
}

func Test_service_GetWord(t *testing.T) {
	allWords := loadWords(t)

	type testCase struct {
		name string
		word string
		err  error
		exp  dictionary.Word
	}

	cases := []testCase{
		{"should return error when alphabet is invalid", "vanci", fmt.Errorf("the word is not found"), dictionary.Word{}},
		{"should return error when alphabet is valid but word not found", "yaeni", fmt.Errorf("the word is not found"), dictionary.Word{}},
	}

	for _, w := range allWords {
		cases = append(cases, testCase{
			name: "should return definition of " + w.Word,
			word: w.Word,
			err:  nil,
			exp:  w,
		})
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(dictionary.MockRepository)
			serv := dictionary.NewService(repo)

			if tt.err != nil {
				repo.On("GetWord", tt.word).Return(dictionary.Word{}, tt.err).Once()
			} else {
				repo.On("GetWord", tt.word).Return(tt.exp, nil).Once()
			}

			act, er := serv.GetWord(tt.word)

			if tt.err == nil {
				assert.Equal(t, tt.word, act.Word)
				assert.Equal(t, tt.exp.Word, act.Word)
				assert.Equal(t, tt.exp.Alphabet, act.Alphabet)
				assert.Equal(t, tt.exp.Syllable, act.Syllable)
			} else {
				assert.EqualError(t, er, tt.err.Error())
			}
			repo.AssertExpectations(t)
		})
	}
}

func Test_service_GetWordsByAlphabet(t *testing.T) {
	allAlphabets := loadAlphabets(t)

	type testCase struct {
		name     string
		alphabet string
		err      error
		total    int
		words    []dictionary.Word
	}

	var cases []testCase

	curDir, _ := os.Getwd()
	root := filepath.Join(curDir, "..", "..")

	for _, a := range allAlphabets {
		b, err := os.ReadFile(filepath.Join(root, "data", a.Letter+".json"))
		if err != nil {
			t.Fatalf("read %s.json: %v", a.Letter, err)
		}
		var words []dictionary.Word
		if err = json.Unmarshal(b, &words); err != nil {
			t.Fatalf("unmarshal %s.json: %v", a.Letter, err)
		}
		cases = append(cases, testCase{
			name:     "should return alphabet " + a.Letter + " and total words",
			alphabet: a.Letter,
			err:      nil,
			total:    len(words),
			words:    words,
		})
	}

	// "z" is a valid letter a-z but not in the Banjar alphabet list —
	// service calls GetAlphabets then GetWordsByAlphabet before returning the error.
	type errorCase struct {
		testCase
		needsAlphabetCall bool
		needsWordsCall    bool
	}

	errorCases := []errorCase{
		{testCase{"should return error when alphabet is not in the list", "z", fmt.Errorf("the alphabet is not available"), 0, nil}, true, true},
		{testCase{"should return error when request is not a letter", "+", fmt.Errorf("symbols are not allowed"), 0, nil}, false, false},
		{testCase{"should return error when request is not a single character", "go", fmt.Errorf("alphabet only has one character"), 0, nil}, false, false},
	}

	for _, ec := range errorCases {
		cases = append(cases, ec.testCase)
	}

	// index error cases by alphabet for mock setup
	errorCaseMap := map[string]errorCase{}
	for _, ec := range errorCases {
		errorCaseMap[ec.alphabet] = ec
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(dictionary.MockRepository)
			serv := dictionary.NewService(repo)

			ec, isErrorCase := errorCaseMap[tt.alphabet]
			if tt.err == nil {
				repo.On("GetAlphabets").Return(allAlphabets, nil).Once()
				repo.On("GetWordsByAlphabet", tt.alphabet).Return(tt.words, nil).Once()
			} else if isErrorCase {
				if ec.needsAlphabetCall {
					repo.On("GetAlphabets").Return(allAlphabets, nil).Once()
				}
				if ec.needsWordsCall {
					repo.On("GetWordsByAlphabet", tt.alphabet).Return([]dictionary.Word{}, fmt.Errorf("the alphabet is not available")).Once()
				}
			}

			act, words, er := serv.GetWordsByAlphabet(tt.alphabet)

			if tt.err == nil {
				assert.Equal(t, tt.alphabet, act.Letter)
				assert.Equal(t, tt.total, len(words))
				assert.Nil(t, er)
			} else {
				assert.EqualError(t, er, tt.err.Error())
			}
			repo.AssertExpectations(t)
		})
	}
}

func Test_service_Search(t *testing.T) {
	tests := []struct {
		name    string
		keyword string
		want    []string
		wantErr string
	}{
		{
			name:    "Exact match - word found",
			keyword: "abadan",
			want:    []string{"abadan"},
			wantErr: "",
		},
		{
			name:    "No match - word not found",
			keyword: "xyzabc",
			want:    []string{},
			wantErr: "no matching word found",
		},
		{
			name:    "Fuzzy search - similar words found",
			keyword: "aban",
			want:    []string{"abah", "abat", "abun", "acan", "aman", "amban", "saban"},
			wantErr: "",
		},
		{
			name:    "Case insensitive search",
			keyword: "ABUH",
			want:    []string{"abuh"},
			wantErr: "",
		},
		{
			name:    "Empty input",
			keyword: "",
			want:    nil,
			wantErr: "search keyword is required",
		},
		{
			name:    "Input with numbers",
			keyword: "ab4d4n",
			want:    nil,
			wantErr: "symbols are not allowed",
		},
		{
			name:    "Input with whitespace",
			keyword: " abadan ",
			want:    []string{"abadan"},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(dictionary.MockRepository)
			serv := dictionary.NewService(repo)

			if tt.wantErr == "" || (tt.wantErr == "no matching word found") {
				normalized := normalizeKeyword(tt.keyword)
				result := dictionary.SearchResult{Search: normalized, Words: tt.want, Total: len(tt.want)}
				var repoErr error
				if tt.wantErr == "no matching word found" {
					repoErr = fmt.Errorf("no matching word found")
				}
				repo.On("Search", normalized).Return(result, repoErr).Once()
			}

			got, err := serv.Search(tt.keyword)
			if tt.wantErr != "" {
				if err == nil {
					t.Errorf("Search(%s) error = nil, wantErr %v", got.Search, tt.wantErr)
					return
				}
				if err.Error() != tt.wantErr {
					t.Errorf("Search(%s) error = %v, wantErr %v", got.Search, err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("Search(%s) unexpected error = %v", got.Search, err)
				return
			}

			if !reflect.DeepEqual(got.Words, tt.want) {
				t.Errorf("Search(%s) = %v, want %v", got.Search, got.Words, tt.want)
			}
			repo.AssertExpectations(t)
		})
	}
}

func normalizeKeyword(keyword string) string {
	return strings.ToLower(strings.TrimSpace(keyword))
}
