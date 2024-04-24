package dictionary

// Alphabet model to contain alphabets
type Alphabet struct {
	Letter string `json:"letter"`
	Total  int    `json:"total"`
}

// Word model to contain word definition
type Word struct {
	Word        string           `json:"word"`
	Alphabet    string           `json:"alphabet"`
	Meanings    []WordMeaning    `json:"meanings"`
	Derivatives []WordDerivative `json:"derivatives"`
}

type WordMeaning struct {
	Definitions []WordDefinition `json:"definitions"`
}

type WordDefinition struct {
	Definition        string        `json:"definition,omitempty"`
	Refer             string        `json:"refer,omitempty"`
	LevelOfPoliteness int           `json:"levelOfPoliteness"`
	Examples          []WordExample `json:"examples,omitempty"`
}

type WordExample struct {
	Bjn string `json:"bjn,omitempty"`
	Id  string `json:"id,omitempty"`
}

type WordDerivative struct {
	Word        string           `json:"word"`
	Definitions []WordDefinition `json:"definitions"`
}

type AlphabeticWord map[string][]Word
