package enron

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"unicode"

	"github.com/ChimeraCoder/textcorpora"
	"github.com/Wessie/appdirs"
)

var app = appdirs.New("enroncorpus", "chimeracoder", ".1")

type enronCorpus map[string][]string

var enronCorpusCached enronCorpus = map[string][]string{}

const _CorpusFilename = "enron_mail_20110402.tgz"
const _CorpusUrl = "https://www.cs.enron.edu/~./enron/enron_mail_20110402.tgz"

func EnronCorpus() (textcorpora.Corpus, error) {

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

	return nil, nil
}

// Syllables returns the number of syllables for the word, according to the corpus
// If the word is not in the corpus, it will return 0
func (c enronCorpus) Syllables(word string) int {
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
func (c enronCorpus) Words() int {
	return len(c)
}

// Words cursor returns a channel that can be used to
// iterate over the words in the corpus
func (c enronCorpus) WordsCursor() (cursor chan string) {
	cursor = make(chan string)

	go func() {
		for word := range map[string][]string(c) {
			cursor <- word
		}
		close(cursor)
	}()

	return
}

func untar(infile string) error {
	fin, err := os.Open(infile)
	if err != nil {
		return err
	}
	defer fin.Close()

	// Open the gzip file for decompression into a tar archive
	gr, err := gzip.NewReader(bufio.NewReader(fin))
	if err != nil {
		return err
	}

	// Open the tar archive for reading.
	tr := tar.NewReader(bufio.NewReader(gr))

	// Iterate through the files in the archive.
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Contents of %s:\n", hdr.Name)
		if _, err := io.Copy(os.Stdout, tr); err != nil {
			log.Fatalln(err)
		}
		fmt.Println()
	}
	return err
}
