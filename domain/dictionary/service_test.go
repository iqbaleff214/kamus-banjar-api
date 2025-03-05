package dictionary_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iqbaleff214/kamus-banjar-api/domain/dictionary"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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
	curDir, _ := os.Getwd()
	curDir = filepath.Join(curDir, "..", "..")

	repo := dictionary.NewRepository(curDir)
	serv := dictionary.NewService(repo)

	type testCase struct {
		name string
		word string
		err  error
		exp  dictionary.Word
	}

	cases := []testCase{
		{"should return error when alphabet is invalid", "vanci", fmt.Errorf("the word is not found"), dictionary.Word{}},                  // V is invalid alphabet in Banjarese
		{"should return error when alphabet is valid but word not found", "yaeni", fmt.Errorf("the word is not found"), dictionary.Word{}}, // Y is a valid alphabet, but yaeni is not a word
	}


	for _, l := range []string{
		"a", "b", "c", "d", "g", "h",
		"i", "j", "k", "l", "m", "n",
		"p", "r", "s", "t", "u", "w",
		"y",
	} {
		b, err := os.ReadFile(filepath.Join(curDir, "data", l + ".json"))
		if err != nil {
			t.Fatal(err)
		}
		var words []dictionary.Word
		if err = json.Unmarshal(b, &words); err != nil {
			t.Fatal(err)
		}
		for _, w := range words {
			cases = append(cases, testCase{
				name: "should return definition of " + w.Word,
				word: w.Word,
				err:  nil,
				exp:  w,
			})
		}
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			act, er := serv.GetWord(tt.word)

			if tt.err == nil {
				assert.Equal(t, tt.word, act.Word)
				assert.Equal(t, tt.exp.Word, act.Word)
				assert.Equal(t, tt.exp.Alphabet, act.Alphabet)
				assert.Equal(t, tt.exp.Syllable, act.Syllable)
			} else {
				assert.EqualError(t, er, tt.err.Error())
			}
		})
	}
}

func Test_service_GetWordsByAlphabet(t *testing.T) {
	curDir, _ := os.Getwd()
	curDir = filepath.Join(curDir, "..", "..")

	repo := dictionary.NewRepository(curDir)
	serv := dictionary.NewService(repo)

	type testCase struct {
		name     string
		alphabet string
		err      error
		total    int
	}

	var cases []testCase

	letters := []string{
		"a", "b", "c", "d", "g", "h",
		"i", "j", "k", "l", "m", "n",
		"p", "r", "s", "t", "u", "w",
		"y",
	}
	for _, l := range letters {
		if b, err := os.ReadFile(filepath.Join(curDir, "data", l+".json")); err == nil {
			var words []any
			if err = json.Unmarshal(b, &words); err == nil {
				cases = append(cases, testCase{
					name:     "should return alphabet " + l + " and total words",
					alphabet: l,
					err:      nil,
					total:    len(words),
				})
			}
		}
	}

	cases = append(cases, testCase{"should return error when alphabet is invalid", "z", fmt.Errorf("the alphabet is not available"), 0})                  // V is invalid alphabet in Banjarese
	cases = append(cases, testCase{"should return error when request is not a letter", "+", fmt.Errorf("symbols are not allowed"), 0})                    // Y is a valid alphabet, but yaeni is not a word
	cases = append(cases, testCase{"should return error when request is not a single character", "go", fmt.Errorf("alphabet only has one character"), 0}) // Y is a valid alphabet, but yaeni is not a word

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			act, words, er := serv.GetWordsByAlphabet(tt.alphabet)

			if tt.err == nil {
				assert.Equal(t, tt.alphabet, act.Letter)
				assert.Equal(t, tt.total, len(words))
				assert.Nil(t, er)
			} else {
				assert.EqualError(t, er, tt.err.Error())
			}
		})
	}
}

func Test_service_Search(t *testing.T) {
	curDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}

	curDir = filepath.Join(curDir, "..", "..")

	repo := dictionary.NewRepository(curDir)
	serv := dictionary.NewService(repo)

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
		})
	}
}
