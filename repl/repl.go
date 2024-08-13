package repl

import (
	"bufio"
	"fmt"
	"github.com/jacksonopp/monkey/evaluator"
	"github.com/jacksonopp/monkey/lexer"
	"github.com/jacksonopp/monkey/object"
	"github.com/jacksonopp/monkey/parser"
	"github.com/jacksonopp/monkey/token"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer, flags [3]*bool) {
	_, astFlag, tokenFlag := *flags[0], *flags[1], *flags[2]

	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		io.WriteString(out, PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()

		if line == ":exit" {
			io.WriteString(out, "Bye!\n")
			return
		}

		l := lexer.New(line)

		if tokenFlag {
			for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
				t := fmt.Sprintf("%+v\n", tok)
				io.WriteString(out, t)
			}
			continue
		}

		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		if astFlag {
			io.WriteString(out, program.String())
			io.WriteString(out, "\n")
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
