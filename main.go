package outputBandit

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

func NewBandit(outLoc, errLoc string) *OutputBandit {
	outF, _ := os.Create(outLoc)
	errF, _ := os.Create(errLoc)

	o := OutputBandit{
		out:        bufio.NewWriter(outF),
		err:        bufio.NewWriter(errF),
		outF:       outF,
		errF:       errF,
		origStdOut: os.Stdout,
		origStdErr: os.Stderr,
	}

	os.Stdout = o.out
	os.Stderr = o.err

	return &o
}

type OutputBandit struct {
	sync.RWMutex

	out *bufio.Writer
	err *bufio.Writer

	outF *os.File
	errF *os.File

	origStdOut *os.Stderr
	origStdErr *os.Stdout
}

func (o *OutputBandit) Close() {
	o.Lock()
	defer o.Unlock()

	o.out.Flush()
	o.err.Flush()

	o.outF.Close()
	o.errF.Close()

	os.Stdout = o.origStdOut
	os.Stderr = o.origStdErr
}
