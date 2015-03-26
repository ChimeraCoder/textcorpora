package enron

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"os"
	"path"

	"github.com/Wessie/appdirs"
)

var app = appdirs.New("enroncorpus", "chimeracoder", ".1")

const _CorpusFilename = "enron_mail_20110402.tgz"
const _CorpusUrl = "https://www.cs.enron.edu/~./enron/enron_mail_20110402.tgz"


// TODO make this satisfy the textcorpora interface

type EnronCorpus struct {
	EmailsCursor chan Email
}

func NewCorpus() (*EnronCorpus, error) {

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

	emailsCursor, err := untar(filename)
	if err != nil {
		return nil, err
	}

	return &EnronCorpus{emailsCursor}, nil
}

type Email struct {
	Name     string
	Header   mail.Header
	Contents []byte

	TarHeader *tar.Header
	Error     error // If any error was encountered when parsing this email, it will be stored here
}

func (e Email) Parse() (*mail.Message, error) {
	return mail.ReadMessage(bytes.NewBuffer(e.Contents))
}

func untar(infile string) (chan Email, error) {
	result := make(chan Email)
	fin, err := os.Open(infile)
	if err != nil {
		return nil, err
	}

	// Open the gzip file for decompression into a tar archive
	gr, err := gzip.NewReader(bufio.NewReader(fin))
	if err != nil {
		return nil, err
	}

	// Open the tar archive for reading.
	tr := tar.NewReader(bufio.NewReader(gr))

	go func() {
		defer fin.Close()
		// Iterate through the files in the archive.
		for {
			e := &Email{}
			hdr, err := tr.Next()
			e.TarHeader = hdr
			if err == io.EOF {
				// end of tar archive
				break
			}
			if err != nil {
				e.Error = err
				result <- *e
				continue
			}
			message, err := mail.ReadMessage(tr)
			if err != nil {
				if err == io.EOF {
					// TODO figure out how to to handle this
					continue
				}
				e.Error = err
				result <- *e
				break
			}
			bts, err := ioutil.ReadAll(message.Body)
			if err != nil {
				e.Error = err
				result <- *e
				continue
			}
			e.Contents = bts
			result <- *e
		}
		close(result)
	}()
	return result, nil
}
