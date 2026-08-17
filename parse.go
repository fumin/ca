package ca

import (
	"bytes"
	"fmt"
	"math/big"
	"strconv"

	"github.com/pkg/errors"

	"github.com/fumin/nag/parse"
	"github.com/fumin/nag/parse/scan"
)

// Parse parses input and returns the polynomial it represents.
func Parse(variables map[string]Symbol, order Order, input string) (*Polynomial[*Rat], error) {
	n, err := parse.Parse(scan.NewScanner(bytes.NewBufferString(input)))
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	p, err := evaluate(n, variables, order)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	symMap := make(map[Symbol]string, len(variables))
	for v, sym := range variables {
		symMap[sym] = v
	}
	p.SymbolStringer = func(s Symbol) string { return symMap[s] }

	return p, nil
}

func evaluate(n *parse.Node, variables map[string]Symbol, order Order) (*Polynomial[*Rat], error) {
	switch n.Token.Type {
	case scan.Parenthesis:
		return evaluateParenthesis(n, variables, order)
	case scan.Operator:
		return evaluateOperator(n, variables, order)
	case scan.Int:
		return evaluateInt(n, order)
	case scan.Identifier:
		return evaluateIdentifier(n, variables, order)
	default:
		return nil, errors.Errorf("unknown node %#v", n)
	}
}

func evaluateParenthesis(n *parse.Node, variables map[string]Symbol, order Order) (*Polynomial[*Rat], error) {
	if n.Left == nil {
		return nil, errors.Errorf("%#v", n)
	}
	return evaluate(n.Left, variables, order)
}

func evaluateOperator(n *parse.Node, variables map[string]Symbol, order Order) (*Polynomial[*Rat], error) {
	switch n.Token.Text {
	case "+":
		return evaluatePlus(n, variables, order)
	case "-":
		return evaluateMinus(n, variables, order)
	case "*":
		return evaluateMultiply(n, variables, order)
	case "/":
		return evaluateDivide(n, order)
	case "^":
		return evaluatePower(n, variables, order)
	default:
		return nil, errors.Errorf("%#v", n)
	}
}

func evaluateIdentifier(n *parse.Node, variables map[string]Symbol, order Order) (*Polynomial[*Rat], error) {
	s, ok := variables[n.Token.Text]
	if !ok {
		return nil, errors.Errorf("%#v", n)
	}
	m := makeMonomial(int(s) + 1)
	m[s] = big.NewInt(1)
	p := NewPolynomial(NewRat(0, 1), order, PolynomialTerm[*Rat]{Coefficient: NewRat(0, 1).NewOne(), Monomial: m})
	return p, nil
}

func evaluatePlus(n *parse.Node, variables map[string]Symbol, order Order) (*Polynomial[*Rat], error) {
	left, right, err := evaluateLeftRight(n, variables, order)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("%#v", n))
	}
	z := NewPolynomial(NewRat(0, 1), order).Add(left, right)
	return z, nil
}

func evaluateMinus(n *parse.Node, variables map[string]Symbol, order Order) (*Polynomial[*Rat], error) {
	left, right, err := evaluateLeftRight(n, variables, order)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("%#v", n))
	}
	negRight := right.mulScalar(NewRat(-1, 1), right)
	z := NewPolynomial(NewRat(0, 1), order).Add(left, negRight)
	return z, nil
}

func evaluateMultiply(n *parse.Node, variables map[string]Symbol, order Order) (*Polynomial[*Rat], error) {
	left, right, err := evaluateLeftRight(n, variables, order)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("%#v", n))
	}
	z := NewPolynomial(NewRat(0, 1), order).Mul(left, right)
	return z, nil
}

func evaluateDivide(n *parse.Node, order Order) (*Polynomial[*Rat], error) {
	if n.Left == nil {
		return nil, errors.Errorf("%#v", n)
	}
	num, err := strconv.ParseInt(n.Left.Token.Text, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("%#v", n))
	}
	if n.Right == nil {
		return nil, errors.Errorf("%#v", n)
	}
	denom, err := strconv.ParseInt(n.Right.Token.Text, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("%#v", n))
	}
	p := NewPolynomial(NewRat(0, 1), order, PolynomialTerm[*Rat]{Coefficient: NewRat(num, denom)})
	return p, nil
}

func evaluatePower(n *parse.Node, variables map[string]Symbol, order Order) (*Polynomial[*Rat], error) {
	if n.Left == nil {
		return nil, errors.Errorf("%#v", n)
	}
	left, err := evaluate(n.Left, variables, order)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("%#v", n))
	}
	right, err := strconv.Atoi(n.Right.Token.Text)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("%#v", n))
	}
	z := NewPolynomial(NewRat(0, 1), order).Pow(left, right)
	return z, nil
}

func evaluateInt(n *parse.Node, order Order) (*Polynomial[*Rat], error) {
	i, err := strconv.ParseInt(n.Token.Text, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("%#v", n))
	}
	p := NewPolynomial(NewRat(0, 1), order, PolynomialTerm[*Rat]{Coefficient: NewRat(i, 1)})
	return p, nil
}

func evaluateLeftRight(n *parse.Node, variables map[string]Symbol, order Order) (*Polynomial[*Rat], *Polynomial[*Rat], error) {
	if n.Left == nil {
		return nil, nil, errors.Errorf("%#v", n)
	}
	left, err := evaluate(n.Left, variables, order)
	if err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("%#v", n))
	}
	if n.Right == nil {
		return nil, nil, errors.Errorf("%#v", n)
	}
	right, err := evaluate(n.Right, variables, order)
	if err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("%#v", n))
	}
	return left, right, nil
}
