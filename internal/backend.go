package internal

import (
	"fmt"
)

type RuntimeError struct {
	Token Token
	Msg   string
}

func (r RuntimeError) Error() string {
	return r.Msg
}

type Interpreter struct {
	reporter ErrorReporter
}

func NewInterpreter(reporter ErrorReporter) Interpreter {
	return Interpreter{
		reporter: reporter,
	}
}

// Interpret interprets the expression and returns a regular Golang value, e.g. nil, string, float64, etc.
func (interpreter Interpreter) Interpret(expr Expr) {
	if e, r := interpreter.visit(expr); e != nil {
		switch err := e.(type) {
		case RuntimeError:
			interpreter.reporter.RuntimeError(err)
		default:
			panic(err)
		}
	} else {
		fmt.Println(stringify(r))
	}
}

func (interpreter Interpreter) visit(expr Expr) (error, interface{}) {
	return expr.Visit(interpreter)
}

func (interpreter Interpreter) VisitBinary(binary Binary) (error, interface{}) {
	// Important: left to right evaluation.
	e, left := interpreter.visit(binary.Left)
	if e != nil {
		return e, nil
	}
	e, right := interpreter.visit(binary.Right)
	if e != nil {
		return e, nil
	}

	switch binary.Operator.Type {
	case TokenMinus:
		if e, leftV := interpreter.assertNumber(binary.Operator, left); e != nil {
			return e, nil
		} else if e, rightV := interpreter.assertNumber(binary.Operator, right); e != nil {
			return e, nil
		} else {
			return nil, leftV - rightV
		}
	case TokenSlash:
		if e, leftV := interpreter.assertNumber(binary.Operator, left); e != nil {
			return e, nil
		} else if e, rightV := interpreter.assertNumber(binary.Operator, right); e != nil {
			return e, nil
		} else {
			return nil, leftV / rightV
		}
	case TokenStar:
		if e, leftV := interpreter.assertNumber(binary.Operator, left); e != nil {
			return e, nil
		} else if e, rightV := interpreter.assertNumber(binary.Operator, right); e != nil {
			return e, nil
		} else {
			return nil, leftV * rightV
		}
	case TokenPlus:
		switch leftV := left.(type) {
		case string:
			if e, rightV := interpreter.assertString(right); e != nil {
				return e, nil
			} else {
				return nil, leftV + rightV
			}
		case float64:
			if e, rightV := interpreter.assertNumber(binary.Operator, right); e != nil {

			} else {
				return nil, leftV + rightV
			}
		default:
			return RuntimeError{
				Token: binary.Operator,
				Msg:   fmt.Sprintf("expected two strings or two numbers but got %v + %v", left, right),
			}, nil
		}
	case TokenGreaterEqual:
		if e, leftV := interpreter.assertNumber(binary.Operator, left); e != nil {
			return e, nil
		} else if e, rightV := interpreter.assertNumber(binary.Operator, right); e != nil {
			return e, nil
		} else {
			return nil, leftV >= rightV
		}
	case TokenGreater:
		if e, leftV := interpreter.assertNumber(binary.Operator, left); e != nil {
			return e, nil
		} else if e, rightV := interpreter.assertNumber(binary.Operator, right); e != nil {
			return e, nil
		} else {
			return nil, leftV > rightV
		}
	case TokenLessEqual:
		if e, leftV := interpreter.assertNumber(binary.Operator, left); e != nil {
			return e, nil
		} else if e, rightV := interpreter.assertNumber(binary.Operator, right); e != nil {
			return e, nil
		} else {
			return nil, leftV <= rightV
		}
	case TokenLess:
		if e, leftV := interpreter.assertNumber(binary.Operator, left); e != nil {
			return e, nil
		} else if e, rightV := interpreter.assertNumber(binary.Operator, right); e != nil {
			return e, nil
		} else {
			return nil, leftV < rightV
		}
	case TokenBangEqual:
		return nil, !interpreter.isEqual(left, right)
	case TokenEqualEqual:
		return nil, interpreter.isEqual(left, right)
	}

	return RuntimeError{
		Token: binary.Operator,
		Msg:   "unknown binary operation",
	}, nil
}

func (interpreter Interpreter) VisitGrouping(grouping Grouping) (error, interface{}) {
	return interpreter.visit(grouping.Expression)
}

func (interpreter Interpreter) VisitLiteral(literal Literal) (error, interface{}) {
	switch v := literal.Value.(type) {
	case Number:
		return nil, v.V
	case String:
		return nil, v.V
	case Boolean:
		return nil, v.V
	default:
		return nil, v
	}
}

func (interpreter Interpreter) VisitUnary(unary Unary) (error, interface{}) {
	e, right := interpreter.visit(unary)
	if e != nil {
		return e, nil
	}

	switch unary.Operator.Type {
	case TokenMinus:
		if e, v := interpreter.assertNumber(unary.Operator, right); e != nil {
			return e, nil
		} else {
			return nil, -v
		}
	case TokenBang:
		return nil, !interpreter.isTruthy(right)
	}

	return RuntimeError{
		Token: unary.Operator,
		Msg:   "unexpected unary operator",
	}, nil
}

func (interpreter Interpreter) VisitTernary(ternary Ternary) (error, interface{}) {
	e, cond := interpreter.visit(ternary.Cond)
	if e != nil {
		return e, nil
	}

	if interpreter.isTruthy(cond) {
		return interpreter.visit(ternary.TrueBranch)
	} else {
		return interpreter.visit(ternary.FalseBranch)
	}
}

// Lox implements truthy as anything that is not nil and not false (strict boolean).
// This mimics Ruby's definition of truthy.
func (interpreter Interpreter) isTruthy(right interface{}) bool {
	if right == nil {
		return false
	}

	switch t := right.(type) {
	case Boolean:
		return t.V
	default:
		return true
	}
}

func (interpreter Interpreter) assertNumber(operator Token, v interface{}) (error, float64) {
	switch t := v.(type) {
	case Number:
		return nil, t.V
	case float64:
		return nil, t
	default:
		return RuntimeError{
			Token: operator,
			Msg:   "operand must be a number.",
		}, 0
	}
}

func (interpreter Interpreter) assertString(v interface{}) (error, string) {
	switch t := v.(type) {
	case String:
		return nil, t.V
	case string:
		return nil, t
	default:
		return RuntimeError{
			Token: Token{},
			Msg:   "operand must be string",
		}, ""
	}
}

func (interpreter Interpreter) isEqual(left interface{}, right interface{}) bool {
	return left == right
}
