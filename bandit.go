package oBandit

import (
	"os"
	"syscall"
)

// Creates and returns a new instance of Output Bandit:
// 		- Opens both an output and error file based on the locations provided
// 		- Sets reference to original Stdout and Stderr files
// 		- Sets file descriptor values for original Stdout and Stderr
// 		- Calls set files, set tmp files, and hijack
// 		- Returns pointer to Bandit struct
func New(outLoc, errLoc string) (*Bandit, error) {
	o := Bandit{
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

type Bandit struct {
	out   *os.File
	err   *os.File
	outFd int
	errFd int

	outOrig   *os.File
	errOrig   *os.File
	outOrigFd int
	errOrigFd int

	outTmpLoc string
	errTmpLoc string
	outTmp    *os.File
	errTmp    *os.File
	outTmpFd  int
	errTmpFd  int
}

// Sets the files which will be used as the new Stdout and Stderr:
// 		- Creates out and err files based on provided locs
// 		- Sets a reference to the out and err files
// 		- Sets the file descriptor value for the out and err files
func (o *Bandit) setFiles(outLoc, errLoc string) error {
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

// Sets the tmp files which will be used as the file descriptor placeholders:
// 		- Creates outTmp and errTmp files which are placed in the TempDir
// 		- Sets a reference to the out and err files
// 		- Sets the file descriptor value for the outTmp and errTmp files
func (o *Bandit) setTmpFiles() error {
	outTmpLoc := os.TempDir() + ".bndtOut"
	errTmpLoc := os.TempDir() + ".bndtErr"

	outTmpF, err := os.Create(outTmpLoc)
	if err != nil {
		return err
	}

	errTmpF, err := os.Create(errTmpLoc)
	if err != nil {
		return err
	}

	o.outTmpLoc = outTmpLoc
	o.errTmpLoc = errTmpLoc
	o.outTmp = outTmpF
	o.errTmp = errTmpF

	return nil
}

// Hijacks the Stdout and Stderr:
//		- Sets the orig Stdout file descriptor value to the outTmp file descriptor value
// 		- Sets the new out file descriptor value to the original Stdout file descriptor value
//		- Sets the orig Stderr file descriptor value to the errTmp file descriptor value
// 		- Sets the new err file descriptor value to the original Stderr file descriptor value
func (o *Bandit) hijack() {
	syscall.Dup2(o.outOrigFd, o.outTmpFd)
	syscall.Dup2(o.outFd, o.outOrigFd)

	syscall.Dup2(o.errOrigFd, o.errTmpFd)
	syscall.Dup2(o.errFd, o.errOrigFd)
}

// Un-Hijacks the Stdout and Stderr:
//		- Sets the orig Stdout file descriptor value to its original value
// 		- Sets the new out file descriptor value to its original value
//		- Sets the orig Stderr file descriptor value to its original value
// 		- Sets the new err file descriptor value to its original value
func (o *Bandit) unhijack() {
	syscall.Dup2(o.outOrigFd, o.outFd)
	syscall.Dup2(o.outTmpFd, o.outOrigFd)

	syscall.Dup2(o.errOrigFd, o.errFd)
	syscall.Dup2(o.errTmpFd, o.errOrigFd)
}

// Closes Output Bandit:
// 		- Sets global Stdout and Stderr back to the original
// 		- Sets the output for the log package to the original Stderr
// 		- Closes the output and error file
func (o *Bandit) Close() {
	o.unhijack()
	o.out.Close()
	o.err.Close()
	o.outTmp.Close()
	o.errTmp.Close()

	os.Remove(o.outTmpLoc)
	os.Remove(o.errTmpLoc)
}
