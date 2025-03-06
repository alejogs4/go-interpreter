package parser

import (
	"fmt"
	"strconv"
	"testing"

	"go-interpreter.com/m/ast"
	"go-interpreter.com/m/lexer"
	"go-interpreter.com/m/token"
)

func TestLetStatements(t *testing.T) {
	testCases := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tc := range testCases {

		l := lexer.New(tc.input)
		p := New(l)
		program := p.ParseProgram()
		checkErrors(t, p)

		if program == nil {
			t.Fatalf("ParseProgram() returned nil\n")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tc.expectedIdentifier) {
			return
		}

		letStatement := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, letStatement, tc.expectedValue) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, stm := range program.Statements {
		returnStm, ok := stm.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%T", stm)
			continue
		}

		if returnStm.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
				returnStm.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
}

func TestIntegerLiteral(t *testing.T) {
	input := "5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stm, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	integerExpression, ok := stm.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stm.Expression)
	}

	if integerExpression.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, integerExpression.Value)
	}

	if integerExpression.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5",
			integerExpression.TokenLiteral())
	}

}

func TestPrefixIntegerOperators(t *testing.T) {
	testCases := []struct {
		input          string
		expectedPrefix string
		expression     interface{}
	}{
		{"!5", "!", 5},
		{"-10", "-", 10},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range testCases {
		l := lexer.New(tt.input)
		parser := New(l)
		program := parser.ParseProgram()

		checkErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("Program prefix operators expected statements %d, got %d", 1, len(program.Statements))
		}

		stm, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Program prefix operators expected expression statement for input %s", tt.input)
		}

		exp, ok := stm.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("Program prefix operators expected prefix expression for input %s", tt.input)
		}

		if exp.Operator != tt.expectedPrefix {
			t.Fatalf("Program prefix operators expected prefix operator %s, got %s", tt.expectedPrefix, exp.Operator)
		}

		if ok := testLiteralExpression(t, exp.Right, tt.expression); !ok {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}
		testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
	}
}

func TestBooleanExpressions(t *testing.T) {
	input := "true"
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	checkErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not a expression statement got=%s\n", stmt.TokenLiteral())
	}

	booleanExpression, ok := stmt.Expression.(*ast.Boolean)
	if !ok {
		t.Fatalf("Expected boolean expression got=%v\n", stmt.Token.Type)
	}

	if booleanExpression.Value != true {
		t.Fatalf("Expected boolean value=true got=%s\n", strconv.FormatBool(booleanExpression.Value))
	}
}

func TestIfExpression(t *testing.T) {
	input := "if (x < y) { x }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}
	exp, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	ifExp := testIfExpression(
		t,
		exp,
		func(t *testing.T, ifCond ast.Expression) bool {
			return testInfixExpression(t, ifCond, "x", "<", "y")
		},
		func(t *testing.T, consequence *ast.ExpressionStatement) bool {
			return testIdentifier(t, consequence.Expression, "x")
		},
	)

	if ifExp == nil {
		return
	}

	if ifExp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", ifExp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := "if (x < y) { x } else { y }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkErrors(t, p)

	exp, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	ifExp := testIfExpression(
		t,
		exp,
		func(t *testing.T, ifCond ast.Expression) bool {
			return testInfixExpression(t, ifCond, "x", "<", "y")
		},
		func(t *testing.T, consequence *ast.ExpressionStatement) bool {
			return testIdentifier(t, consequence.Expression, "x")
		},
	)

	if ifExp == nil {
		return
	}

	if ifExp.Alternative == nil {
		t.Fatalf("exp.Alternative.Statements was nil")
	}

	if len(ifExp.Alternative.Statements) != 1 {
		t.Fatalf("ifExp.Alternative does not contain %d statements. got=%d\n",
			1, len(ifExp.Alternative.Statements))
	}

	alternativeStm, ok := ifExp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ifExp.Alternative.Statements[0] is not ast.ExpressionStatement. got=%T",
			ifExp.Alternative.Statements[0])
	}

	alternativeExp, ok := alternativeStm.Expression.(*ast.Identifier)
	if !ok {
		t.Fatal("ifExp.Alternative is not ast.Idenfier")
	}

	if !testIdentifier(t, alternativeExp, "y") {
		return
	}
}

func TestFunctionLiterals(t *testing.T) {
	input := `fn(x, y) { return x + y; }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	exp, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	fnExp, ok := exp.Expression.(*ast.FunctionExpression)
	if !ok {
		t.Fatalf("exp.Expression is not ast.FunctionExpression. got=%T",
			exp.Expression)
	}

	if fnExp.Token.Type != token.Function {
		t.Fatalf("fnExp.Token is not a function got=%s",
			fnExp.Token.Type)
	}
}

func TestFunctionParametersParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionExpression)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(function.Parameters))
		}

		for i, parsedParam := range function.Parameters {
			testLiteralExpression(t, parsedParam, tt.expectedParams[i])
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2*3, 4+5);`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func checkErrors(t *testing.T, p *Parser) {
	if len(p.Errors) == 0 {
		return
	}

	t.Error(p.Error())
	t.FailNow()
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

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name)
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.TokenLiteral())
		return false
	}
	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, v)
	case int64:
		return testIntegerLiteral(t, exp, int(v))
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testInfixExpression(
	t *testing.T,
	exp ast.Expression,
	left interface{},
	operator string,
	right interface{},
) bool {
	infixExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, infixExp.Left, left) {
		return false
	}

	if infixExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, infixExp.Operator)
		return false
	}

	if !testLiteralExpression(t, infixExp.Right, right) {
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	boolExp, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if boolExp.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, boolExp.Value)
		return false
	}

	if boolExp.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			value, boolExp.TokenLiteral())
		return false
	}

	return true
}

type conditionExpression func(t *testing.T, ifExp ast.Expression) bool
type conditionConsequenceExpression func(t *testing.T, consequence *ast.ExpressionStatement) bool

func testIfExpression(t *testing.T, exp *ast.ExpressionStatement, cond conditionExpression, condConsequence conditionConsequenceExpression) *ast.IfExpression {
	ifExp, ok := exp.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
			exp.Expression)
	}

	if ok := cond(t, ifExp.Condition); !ok {
		return nil
	}

	if len(ifExp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(ifExp.Consequence.Statements))
	}

	consequence, ok := ifExp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			ifExp.Consequence.Statements[0])
	}

	if !condConsequence(t, consequence) {
		return nil
	}

	return ifExp
}
