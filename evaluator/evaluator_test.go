package evaluator

import (
	"github.com/jacksonopp/monkey/lexer"
	"github.com/jacksonopp/monkey/object"
	"github.com/jacksonopp/monkey/parser"
	"testing"
)

func TestExpressions(t *testing.T) {
	t.Run("integer expression", func(t *testing.T) {
		tests := []struct {
			input    string
			expected int64
		}{
			{
				"5",
				5,
			},
			{
				"10",
				10,
			},
			{
				"-5",
				-5,
			},
			{
				"5 + 5",
				5 + 5,
			},
			{
				"5 * 5",
				5 * 5,
			},
			{
				"5 + 2 * 2",
				5 + 2*2,
			},
			{
				"50 / 5 - 3",
				50/5 - 3,
			},
			{
				"2 * (3+4)",
				2 * (3 + 4),
			},
			{
				"-30 + 10",
				-30 + 10,
			},
		}

		for _, tt := range tests {
			t.Run(tt.input, func(t *testing.T) {
				evaluated := testEval(tt.input)
				testIntegerObject(t, evaluated, tt.expected)
			})
		}
	})

	t.Run("boolean expressions", func(t *testing.T) {
		tests := []struct {
			input    string
			expected bool
		}{
			{
				"true",
				true,
			},
			{
				"false",
				false,
			},
			{
				"1 < 2",
				true,
			},
			{
				"1 > 2",
				false,
			},
			{
				"1 == 1",
				true,
			},
			{
				"1 != 1",
				false,
			},
			{
				"1 != 2",
				true,
			},
			{
				"1 == 2",
				false,
			},
			{
				"true == true",
				true,
			},
			{
				"false == false",
				true,
			},
			{
				"true == false",
				false,
			},
			{
				"true != false",
				true,
			},
			{
				"false != true",
				true,
			},
			{
				"(1 < 2) == true",
				true,
			},
			{
				"(1 > 2) == true",
				false,
			}, {
				"(1 < 2) == false",
				false,
			},
			{
				"(1 > 2) == false",
				true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.input, func(t *testing.T) {
				evaluated := testEval(tt.input)
				testBooleanObject(t, evaluated, tt.expected)
			})
		}
	})

	t.Run("bang operator", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected bool
		}{
			{
				"not true",
				"!true",
				false,
			},
			{
				"not false",
				"!false",
				true,
			},
			{
				"not 5",
				"!5",
				false,
			},
			{
				"not not true",
				"!!true",
				true,
			},
			{
				"not not false",
				"!!false",
				false,
			},
			{
				"not not 5",
				"!!5",
				true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				evaluated := testEval(tt.input)
				testBooleanObject(t, evaluated, tt.expected)
			})
		}
	})

	t.Run("if else expressions", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected interface{}
		}{
			{
				"true",
				"if (true) { 10 }",
				10,
			},
			{
				"false",
				"if (false) { 10 }",
				nil,
			},
			{
				"truthy",
				"if (1) { 10 }",
				10,
			},
			{
				"comparison true",
				"if (1 < 10) { 10 }",
				10,
			},
			{
				"comparison false",
				"if (1 > 10) { 10 }",
				nil,
			},
			{
				"comparison else is true",
				"if (1 > 10) { 10 } else { 20 }",
				20,
			},
			{
				"comparison false",
				"if (1 < 10) { 10 } else { 20 }",
				10,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				evaluated := testEval(tt.input)
				integer, ok := tt.expected.(int)
				if ok {
					testIntegerObject(t, evaluated, int64(integer))
				} else {
					testNullObject(t, evaluated)
				}
			})
		}
	})

	t.Run("function expression", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected int64
		}{
			{
				"identity",
				"let identity = fn(x) { x; }; identity(5);",
				5,
			},
			{
				"identity with return",
				"let identity = fn(x) { return x; }; identity(5);",
				5,
			},
			{
				"double",
				"let double = fn(x) { x * 2; }; double(5);",
				10,
			},
			{
				"add",
				"let add = fn(x, y) { x + y; }; add(5, 5);",
				10,
			},
			{
				"add with calling add",
				"let add = fn(x, y) { x + y; }; add(5, add(5, 5));",
				15,
			},
			{
				"anonymous iife",
				"fn(x){ x; }(5);",
				5,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				testIntegerObject(t, testEval(tt.input), tt.expected)
			})
		}
	})
}

func TestStatements(t *testing.T) {
	t.Run("return statement", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected int64
		}{
			{"basic return", "return 10;", 10},
			{"not post return", "return 10; 9;", 10},
			{"not post return, with calculation", "return 2 * 5; 9;", 10},
			{"not pre return or post return", "9; return 2 * 5; 9;", 10},
			{
				"inside a nested statement",
				`
				if (10 > 1) {
				  if (10 > 1) {
					return 10;
				  }
				  return 1;
				}
				`,
				10,
			},
		}

		for _, tt := range tests {
			t.Run(tt.input, func(t *testing.T) {
				evaluated := testEval(tt.input)
				testIntegerObject(t, evaluated, tt.expected)
			})
		}
	})

	t.Run("let statement", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected int64
		}{
			{
				"simple assignment",
				"let a = 10; a;",
				10,
			},
			{
				"simple assignment with calculation",
				"let a = 5 * 2; a;",
				10,
			},
			{
				"multiple assignments",
				"let a = 10; let b = a; b;",
				10,
			},
			{
				"multiple assignments with calculations",
				"let a = 5; let b = 3; let c = a * b - 5; c;",
				10,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				testIntegerObject(t, testEval(tt.input), tt.expected)
			})
		}
	})
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) {x + 2};"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)

	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. got=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestErrorHandling(t *testing.T) {
	t.Run("error handling", func(t *testing.T) {
		tests := []struct {
			name            string
			input           string
			expectedMessage string
		}{
			{
				"type mismatch",
				"5 + true",
				"type mismatch: INTEGER + BOOLEAN",
			},
			{
				"adding int and bool",
				"5 + true; 8;",
				"type mismatch: INTEGER + BOOLEAN",
			},
			{
				"negative true",
				"-true",
				"unknown operator: -BOOLEAN",
			},
			{
				"adding bools",
				"true + true",
				"unknown operator: BOOLEAN + BOOLEAN",
			},
			{
				"adding bools nested",
				"if (20 > 1) { true + false }",
				"unknown operator: BOOLEAN + BOOLEAN",
			},
			{
				"adding bools deep nested",
				`
				if (10 > 1) {
				  if (10 > 1) {
					return true + false;
				  }
				return 1;
				}
				`,
				"unknown operator: BOOLEAN + BOOLEAN",
			},
			{
				"identifier not found",
				"foobar",
				"identifier not found: foobar",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				evaluated := testEval(tt.input)

				errObj, ok := evaluated.(*object.Error)
				if !ok {
					t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
				}
				if errObj.Message != tt.expectedMessage {
					t.Errorf("wrong error message. want=%q, got=%q", tt.expectedMessage, errObj.Message)
				}
			})
		}
	})
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not null. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. want=%d, got=%d", expected, result.Value)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. want=%t, got=%t", expected, result.Value)
		return false
	}
	return true
}
