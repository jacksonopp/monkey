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
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
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
