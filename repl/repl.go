package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/icholy/monkey/lexer"
	"github.com/icholy/monkey/token"
)

var Prefix = ">> "

func Run(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	fmt.Fprint(out, Prefix)
	for scanner.Scan() {
		line := scanner.Text()
		lex := lexer.New(line)
		for {
			tok := lex.NextToken()
			if tok.Type == token.ILLEGAL || tok.Type == token.EOF {
				break
			}
			fmt.Fprintln(out, tok)
		}
		fmt.Fprint(out, Prefix)
	}
}
