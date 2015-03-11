# oBandit [![GoDoc](https://godoc.org/github.com/missionMeteora/oBandit?status.svg)](https://godoc.org/github.com/missionMeteora/oBandit) ![Status](https://img.shields.io/badge/status-beta-yellow.svg)

oBandit is a library which will hijack your Stdout and Stderr and replace each with an output file.

## Usage
``` go
package main

import (
	"fmt"
	"log"

	"github.com/missionMeteora/oBandit"
)

func main() {
	b, err := bandit.New(`output.txt`, `error.txt`)
	if err != nil {
		log.Println("There was an issue creating a new Bandit:", err)
	}

	fmt.Println("Here is some stdout")
	log.Println("Here is some stderr")

	b.Close()

	fmt.Println("Stdout back in town!")
	log.Println("Stderr back in town!")
}
```