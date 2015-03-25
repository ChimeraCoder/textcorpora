package cmu

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"unicode"

	"github.com/ChimeraCoder/textcorpora"
	"github.com/Wessie/appdirs"
)

var app = appdirs.New("cmudict", "chimeracoder", ".1")

type cmuCorpus map[string][]string

var cmuCorpusCached cmuCorpus = map[string][]string{}

const _CorpusFilename = "cmudict.0.7a.corpus"
const _CorpusUrl = "http://svn.code.sf.net/p/cmusphinx/code/trunk/cmudict/cmudict.0.7a"

// CMUCorpus returns the CMU corpus
func CMUCorpus() (textcorpora.Corpus, error) {

	filename := path.Join(app.UserData(), _CorpusFilename)
	// Check if file already exists
	if _, err := os.Stat(filename); err != nil {

		log.Printf("Writing to filename %s", filename)

		err := os.MkdirAll(app.UserData(), os.ModePerm)
		if err != nil {
			return nil, err
		}
		out, err := os.Create(filename)
		if err != nil {
			panic(err)
			return nil, err
		}
		defer out.Close()
		log.Printf("Fetching url %s", _CorpusUrl)
		resp, err := http.Get(_CorpusUrl)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		n, err := io.Copy(out, resp.Body)
		log.Printf("Wrote %d bytes", n)
	}

	bts, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cmu := string(bts)

	re := regexp.MustCompile(`^[A-Z]`)
	var tmpCorpus cmuCorpus = map[string][]string{}
	for _, line := range strings.Split(cmu, "\n") {
		line = strings.TrimSpace(line)

		//TODO account for the lines that represent alternate pronunciations

		// Ignore the lines that don't start with a character A-Z
		if len(line) == 0 || !re.MatchString(line[:1]) {
			continue
		}

		linesplit := strings.Split(line, " ")
		word := linesplit[0]
		tmpCorpus[strings.ToUpper(word)] = linesplit[1:]
	}

	return tmpCorpus, nil
}

// Syllables returns the number of syllables for the word, according to the corpus
// If the word is not in the corpus, it will return 0
func (c cmuCorpus) Syllables(word string) int {
	phonemes, ok := c[strings.ToUpper(word)]
	if !ok {
		return 0
	}

	count := 0
	for _, phoneme := range phonemes {
		for _, r := range phoneme {
			if unicode.IsNumber(r) {
				count++
			}
		}
	}
	return count
}

// Words returns the number of words in the corpus
func (c cmuCorpus) Words() int {
	return len(c)
}

// Words cursor returns a channel that can be used to
// iterate over the words in the corpus
func (c cmuCorpus) WordsCursor() (cursor chan string) {
	cursor = make(chan string)

	go func() {
		for word := range map[string][]string(c) {
			cursor <- word
		}
		close(cursor)
	}()

	return
}
