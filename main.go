package outputBandit

import (
	"log"
	"os"
)

func New(outLoc, errLoc string) *OutputBandit {
	out, _ := os.Create(outLoc)
	err, _ := os.Create(errLoc)

	o := OutputBandit{
		out: out,
		err: err,

		outOrig: os.Stdout,
		errOrig: os.Stderr,
	}

	os.Stdout = out
	os.Stderr = err
	log.SetOutput(err)

	return &o
}

type OutputBandit struct {
	out *os.File
	err *os.File

	outOrig *os.File
	errOrig *os.File
}

func (o *OutputBandit) Close() {
	o.out.Close()
	o.err.Close()

	os.Stdout = o.outOrig
	os.Stderr = o.errOrig
	log.SetOutput(o.errOrig)
}
