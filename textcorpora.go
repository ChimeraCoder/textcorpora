package textcorpora

// A Corpus is a body of text that supports certain features
// Currently the only required queries are Syllables (number of syllables)
// and Words (number of words)
type Corpus interface {
	Syllables(string) int
	Words() int
}
