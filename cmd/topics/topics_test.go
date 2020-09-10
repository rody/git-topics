package topics

import (
	"testing"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestToCommitInfo(t *testing.T) {
	c := object.Commit{
		Hash: plumbing.NewHash("01234567890123456789"),
		Author: object.Signature{
			Name:  "John Doe",
			Email: "john.doe@example.com",
			When:  time.Now(),
		},
		Committer: object.Signature{
			Name:  "Jane Doe",
			Email: "jane.doe@example.com",
			When:  time.Now(),
		},
		Message: "First line of message.\nSecond line\n",
	}

	ci := toCommitInfo(&c)

	if ci == nil {
		t.Fatal("expected commit info to be non nil")
	}

}
