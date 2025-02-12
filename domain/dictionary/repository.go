package dictionary

import "embed"

// Repository contains method to interact with data source
type Repository interface {
	GetAlphabets() []Alphabet
	GetWordsByAlphabet(alphabet string) ([]Word, error)
}

// NewRepository is a function to instantiate new Repository object
func NewRepository(config map[string]any) Repository {
	if embedded, ok := config["embed"]; ok {
		repo := embedRepository{}
		fs := embedded.(embed.FS)
		repo.init(fs)
		return &repo
	}

	return jsonRepository{
		alphabets: []string{
			"a", "b", "c", "d", "g", "h",
			"i", "j", "k", "l", "m", "n",
			"p", "r", "s", "t", "u", "w",
			"y",
		},
	}
}
