package enron

import (
	"path"
	"testing"
)

func Test_Untar(t *testing.T) {
	emails, err := untar(path.Join(app.UserData(), _CorpusFilename))
	if err != nil {
		t.Error(err)
	}
	// Try reading some emails for at least one non-empty, correctly-parsed one
	const numTries = 20
	i := 0
	for email := range emails {
		if i >= numTries {
			break
		}
		if email.Error == nil && len(email.Contents) > 0 {
			// This email was correctly parsed
			return
		}
		i++
	}

	t.Errorf("Read %d emails without finding a correctly-parsed one", numTries)

}
