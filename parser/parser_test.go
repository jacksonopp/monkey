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
			checkParserErrors(t, p)

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
				if !testLetStatement(t, stmt, tt.identifier) {
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
		t.Run("let statements with expressions", func(t *testing.T) {
			tests := []struct {
				name               string
				input              string
				expectedIdentifier string
				expectedValue      interface{}
			}{
				{
					"assign number",
					"let x = 5;",
					"x",
					5,
				},
				{
					"assign bool",
					"let y = true;",
					"y",
					true,
				},
				{
					"assign variable",
					"let foobar = z;",
					"foobar",
					"z",
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					l := lexer.New(tt.input)
					p := New(l)
					program := p.ParseProgram()
					checkParserErrors(t, p)

					checkProgramStatementsLength(t, program.Statements, 1)

					stmt := program.Statements[0]

					if !testLetStatement(t, stmt, tt.expectedIdentifier) {
						return
					}

					val := stmt.(*ast.LetStatement).Value
					if !testLiteralExpression(t, val, tt.expectedValue) {
						return
					}
				})
			}
		})
	})

	t.Run("return statement", func(t *testing.T) {
		t.Run("basic return statement", func(t *testing.T) {
			input := `
return 5;
return 10;
return 123123;`
			l := lexer.New(input)
			p := New(l)

			program := p.ParseProgram()
			checkParserErrors(t, p)

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

		t.Run("return statement with expressions", func(t *testing.T) {
			tests := []struct {
				input         string
				expectedValue interface{}
			}{
				{"return 5;", 5},
				{"return true;", true},
				{"return foobar;", "foobar"},
			}

			for _, tt := range tests {
				t.Run(tt.input, func(t *testing.T) {
					l := lexer.New(tt.input)
					p := New(l)
					program := p.ParseProgram()
					checkParserErrors(t, p)

					if len(program.Statements) != 1 {
						t.Fatalf("program.Statements does not contain 1 statements. got=%d",
							len(program.Statements))
					}

					stmt := program.Statements[0]
					returnStmt, ok := stmt.(*ast.ReturnStatement)
					if !ok {
						t.Fatalf("stmt not *ast.ReturnStatement. got=%T", stmt)
					}
					if returnStmt.TokenLiteral() != "return" {
						t.Fatalf("returnStmt.TokenLiteral not 'return', got %q",
							returnStmt.TokenLiteral())
					}
					if testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue) {
						return
					}

				})
			}
		})

		// FIXME This test causes suite to hang?
		//t.Run("return statment with expressions", func(t *testing.T) {
		//	tests := []struct {
		//		name          string
		//		input         string
		//		expectedValue interface{}
		//	}{
		//		{
		//			"return number",
		//			"return 5;",
		//			5,
		//		},
		//		{
		//			"return bool",
		//			"return true;",
		//			true,
		//		},
		//		{
		//			"return variable",
		//			"return foobar",
		//			"foobar",
		//		},
		//	}
		//
		//	for _, tt := range tests {
		//		t.Run(tt.name, func(t *testing.T) {
		//			l := lexer.New(tt.input)
		//			p := New(l)
		//			program := p.ParseProgram()
		//
		//			checkParserErrors(t, p)
		//			checkProgramStatementsLength(t, program.Statements, 1)
		//			stmt := program.Statements[0]
		//
		//			returnStmt, ok := stmt.(*ast.ReturnStatement)
		//			if !ok {
		//				t.Fatalf("stmt not *ast.ReturnStatement. got=%T", stmt)
		//			}
		//			if returnStmt.TokenLiteral() != "return" {
		//				t.Fatalf("returnStmtn.TokenLiteral not 'return'. got=%q", returnStmt.TokenLiteral())
		//			}
		//
		//			if testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue) {
		//				return
		//			}
		//		})
		//	}
		//})
	})
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
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
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{"basic identifier", "foobar;", "foobar"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {

				l := lexer.New(tt.input)
				p := New(l)
				program := p.ParseProgram()
				checkParserErrors(t, p)

				checkProgramStatementsLength(t, program.Statements, 1)
				stmt := checkStatementIsExpressionStatement(t, program.Statements[0])

				testIdentifier(t, stmt.Expression, tt.expected)
			})
		}
	})

	t.Run("integer literal expression", func(t *testing.T) {
		input := "5;"

		l := lexer.New(input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		stmt := checkStatementIsExpressionStatement(t, program.Statements[0])

		testIntegerLiteral(t, stmt.Expression, 5)
		checkProgramStatementsLength(t, program.Statements, 1)
	})

	t.Run("prefix operator expressions", func(t *testing.T) {
		prefixTests := []struct {
			name     string
			input    string
			operator string
			value    interface{}
		}{
			{"not integer", "!5", "!", 5},
			{"minus integer", "-15", "-", 15},
			{"not true", "!true", "!", true},
			{"not false", "!false", "!", false},
		}

		for _, tt := range prefixTests {
			t.Run(tt.name, func(t *testing.T) {
				l := lexer.New(tt.input)
				p := New(l)
				program := p.ParseProgram()
				checkParserErrors(t, p)
				checkProgramStatementsLength(t, program.Statements, 1)
				stmt := checkStatementIsExpressionStatement(t, program.Statements[0])

				exp, ok := stmt.Expression.(*ast.PrefixExpression)
				if !ok {
					t.Fatalf("stmt.Expression is not ast.PrefixExpression. got=%T", stmt.Expression)
				}
				if exp.Operator != tt.operator {
					t.Fatalf("exp.Operator is not %s. got %s", tt.operator, exp.Operator)
				}
				testLiteralExpression(t, exp.Right, tt.value)
			})
		}
	})

	t.Run("infix operator expressions", func(t *testing.T) {
		tests := []struct {
			name       string
			input      string
			leftValue  interface{}
			operator   string
			rightValue interface{}
		}{
			{"add", "5 + 5;", 5, "+", 5},
			{"sub", "5 - 5;", 5, "-", 5},
			{"mult", "5 * 5;", 5, "*", 5},
			{"div", "5 / 5;", 5, "/", 5},
			{"lt", "5 < 5;", 5, "<", 5},
			{"gt", "5 > 5;", 5, ">", 5},
			{"eq", "5 == 5;", 5, "==", 5},
			{"neq", "5 != 5;", 5, "!=", 5},
			{"eq bools true", "true == true;", true, "==", true},
			{"eq bools false", "false == false;", false, "==", false},
			{"neq bools", "false != true;", false, "!=", true},
		}

		for _, tt := range tests {
			t.Run(fmt.Sprintf("%s infix operator", tt.operator), func(t *testing.T) {
				l := lexer.New(tt.input)
				p := New(l)
				program := p.ParseProgram()
				checkParserErrors(t, p)

				checkProgramStatementsLength(t, program.Statements, 1)
				stmt := checkStatementIsExpressionStatement(t, program.Statements[0])

				testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
			})
		}
	})

	t.Run("operator precedence", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{
				"prefix before mult",
				"-a * b",
				"((-a) * b)",
			},
			{
				"prefix before bang",
				"!-a",
				"(!(-a))",
			},
			{
				"left addition first",
				"a + b + c",
				"((a + b) + c)",
			},
			{
				"left addition before subtraction",
				"a + b - c",
				"((a + b) - c)",
			},
			{
				"left mult first",
				"a * b * c",
				"((a * b) * c)",
			},
			{
				"left mult before div",
				"a * b / c",
				"((a * b) / c)",
			},
			{
				"products before sums",
				"a + b * c + d / e - f",
				"(((a + (b * c)) + (d / e)) - f)",
			},
			{
				"two statements",
				"3 + 4;-5 * 5",
				"(3 + 4)((-5) * 5)",
			},
			{
				"comparisons eq",
				"5 > 4 == 3 < 4",
				"((5 > 4) == (3 < 4))",
			},
			{
				"comparisons neq",
				"5 > 4 != 3 < 4",
				"((5 > 4) != (3 < 4))",
			},
			{
				"kitchen sink",
				"3 + 4 * 5 == 3 * 1 + 4 * 5",
				"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
			},
			{
				"true",
				"true",
				"true",
			},
			{
				"false",
				"false",
				"false",
			},
			{
				"comp before false",
				"3 > 5 == false",
				"((3 > 5) == false)",
			},
			{
				"comp before true",
				"3 < 5 == true",
				"((3 < 5) == true)",
			},
			{
				"grouping addition",
				"1 + (2 + 3) + 4",
				"((1 + (2 + 3)) + 4)",
			},
			{
				"grouping mult",
				"(5 + 5) * 2",
				"((5 + 5) * 2)",
			},
			{
				"grouping div",
				"2 / (5 + 5)",
				"(2 / (5 + 5))",
			},
			{
				"grouping with prefix",
				"-(5 + 5)",
				"(-(5 + 5))",
			},
			{
				"grouping bools",
				"!(true == true)",
				"(!(true == true))",
			},
			{
				"functions a",
				"a + add(b * c) + d",
				"((a + add((b * c))) + d)",
			},
			{
				"functions b",
				"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
				"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
			},
			{
				"functions c",
				"add(a + b + c * d / f + g)",
				"add((((a + b) + ((c * d) / f)) + g))",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				l := lexer.New(tt.input)
				p := New(l)
				program := p.ParseProgram()
				checkParserErrors(t, p)

				actual := program.String()
				if actual != tt.expected {
					t.Errorf("expected=%q, got=%q", tt.expected, actual)
				}
			})
		}
	})

	t.Run("if expressions", func(t *testing.T) {
		input := "if (x < y) { x }"

		l := lexer.New(input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		checkProgramStatementsLength(t, program.Statements, 1)
		stmt := checkStatementIsExpressionStatement(t, program.Statements[0])

		exp, ok := stmt.Expression.(*ast.IfExpression)
		if !ok {
			t.Fatalf("stmt.Expression not ast.IfExpression. got=%T", stmt.Expression)
		}

		if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
			return
		}

		checkProgramStatementsLength(t, exp.Consequence.Statements, 1)

		consequence := checkStatementIsExpressionStatement(t, exp.Consequence.Statements[0])
		if !testIdentifier(t, consequence.Expression, "x") {
			return
		}

		if exp.Alternative != nil {
			t.Errorf("exp.Alternative statement was not nil. got=%+v", exp.Alternative)
		}
	})

	t.Run("if else expression", func(t *testing.T) {
		input := "if (x < y) { x } else { y }"

		l := lexer.New(input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		checkProgramStatementsLength(t, program.Statements, 1)
		stmt := checkStatementIsExpressionStatement(t, program.Statements[0])

		exp, ok := stmt.Expression.(*ast.IfExpression)
		if !ok {
			t.Fatalf("stmt.Expression not ast.IfExpression. got=%T", stmt.Expression)
		}

		if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
			return
		}

		checkProgramStatementsLength(t, exp.Consequence.Statements, 1)

		consequence := checkStatementIsExpressionStatement(t, exp.Consequence.Statements[0])
		if !testIdentifier(t, consequence.Expression, "x") {
			return
		}

		if exp.Alternative == nil {
			t.Fatalf("exp.Alternative was nil")
		}

		alternative := checkStatementIsExpressionStatement(t, exp.Alternative.Statements[0])
		if !testIdentifier(t, alternative.Expression, "y") {
			return
		}
	})

	t.Run("function expressions", func(t *testing.T) {
		input := `fn(x, y) { x + y }`
		l := lexer.New(input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		checkProgramStatementsLength(t, program.Statements, 1)

		stmt := checkStatementIsExpressionStatement(t, program.Statements[0])

		function, ok := stmt.Expression.(*ast.FunctionLiteral)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", program.Statements[0])
		}

		if len(function.Parameters) != 2 {
			t.Fatalf("function literal parameters wrong. want 2, got=%d", len(function.Parameters))
		}

		testLiteralExpression(t, function.Parameters[0], "x")
		testLiteralExpression(t, function.Parameters[1], "y")

		if len(function.Body.Statements) != 1 {
			t.Fatalf("function.Body.Statements has not 1 statement. got=%d", function.Body.Statements)
		}

		bodyStmt := checkStatementIsExpressionStatement(t, function.Body.Statements[0])

		testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
	})

	t.Run("function parameter", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected []string
		}{
			{
				"empty",
				"fn(){};",
				[]string{},
			},
			{
				"one param",
				"fn(x) {};",
				[]string{"x"},
			},
			{
				"three params",
				"fn(x, y, z) {};",
				[]string{"x", "y", "z"},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				l := lexer.New(tt.input)
				p := New(l)
				program := p.ParseProgram()
				checkParserErrors(t, p)

				stmt := checkStatementIsExpressionStatement(t, program.Statements[0])
				fn, ok := stmt.Expression.(*ast.FunctionLiteral)
				if !ok {
					t.Fatalf("fn not *ast.FunctionLiteral. got=%T", stmt.Expression)
				}
				if len(fn.Parameters) != len(tt.expected) {
					t.Fatalf("fn.Parameters has incorrect number of parameters. want=%d, got=%d", len(fn.Parameters), len(tt.expected))
				}

				for i, ident := range tt.expected {
					testLiteralExpression(t, fn.Parameters[i], ident)
				}
			})
		}
	})

	t.Run("call expression", func(t *testing.T) {
		input := "add(1, 2 * 3, 4 + 5)"
		l := lexer.New(input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)
		checkProgramStatementsLength(t, program.Statements, 1)

		stmt := checkStatementIsExpressionStatement(t, program.Statements[0])

		exp, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt.Expression not *ast.CallExpression. got=%T", stmt.Expression)
		}

		if !testIdentifier(t, exp.Function, "add") {
			return
		}

		if len(exp.Arguments) != 3 {
			t.Fatalf("wrong length of arguments. want=%d, got=%d", 3, len(exp.Arguments))
		}

		testLiteralExpression(t, exp.Arguments[0], 1)
		testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
		testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
	})

	t.Run("call expression parameter", func(t *testing.T) {
		tests := []struct {
			name          string
			input         string
			expectedIdent string
			expectedArgs  []string
		}{
			{
				"empty",
				"add()",
				"add",
				[]string{},
			},
			{
				"one param",
				"log(1)",
				"log",
				[]string{"1"},
			},
			{
				"multiple params",
				"foo(1, 2 * 3, 4 + 5)",
				"foo",
				[]string{"1", "(2 * 3)", "(4 + 5)"},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				l := lexer.New(tt.input)
				p := New(l)
				program := p.ParseProgram()

				checkParserErrors(t, p)

				checkProgramStatementsLength(t, program.Statements, 1)
				stmt := checkStatementIsExpressionStatement(t, program.Statements[0])

				exp, ok := stmt.Expression.(*ast.CallExpression)
				if !ok {
					t.Fatalf("stmt.Expression not *ast.CallExpression. got=%T", stmt.Expression)
				}

				testIdentifier(t, exp.Function, tt.expectedIdent)

				if len(exp.Arguments) != len(tt.expectedArgs) {
					t.Fatalf("wrong number of arguments. want=%d, got=%d", len(tt.expectedArgs), len(exp.Arguments))
				}

				for i, arg := range tt.expectedArgs {
					if exp.Arguments[i].String() != arg {
						t.Errorf("argument %d wrong. want=%q, got=%q", i, arg, exp.Arguments[i].String())
					}
				}
			})
		}
	})

}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got %T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}

	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}
	if bo.Value != value {
		t.Errorf("bo.Value not %t. got %t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got %s", value, bo.TokenLiteral())
		return false
	}
	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	op, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not InfixExpression. got %T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, op.Left, left) {
		return false
	}
	if op.Operator != operator {
		t.Errorf("exp.Operator is not '%s', got=%q", operator, op.Operator)
		return false
	}
	if !testLiteralExpression(t, op.Right, right) {
		return false
	}
	return true
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

func checkProgramStatementsLength(t *testing.T, stmts []ast.Statement, length int) {
	if len(stmts) != length {
		t.Errorf("program has not enough statements. got=%d", len(stmts))
		t.FailNow()
	}
}

func checkStatementIsExpressionStatement(t *testing.T, stmt ast.Statement) *ast.ExpressionStatement {
	s, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("statement is not ast.ExpressionStatement. got=%T", stmt)
		t.FailNow()
	}
	return s
}

func checkParserErrors(t *testing.T, p *Parser) {
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
