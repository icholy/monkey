package main

import (
	"log"
	"os"

	"github.com/icholy/monkey/repl"
)

func main() {

	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		if err := repl.Exec(f); err != nil {
			log.Fatal(err)
		}
		return
	}

	repl.Run(os.Stdin, os.Stdout)
}
