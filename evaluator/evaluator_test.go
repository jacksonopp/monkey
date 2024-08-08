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
			name     string
			input    string
			expected int64
		}{
			{
				"5",
				"5",
				5,
			},
			{
				"10",
				"10",
				10,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				evaluated := testEval(tt.input)
				testIntegerObject(t, evaluated, tt.expected)
			})
		}
	})

	t.Run("boolean expressions", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected bool
		}{
			{
				"true",
				"true",
				true,
			},
			{
				"false",
				"false",
				false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
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
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
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
