package outputBandit

import (
	"log"
	"os"
	"syscall"
)

// Creates and returns a new instance of Output Bandit:
// 		- Opens both an output and error file based on the locations provided
// 		- Sets Stdout to the output file and Stderr to the error file
// 		- Sets the output for the log package to the new Stderr
func New(outLoc, errLoc string) (*OutputBandit, error) {
	outF, err := os.Create(outLoc)
	if err != nil {
		return nil, err
	}

	errF, err := os.Create(errLoc)
	if err != nil {
		return nil, err
	}

	o := OutputBandit{
		out: outF,
		err: errF,

		outOrig: os.Stdout,
		errOrig: os.Stderr,
	}

	os.Stdout = outF
	os.Stderr = errF
	log.SetOutput(errF)
	syscall.Dup2(int(errF.Fd()), 2)

	return &o, nil
}

type OutputBandit struct {
	out *os.File
	err *os.File

	outOrig *os.File
	errOrig *os.File
}

// Closes Output Bandit:
// 		- Sets global Stdout and Stderr back to the original
// 		- Sets the output for the log package to the original Stderr
// 		- Closes the output and error file
func (o *OutputBandit) Close() {
	os.Stdout = o.outOrig
	os.Stderr = o.errOrig
	log.SetOutput(o.errOrig)

	o.out.Close()
	o.err.Close()
}
