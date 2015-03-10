package outputBandit

import (
	"os"
	"syscall"
)

const (
	outTmpLoc = ".bndtOut"
	errTmpLoc = ".bndtErr"
)

// Creates and returns a new instance of Output Bandit:
// 		- Opens both an output and error file based on the locations provided
// 		- Sets Stdout to the output file and Stderr to the error file
// 		- Sets the output for the log package to the new Stderr
func New(outLoc, errLoc string) (*OutputBandit, error) {
	o := OutputBandit{
		outOrig:   os.Stdout,
		errOrig:   os.Stderr,
		outOrigFd: int(os.Stdout.Fd()),
		errOrigFd: int(os.Stderr.Fd()),
	}

	o.setFiles(outLoc, errLoc)
	o.setTmpFiles()
	o.hijack()

	return &o, nil
}

type OutputBandit struct {
	out   *os.File
	err   *os.File
	outFd int
	errFd int

	outOrig   *os.File
	errOrig   *os.File
	outOrigFd int
	errOrigFd int

	outTmp   *os.File
	errTmp   *os.File
	outTmpFd int
	errTmpFd int
}

func (o *OutputBandit) setFiles(outLoc, errLoc string) error {
	outF, err := os.Create(outLoc)
	if err != nil {
		return err
	}

	errF, err := os.Create(errLoc)
	if err != nil {
		return err
	}

	o.out = outF
	o.err = errF

	o.outFd = int(outF.Fd())
	o.errFd = int(errF.Fd())

	return nil
}

func (o *OutputBandit) setTmpFiles() error {
	outTmpF, err := os.Create(outTmpLoc)
	if err != nil {
		return err
	}

	errTmpF, err := os.Create(errTmpLoc)
	if err != nil {
		return err
	}

	o.outTmp = outTmpF
	o.errTmp = errTmpF

	return nil
}

func (o *OutputBandit) hijack() {
	syscall.Dup2(o.outOrigFd, o.outTmpFd)
	syscall.Dup2(o.outFd, o.outOrigFd)

	syscall.Dup2(o.errOrigFd, o.errTmpFd)
	syscall.Dup2(o.errFd, o.errOrigFd)
}

func (o *OutputBandit) unhijack() {
	syscall.Dup2(o.outOrigFd, o.outFd)
	syscall.Dup2(o.outTmpFd, o.outOrigFd)

	syscall.Dup2(o.errOrigFd, o.errFd)
	syscall.Dup2(o.errTmpFd, o.errOrigFd)
}

// Closes Output Bandit:
// 		- Sets global Stdout and Stderr back to the original
// 		- Sets the output for the log package to the original Stderr
// 		- Closes the output and error file
func (o *OutputBandit) Close() {
	o.unhijack()
	o.out.Close()
	o.err.Close()
	o.outTmp.Close()
	o.errTmp.Close()

	os.Remove(outTmpLoc)
	os.Remove(errTmpLoc)
}
