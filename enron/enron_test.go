package enron

import (
	"path"
	"testing"
)

func Test_GetCorpus(t *testing.T) {
	_, err := EnronCorpus()
	if err != nil {
		t.Error(err)
	}
}

func Test_Untar(t *testing.T) {
	err := untar(path.Join(app.UserData(), _CorpusFilename))
	if err != nil {
		t.Error(err)
	}
}
