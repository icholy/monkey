package evaluator

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/chzyer/readline"

	"github.com/icholy/monkey/object"
	"github.com/icholy/monkey/parser"
)

var Prompt = ">> "

func REPL(in io.Reader, out io.Writer) {
	rl, err := readline.New(Prompt)
	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()
	env := object.NewEnv(nil)
	for {
		line, err := rl.Readline()
		if err != nil {
			log.Fatal(err)
		}
		program, err := parser.Parse(line)
		if err != nil {
			fmt.Println(err)
		} else {
			obj, err := Eval(program, env)
			if err != nil {
				fmt.Fprintf(out, "ERROR: %s\n", err)
			} else {
				fmt.Fprintln(out, obj.Inspect())
			}
		}
	}
}

func Run(in io.Reader) error {
	data, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}
	program, err := parser.Parse(string(data))
	if err != nil {
		return err
	}
	_, err = Eval(program, object.NewEnv(nil))
	return err
}