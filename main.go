package outputBandit

import (
	"os"
)

func New(outLoc, errLoc string) *OutputBandit {
	out, _ := os.Create(outLoc)
	err, _ := os.Create(errLoc)

	os.Stdout = out
	os.Stderr = err

	o := OutputBandit{
		out: out,
		err: err,
	}

	return &o
}

type OutputBandit struct {
	out *os.File
	err *os.File
}

func (o *OutputBandit) Close() {
	o.out.Close()
	o.err.Close()
}
