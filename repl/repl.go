package repl

import (
	"bufio"
	"fmt"
	"io"

	"go-interpreter.com/m/lexer"
	"go-interpreter.com/m/parser"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors) != 0 {
			printParserErrors(out, p.Errors)
		}

		_, _ = fmt.Fprintf(out, "%+v\n", program.String())
	}
}

func printParserErrors(out io.Writer, errors []error) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg.Error()+"\n")
	}
}
