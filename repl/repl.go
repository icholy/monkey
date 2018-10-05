package repl

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/icholy/monkey/evaluator"
	"github.com/icholy/monkey/object"
	"github.com/icholy/monkey/parser"
)

var Prefix = ">> "

func Run(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnv(nil)
	fmt.Fprint(out, Prefix)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "exit" {
			break
		}
		program, err := parser.Parse(line)
		if err != nil {
			fmt.Println(err)
		} else {
			obj, err := evaluator.Eval(program, env)
			if err != nil {
				fmt.Fprintf(out, "ERROR: %s\n", err)
			} else {
				fmt.Fprintln(out, obj.Inspect())
			}
		}
		fmt.Fprint(out, Prefix)
	}
}

func Exec(in io.Reader) error {
	data, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}
	program, err := parser.Parse(string(data))
	if err != nil {
		return err
	}
	_, err = evaluator.Eval(program, object.NewEnv(nil))
	return err
}
