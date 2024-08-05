package parser

import (
	"github.com/jacksonopp/monkey/ast"
	"github.com/jacksonopp/monkey/lexer"
	"testing"
)

type expectedTest struct {
	identifier string
}

func TestLetStatement(t *testing.T) {
	t.Run("basic let statements", func(t *testing.T) {
		input := `
let x = 5;
let y = 10;
let foobar = 838383;
`

		l := lexer.New(input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErorrs(t, p)

		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}
		if len(program.Statements) != 3 {
			t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
		}

		tests := []expectedTest{
			{"x"},
			{"y"},
			{"foobar"},
		}

		for i, tt := range tests {
			stmt := program.Statements[i]
			if !assertLetStatement(t, stmt, tt.identifier) {
				return
			}
		}
	})
	t.Run("expect parser errors", func(t *testing.T) {
		input := `
	let x 5;
	let = 10;
	let 4234234;
`

		l := lexer.New(input)
		p := New(l)

		p.ParseProgram()

		errors := p.Errors()
		if len(errors) != 3 {
			t.Errorf("expected 3 errors, got=%d", len(errors))
			for _, msg := range errors {
				t.Errorf("error: %s", msg)
			}
		}
	})
}

func TestReturnStatement(t *testing.T) {
	t.Run("basic return statement", func(t *testing.T) {
		input := `
return 5;
return 10;
return 423432;`
		l := lexer.New(input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErorrs(t, p)

		if len(program.Statements) != 3 {
			t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
		}

		for _, stmt := range program.Statements {
			returnStmt, ok := stmt.(*ast.ReturnStatement)
			if !ok {
				t.Errorf("stmt not *ast.ReturnStatement. got %T", stmt)
				continue
			}
			if returnStmt.TokenLiteral() != "return" {
				t.Errorf("returnStmt.TokenLiteral not 'return'. got=%q", returnStmt.TokenLiteral())
			}
		}
	})
}

func assertLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func checkParserErorrs(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error %q", msg)
		t.FailNow()
	}
}
