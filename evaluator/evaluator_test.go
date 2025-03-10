package evaluator

import (
	"testing"

	"go-interpreter.com/m/lexer"
	"go-interpreter.com/m/object"
	"go-interpreter.com/m/parser"
)

func TestEvalInteger(t *testing.T) {
	testCases := []struct {
		input    string
		expected int
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
	}

	for _, tc := range testCases {
		evaluated := executeEval(tc.input)
		testIntegerLiteral(t, evaluated, tc.expected)
	}
}

func TestEvalBool(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tc := range testCases {
		evaluated := executeEval(tc.input)
		testBooleanLiteral(t, evaluated, tc.expected)
	}
}

func TestNegationOperator(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tc := range testCases {
		evaluated := executeEval(tc.input)
		testBooleanLiteral(t, evaluated, tc.expected)
	}
}

func TestInfixExpressions(t *testing.T) {
	testCases := []struct {
		input    string
		expected interface{}
	}{
		{"5+5", 10},
		{"17-5", 12},
		{"40/2", 20},
		{"30*2", 60},
		{"40/2 + 21", 41},
		{"5+2*10", 25},
		{"5 == 5", true},
		{"4 == 6", false},
		{"5 > 4", true},
		{"5 > 10", false},
		{"5 < 10", true},
		{"5 < 4", false},
		{"5 != 5", false},
		{"5 != 4", true},
	}

	for _, tc := range testCases {
		evaluated := executeEval(tc.input)
		val, ok := tc.expected.(int)
		valBool, okBool := tc.expected.(bool)

		if !ok && !okBool {
			t.Errorf("expected value is not supoorted, got=%t\n", tc.expected)
			return
		}

		if ok {
			testIntegerLiteral(t, evaluated, val)
			return
		}

		if okBool {
			testBooleanLiteral(t, evaluated, valBool)
			return
		}
	}
}

func executeEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return Eval(program)
}

func testIntegerLiteral(t *testing.T, obj object.Object, expected int) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not boolean. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}

	return true
}
