package parser

import (
	"fmt"
	"github.com/jacksonopp/monkey/ast"
	"github.com/jacksonopp/monkey/lexer"
	"testing"
)

type expectedTest struct {
	identifier string
}

func TestStatements(t *testing.T) {
	t.Run("let statement", func(t *testing.T) {
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
			if len(errors) < 3 {
				t.Errorf("expected at least 3 errors, got=%d", len(errors))
				for _, msg := range errors {
					t.Errorf("error: %s", msg)
				}
			}
		})
	})

	t.Run("return statement", func(t *testing.T) {
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

func TestExpressions(t *testing.T) {
	t.Run("identifier expression", func(t *testing.T) {
		input := "foobar;"
		l := lexer.New(input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErorrs(t, p)

		checkProgramStatmentsLength(t, program.Statements, 1)
		stmt := checkStatmentIsExpressionStatement(t, program.Statements[0])

		ident, ok := stmt.Expression.(*ast.Identifier)
		if !ok {
			t.Fatalf("expr not *ast.Identifier. got=%T", stmt.Expression)
		}
		if ident.Value != "foobar" {
			t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
		}
		if ident.TokenLiteral() != "foobar" {
			t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
		}
	})

	t.Run("integer literal expression", func(t *testing.T) {
		input := "5;"

		l := lexer.New(input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErorrs(t, p)
		stmt := checkStatmentIsExpressionStatement(t, program.Statements[0])

		testIntegerLiteral(t, stmt.Expression, 5)
		checkProgramStatmentsLength(t, program.Statements, 1)
	})

	t.Run("prefix operator expressions", func(t *testing.T) {
		prefixTests := []struct {
			input        string
			operator     string
			integerValue int64
		}{
			{"!5", "!", 5},
			{"-15", "-", 15},
		}

		for _, tt := range prefixTests {
			t.Run(fmt.Sprintf("%s operator", tt.operator), func(t *testing.T) {
				l := lexer.New(tt.input)
				p := New(l)
				program := p.ParseProgram()
				checkParserErorrs(t, p)
				checkProgramStatmentsLength(t, program.Statements, 1)
				stmt := checkStatmentIsExpressionStatement(t, program.Statements[0])

				exp, ok := stmt.Expression.(*ast.PrefixExpression)
				if !ok {
					t.Fatalf("stmt.Expression is not ast.PrefixExpression. got=%T", stmt.Expression)
				}
				if exp.Operator != tt.operator {
					t.Fatalf("exp.Operator is not %s. got %s", tt.operator, exp.Operator)
				}
				if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
					return
				}
			})
		}
	})
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%s", value, integ.TokenLiteral())
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got %s", value, integ.TokenLiteral())
		return false
	}

	return true
}

func checkProgramStatmentsLength(t *testing.T, stmts []ast.Statement, length int) {
	if len(stmts) != length {
		t.Errorf("program has not enough statements. got=%d", len(stmts))
		t.FailNow()
	}
}

func checkStatmentIsExpressionStatement(t *testing.T, stmt ast.Statement) *ast.ExpressionStatement {
	s, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("statement is not ast.ExpressionStatement. got=%T", stmt)
		t.FailNow()
	}
	return s
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
