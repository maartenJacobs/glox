package internal

import (
	"strings"
)

type Printer struct {
	output strings.Builder
}

func (p *Printer) Print(expr Expr) string {
	p.output = strings.Builder{}
	expr.Visit(ExprVisitor{
		VisitBinary:   p.visitBinary,
		VisitGrouping: p.visitGroupingExpr,
		VisitLiteral:  p.visitLiteral,
		VisitUnary:    p.visitUnary,
		VisitTernary:  p.visitTernary,
	})
	return p.output.String()
}

func (p *Printer) visitBinary(binary Binary) interface{} {
	p.parenthesize(binary.Operator.Lexeme, binary.Left, binary.Right)
	return nil
}

func (p *Printer) visitGroupingExpr(grouping Grouping) interface{} {
	p.parenthesize("group", grouping.Expression)
	return nil
}

func (p *Printer) visitUnary(unary Unary) interface{} {
	p.parenthesize(unary.Operator.Lexeme, unary.Right)
	return nil
}

func (p *Printer) visitTernary(ternary Ternary) interface{} {
	p.parenthesize("?:", ternary.Cond, ternary.TrueBranch, ternary.FalseBranch)
	return nil
}

func (p *Printer) visitLiteral(literal Literal) interface{} {
	if literal.Value == nil {
		p.output.WriteString("nil")
	} else {
		p.output.WriteString(literal.Value.String())
	}
	return nil
}

func (p *Printer) parenthesize(name string, exprs ...Expr) {
	p.output.WriteString("(")
	p.output.WriteString(name)
	var eprinter Printer
	for _, expr := range exprs {
		p.output.WriteString(" ")
		p.output.WriteString(eprinter.Print(expr))
	}
	p.output.WriteString(")")
}
