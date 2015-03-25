package textcorpora

// A Corpus is a body of text that supports certain features
// Currently the only required queries are Syllables (number of syllables),
// Words (number of words), and WordsCursor() (for iterating over the words).
type Corpus interface {
	Syllables(string) int
	Words() int
	WordsCursor() chan string
}
