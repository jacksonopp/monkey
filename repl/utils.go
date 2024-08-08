package repl

import "io"

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Whoops! We ran in to some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, e := range errors {
		io.WriteString(out, "\t"+e+"\n")
	}
}
