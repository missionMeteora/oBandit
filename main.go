package outputBandit

import (
	"os"
	"sync"
)

func New(outLoc, errLoc string) *OutputBandit {
	out, _ := os.Create(outLoc)
	err, _ := os.Create(errLoc)

	o := OutputBandit{
		out:        out,
		err:        err,
		origStdOut: os.Stdout,
		origStdErr: os.Stderr,
	}

	os.Stdout = o.out
	os.Stderr = o.err

	return &o
}

type OutputBandit struct {
	sync.RWMutex

	out *os.File
	err *os.File

	origStdOut *os.File
	origStdErr *os.File
}

func (o *OutputBandit) Close() {
	o.Lock()
	defer o.Unlock()

	o.out.Close()
	o.err.Close()

	os.Stdout = o.origStdOut
	os.Stderr = o.origStdErr
}
