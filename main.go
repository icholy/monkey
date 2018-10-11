package main

import (
	"log"
	"os"

	"github.com/icholy/monkey/evaluator"
)

func main() {

	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		if err := evaluator.Run(f); err != nil {
			log.Fatal(err)
		}
		return
	}

	evaluator.REPL(os.Stdin, os.Stdout)
}
