package enron

import (
	"path"
	"testing"
)

func Test_Untar(t *testing.T) {
	err := untar(path.Join(app.UserData(), _CorpusFilename))
	if err != nil {
		t.Error(err)
	}
}
