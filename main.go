package main

import (
	"errors"
	"fmt"
	"github.com/MoonShining/math/lexer"
	"github.com/MoonShining/math/parser"
)

func main() {
	l := lexer.NewLexer("(-1+2)*3")
	p := parser.New(l)
	exp := p.ParseExpression(0)
	fmt.Printf("%+v\n", exp)

	res, _ := compute(exp)
	fmt.Println(res)
}

func compute(exp parser.Expression) (int64, error) {
	switch e := exp.(type) {
	case parser.NumberExpression:
		return e.Value, nil
	case parser.PrefixExpression:
		res, err := compute(e.Expr)
		if err != nil {
			return 0, err
		}
		return -res, nil
	case parser.InfixExpression:
		a, err := compute(e.Left)
		if err != nil {
			return 0, err
		}
		b, err := compute(e.Right)
		if err != nil {
			return 0, err
		}
		switch e.Operator {
		case lexer.ADD:
			return a + b, nil
		case lexer.SUB:
			return a - b, nil
		case lexer.MUL:
			return a * b, nil
		case lexer.DIV:
			return a / b, nil
		default:
			return 0, errors.New("won't happen")
		}
	default:
		return 0, errors.New("won't happen")
	}
}
