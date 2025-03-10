package evaluator

import (
	"go-interpreter.com/m/ast"
	"go-interpreter.com/m/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgramStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToObject(node.Value)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node.Operator, Eval(node.Right))
	case *ast.InfixExpression:
		return evalInfixExpression(node.Operator, Eval(node.Left), Eval(node.Right))
	}

	return nil
}

func evalProgramStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range statements {
		result = Eval(stmt)
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evaluateNegationOperator(right)
	case "-":
		return evaluateMinusOperator(right)
	}

	return NULL
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.Integer_Obj && right.Type() == object.Integer_Obj:
		return evaluateArimethic(operator, left, right)
	}

	return NULL
}

type OperateOnInfixOperators[T any] func(left, right T) object.Object

func evaluateArimethic(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case ">":
		return &object.Boolean{Value: leftVal > rightVal}
	case "<":
		return &object.Boolean{Value: leftVal < rightVal}
	case "==":
		return &object.Boolean{Value: leftVal == rightVal}
	case "!=":
		return &object.Boolean{Value: leftVal != rightVal}
	}

	return NULL
}

func evaluateNegationOperator(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evaluateMinusOperator(right object.Object) object.Object {
	if right.Type() != object.Integer_Obj {
		return NULL
	}

	rightVal := right.(*object.Integer).Value
	return &object.Integer{Value: -rightVal}
}

func nativeBoolToObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}

	return FALSE
}
