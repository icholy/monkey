package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/icholy/monkey/evaluator"
	"github.com/icholy/monkey/lexer"
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
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if errs := p.Errors(); len(errs) > 0 {
			for _, err := range errs {
				fmt.Println(err)
			}
		} else {
			obj := evaluator.Eval(program, env)
			fmt.Fprintln(out, obj.Inspect())
		}
		fmt.Fprint(out, Prefix)
	}
}
