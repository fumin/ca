package ca

import (
	"flag"
	"fmt"
	"log"
	"slices"
	"testing"
)

func TestDivide(t *testing.T) {
	tests := []struct {
		f         *Polynomial[*Rat]
		g         []*Polynomial[*Rat]
		quotient  [][]Quotient[*Rat]
		remainder *Polynomial[*Rat]
	}{
		// Example 1, Chapter 2.3, Ideals, Varieties, and Algorithms, D. Cox, J. Little, D. O'Shea.
		{
			f: NewPolynomial(NewRat(0, 1), Deglex, PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(1, 2)}, PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1)}),
			g: []*Polynomial[*Rat]{
				NewPolynomial(NewRat(0, 1), Deglex, PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(1, 1)}, PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1)}),
				NewPolynomial(NewRat(0, 1), Deglex, PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(0, 1)}, PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1)}),
			},
			quotient: [][]Quotient[*Rat]{
				[]Quotient[*Rat]{{Coefficient: NewRat(1, 1), Monomial: newMonomial(0, 1)}},
				[]Quotient[*Rat]{{Coefficient: NewRat(-1, 1)}},
			},
			remainder: NewPolynomial(NewRat(0, 1), Deglex, PolynomialTerm[*Rat]{Coefficient: NewRat(2, 1)}),
		},
		// Example 2, Chapter 2.3, Ideals, Varieties, and Algorithms, D. Cox, J. Little, D. O'Shea.
		{
			f: NewPolynomial(NewRat(0, 1), Deglex,
				PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(2, 1)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(1, 2)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(0, 2)},
			),
			g: []*Polynomial[*Rat]{
				NewPolynomial(NewRat(0, 1), Deglex, PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(1, 1)}, PolynomialTerm[*Rat]{Coefficient: NewRat(-1, 1)}),
				NewPolynomial(NewRat(0, 1), Deglex, PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(0, 2)}, PolynomialTerm[*Rat]{Coefficient: NewRat(-1, 1)}),
			},
			quotient: [][]Quotient[*Rat]{
				[]Quotient[*Rat]{
					{Coefficient: NewRat(1, 1), Monomial: newMonomial(0, 1)},
					{Coefficient: NewRat(1, 1), Monomial: newMonomial(1)},
				},
				[]Quotient[*Rat]{{Coefficient: NewRat(1, 1)}},
			},
			remainder: NewPolynomial(NewRat(0, 1), Deglex,
				PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(0, 1)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(1)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1)},
			),
		},
		// Example 4, Chapter 2.3, Ideals, Varieties, and Algorithms, D. Cox, J. Little, D. O'Shea.
		{
			f: NewPolynomial(NewRat(0, 1), Lex,
				PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(1, 2)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(2, 1)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(2, 0)},
			),
			g: []*Polynomial[*Rat]{
				NewPolynomial(NewRat(0, 1), Lex, PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(2, 0)}, PolynomialTerm[*Rat]{Coefficient: NewRat(-1, 1)}),
				NewPolynomial(NewRat(0, 1), Lex, PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(1, 1)}, PolynomialTerm[*Rat]{Coefficient: NewRat(-1, 1)}),
			},
			quotient: [][]Quotient[*Rat]{
				[]Quotient[*Rat]{
					{Coefficient: NewRat(1, 1), Monomial: newMonomial(0, 1)},
					{Coefficient: NewRat(1, 1)},
				},
				[]Quotient[*Rat]{{Coefficient: NewRat(1, 1), Monomial: newMonomial(0, 1)}},
			},
			remainder: NewPolynomial(NewRat(0, 1), Lex,
				PolynomialTerm[*Rat]{Coefficient: NewRat(2, 1), Monomial: newMonomial(0, 1)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1)},
			),
		},
		// Example 5, Chapter 2.3, Ideals, Varieties, and Algorithms, D. Cox, J. Little, D. O'Shea.
		// F = (f1, f2).
		{
			f: NewPolynomial(NewRat(0, 1), Lex,
				PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(1, 2)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(-1, 1), Monomial: newMonomial(1)},
			),
			g: []*Polynomial[*Rat]{
				NewPolynomial(NewRat(0, 1), Lex, PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(1, 1)}, PolynomialTerm[*Rat]{Coefficient: NewRat(-1, 1)}),
				NewPolynomial(NewRat(0, 1), Lex, PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(0, 2)}, PolynomialTerm[*Rat]{Coefficient: NewRat(-1, 1)}),
			},
			quotient: [][]Quotient[*Rat]{
				[]Quotient[*Rat]{{Coefficient: NewRat(1, 1), Monomial: newMonomial(0, 1)}},
				[]Quotient[*Rat]{},
			},
			remainder: NewPolynomial(NewRat(0, 1), Lex,
				PolynomialTerm[*Rat]{Coefficient: NewRat(-1, 1), Monomial: newMonomial(1)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(0, 1)},
			),
		},
		// Example 5, Chapter 2.3, Ideals, Varieties, and Algorithms, D. Cox, J. Little, D. O'Shea.
		// F = (f2, f1).
		{
			f: NewPolynomial(NewRat(0, 1), Lex,
				PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(1, 2)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(-1, 1), Monomial: newMonomial(1)},
			),
			g: []*Polynomial[*Rat]{
				NewPolynomial(NewRat(0, 1), Lex, PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(0, 2)}, PolynomialTerm[*Rat]{Coefficient: NewRat(-1, 1)}),
				NewPolynomial(NewRat(0, 1), Lex, PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(1, 1)}, PolynomialTerm[*Rat]{Coefficient: NewRat(-1, 1)}),
			},
			quotient: [][]Quotient[*Rat]{
				[]Quotient[*Rat]{{Coefficient: NewRat(1, 1), Monomial: newMonomial(1)}},
				[]Quotient[*Rat]{},
			},
			remainder: NewPolynomial(NewRat(0, 1), Lex),
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var remainder *Polynomial[*Rat]
			quotient := [][]Quotient[*Rat]{}
			testf := NewPolynomial(NewRat(0, 1), Deglex).Set(test.f)
			quotient, remainder = Divide(quotient, testf, test.g)

			if !remainder.Equal(test.remainder) {
				t.Errorf("%v", remainder)
			}

			// Check if quotient * g + remainder == f.
			f := NewPolynomial(test.f.field, test.f.order)
			for i := range quotient {
				for j := range quotient[i] {
					c := NewPolynomial(NewRat(0, 1), Deglex, PolynomialTerm[*Rat]{Coefficient: quotient[i][j].Coefficient})
					q := NewPolynomial(NewRat(0, 1), Deglex, PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: quotient[i][j].Monomial})
					cwgw := mul(c, q, test.g[i])
					f.Add(f, cwgw)
				}
			}
			f.Add(f, remainder)
			if !f.Equal(test.f) {
				t.Errorf("%v", f)
			}
		})
	}
}

func TestLeadingTerm(t *testing.T) {
	tests := []struct {
		x           *Polynomial[*Rat]
		leadingTerm PolynomialTerm[*Rat]
	}{
		{
			x: NewPolynomial(
				NewRat(0, 1), Lex,
				PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(2)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(-1, 1)},
			),
			leadingTerm: PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(2)},
		},
		{
			x: NewPolynomial(
				NewRat(0, 1), Lex,
				PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(1, 2)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(2, 1)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(2)}),
			leadingTerm: PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(1, 2)},
		},
		{
			x: NewPolynomial(
				NewRat(0, 1), Deglex,
				PolynomialTerm[*Rat]{Coefficient: NewRat(1, 2), Monomial: newMonomial(0, 1)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(1, 3), Monomial: newMonomial(0, 0, 1)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(1, -4), Monomial: newMonomial(0, 1, 1)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(1, -5), Monomial: newMonomial(0, 2)},
			),
			leadingTerm: PolynomialTerm[*Rat]{Coefficient: NewRat(-1, 4), Monomial: newMonomial(0, 1, 1)},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			lt := test.x.LeadingTerm()
			if !termEq(lt, test.leadingTerm) {
				t.Errorf("%v", lt)
			}
		})
	}
}

func TestOrder(t *testing.T) {
	tests := []struct {
		words  []Monomial
		order  Order
		sorted []Monomial
	}{
		{
			words: []Monomial{
				newMonomial(0, 0, 1),
				newMonomial(0, 2, 1),
				newMonomial(0, 1, 1),
			},
			order: Lex,
			sorted: []Monomial{
				newMonomial(0, 2, 1),
				newMonomial(0, 1, 1),
				newMonomial(0, 0, 1),
			},
		},
		{
			words: []Monomial{
				newMonomial(0, 1, 0, 1, 1, 0, 0),
				newMonomial(0, 1, 1, 0, 1, 0, 0),
				newMonomial(0, 1, 0, 1, 0, 0, 1),
				newMonomial(0, 0, 0, 1, 1, 0, 1),
				newMonomial(0, 0, 1, 0, 0, 1, 1),
				newMonomial(0, 1, 0, 0, 0, 1, 1),
				newMonomial(0, 0, 1, 1, 1, 0, 0),
				newMonomial(0, 0, 0, 1, 1, 1, 0),
				newMonomial(0, 1, 0, 0, 1, 1, 0),
				newMonomial(0, 1, 0, 0, 1, 0, 1),
				newMonomial(0, 0, 1, 1, 0, 0, 1),
				newMonomial(0, 0, 1, 0, 1, 1, 0),
				newMonomial(0, 1, 0, 1, 0, 1, 0),
				newMonomial(0, 1, 1, 0, 0, 1, 0),
				newMonomial(0, 0, 1, 1, 0, 1, 0),
				newMonomial(0, 1, 1, 1, 0, 0, 0),
				newMonomial(0, 0, 0, 0, 1, 1, 1),
				newMonomial(0, 0, 0, 1, 0, 1, 1),
				newMonomial(0, 0, 1, 0, 1, 0, 1),
				newMonomial(0, 1, 1, 0, 0, 0, 1),
			},
			order: Lex,
			sorted: []Monomial{
				newMonomial(0, 1, 1, 1, 0, 0, 0),
				newMonomial(0, 1, 1, 0, 1, 0, 0),
				newMonomial(0, 1, 1, 0, 0, 1, 0),
				newMonomial(0, 1, 1, 0, 0, 0, 1),
				newMonomial(0, 1, 0, 1, 1, 0, 0),
				newMonomial(0, 1, 0, 1, 0, 1, 0),
				newMonomial(0, 1, 0, 1, 0, 0, 1),
				newMonomial(0, 1, 0, 0, 1, 1, 0),
				newMonomial(0, 1, 0, 0, 1, 0, 1),
				newMonomial(0, 1, 0, 0, 0, 1, 1),
				newMonomial(0, 0, 1, 1, 1, 0, 0),
				newMonomial(0, 0, 1, 1, 0, 1, 0),
				newMonomial(0, 0, 1, 1, 0, 0, 1),
				newMonomial(0, 0, 1, 0, 1, 1, 0),
				newMonomial(0, 0, 1, 0, 1, 0, 1),
				newMonomial(0, 0, 1, 0, 0, 1, 1),
				newMonomial(0, 0, 0, 1, 1, 1, 0),
				newMonomial(0, 0, 0, 1, 1, 0, 1),
				newMonomial(0, 0, 0, 1, 0, 1, 1),
				newMonomial(0, 0, 0, 0, 1, 1, 1),
			},
		},
		{
			words: []Monomial{
				newMonomial(0, 0, 1, 1),
				newMonomial(0, 1, 0, 1),
				newMonomial(0, 2),
				newMonomial(0, 0, 0, 2),
				newMonomial(0, 1, 1),
				newMonomial(0, 0, 2),
			},
			order: Lex,
			sorted: []Monomial{
				newMonomial(0, 2),
				newMonomial(0, 1, 1),
				newMonomial(0, 1, 0, 1),
				newMonomial(0, 0, 2),
				newMonomial(0, 0, 1, 1),
				newMonomial(0, 0, 0, 2),
			},
		},
		{
			words: []Monomial{
				newMonomial(0, 2),
				newMonomial(0, 2, 1),
				newMonomial(0, 0, 2),
				newMonomial(0, 1, 2),
				newMonomial(0, 0, 3),
				newMonomial(),
				newMonomial(0, 1, 1),
				newMonomial(0, 3),
				newMonomial(0, 0, 1),
				newMonomial(0, 1),
			},
			order: Deglex,
			sorted: []Monomial{
				newMonomial(),
				newMonomial(0, 1),
				newMonomial(0, 0, 1),
				newMonomial(0, 2),
				newMonomial(0, 1, 1),
				newMonomial(0, 0, 2),
				newMonomial(0, 3),
				newMonomial(0, 2, 1),
				newMonomial(0, 1, 2),
				newMonomial(0, 0, 3),
			},
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			slices.SortFunc(test.words, test.order)
			if !slices.EqualFunc(test.words, test.sorted, monomialEq) {
				t.Errorf("%v", test.words)
			}
		})
	}
}

func TestPolynomialAdd(t *testing.T) {
	type testcase struct {
		x *Polynomial[*Rat]
		y *Polynomial[*Rat]
		z *Polynomial[*Rat]
	}
	tests := []testcase{
		{
			x: NewPolynomial(
				NewRat(0, 1), Deglex,
				PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(0, 3)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(0, 3)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(3, 1), Monomial: newMonomial(0, 1, 1)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(4, 1), Monomial: newMonomial(0, 2, 1)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(0, 1), Monomial: newMonomial(0, 0, 2)},
			),
			y: NewPolynomial(
				NewRat(0, 1), Deglex,
				PolynomialTerm[*Rat]{Coefficient: NewRat(-3, 1), Monomial: newMonomial(0, 1, 1)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(-7, 1), Monomial: newMonomial(0, 2, 1)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(-1, 1), Monomial: newMonomial(0, 0, 2)},
			),
			z: NewPolynomial(
				NewRat(0, 1), Deglex,
				PolynomialTerm[*Rat]{Coefficient: NewRat(2, 1), Monomial: newMonomial(0, 3)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(-3, 1), Monomial: newMonomial(0, 2, 1)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(-1, 1), Monomial: newMonomial(0, 0, 2)},
			),
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			z := NewPolynomial(test.z.field, test.z.order)
			z.Add(test.x, test.y)
			if !z.Equal(test.z) {
				t.Errorf("%v", z)
			}

			// z = x.
			x := NewPolynomial(NewRat(0, 1), Deglex).Set(test.x)
			z = x
			z.Add(x, test.y)
			if !z.Equal(test.z) {
				t.Errorf("%v", z)
			}

			// z = y.
			y := NewPolynomial(NewRat(0, 1), Deglex).Set(test.y)
			z = y
			z.Add(test.x, y)
			if !z.Equal(test.z) {
				t.Errorf("%v", z)
			}
		})
	}
}

func TestPolynomialAddZEqXY(t *testing.T) {
	tests := []struct {
		x *Polynomial[*Rat]
		z *Polynomial[*Rat]
	}{
		{
			x: NewPolynomial(
				NewRat(0, 1), Deglex,
				PolynomialTerm[*Rat]{Coefficient: NewRat(2, 1), Monomial: newMonomial(0, 3)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(-3, 1), Monomial: newMonomial(0, 1, 1)},
			),
			z: NewPolynomial(
				NewRat(0, 1), Deglex,
				PolynomialTerm[*Rat]{Coefficient: NewRat(4, 1), Monomial: newMonomial(0, 3)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(-6, 1), Monomial: newMonomial(0, 1, 1)},
			),
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			x := NewPolynomial(NewRat(0, 1), Deglex).Set(test.x)
			z := x
			z.Add(x, x)
			if !z.Equal(test.z) {
				t.Errorf("%v", z)
			}
		})
	}
}

func TestPolynomialMul(t *testing.T) {
	tests := []struct {
		x *Polynomial[*Rat]
		y *Polynomial[*Rat]
		z *Polynomial[*Rat]
	}{
		{
			x: NewPolynomial(
				NewRat(0, 1), Deglex,
				PolynomialTerm[*Rat]{Coefficient: NewRat(2, 1), Monomial: newMonomial(0, 1)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(3, 1), Monomial: newMonomial(0, 0, 1)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(-4, 1), Monomial: newMonomial(0, 1, 1)},
			),
			y: NewPolynomial(
				NewRat(0, 1), Deglex,
				PolynomialTerm[*Rat]{Coefficient: NewRat(5, 2), Monomial: newMonomial(0, 1)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(-6, 1), Monomial: newMonomial(0, 0, 1)},
			),
			z: NewPolynomial(
				NewRat(0, 1), Deglex,
				PolynomialTerm[*Rat]{Coefficient: NewRat(5, 1), Monomial: newMonomial(0, 2)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(-9, 2), Monomial: newMonomial(0, 1, 1)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(-18, 1), Monomial: newMonomial(0, 0, 2)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(-10, 1), Monomial: newMonomial(0, 2, 1)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(24, 1), Monomial: newMonomial(0, 1, 2)},
			),
		},
		{
			x: NewPolynomial(
				NewRat(0, 1), Deglex,
				PolynomialTerm[*Rat]{Coefficient: NewRat(3, 1), Monomial: newMonomial(0, 1)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(-4, 1), Monomial: newMonomial(0, 0, 1)},
			),
			y: NewPolynomial(
				NewRat(0, 1), Deglex,
				PolynomialTerm[*Rat]{Coefficient: NewRat(-2, 1)},
			),
			z: NewPolynomial(
				NewRat(0, 1), Deglex,
				PolynomialTerm[*Rat]{Coefficient: NewRat(-6, 1), Monomial: newMonomial(0, 1)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(8, 1), Monomial: newMonomial(0, 0, 1)},
			),
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			z := NewPolynomial(test.z.field, test.z.order).Mul(test.x, test.y)
			if !z.Equal(test.z) {
				t.Errorf("%v", z)
			}
		})
	}
}

func TestPow(t *testing.T) {
	tests := []struct {
		x *Polynomial[*Rat]
		y int
		z *Polynomial[*Rat]
	}{
		{
			x: NewPolynomial(
				NewRat(0, 1), Deglex,
				PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(0, 1)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(-2, 1), Monomial: newMonomial(0, 0, 1)},
			),
			y: 2,
			z: NewPolynomial(
				NewRat(0, 1), Deglex,
				PolynomialTerm[*Rat]{Coefficient: NewRat(4, 1), Monomial: newMonomial(0, 0, 2)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(-4, 1), Monomial: newMonomial(0, 1, 1)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(0, 2)},
			),
		},
		{
			x: NewPolynomial(
				NewRat(0, 1), Deglex,
				PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(0, 1)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(-1, 1), Monomial: newMonomial(0, 0, 1)},
			),
			y: 3,
			z: NewPolynomial(
				NewRat(0, 1), Deglex,
				PolynomialTerm[*Rat]{Coefficient: NewRat(-1, 1), Monomial: newMonomial(0, 0, 3)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(3, 1), Monomial: newMonomial(0, 1, 2)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(-3, 1), Monomial: newMonomial(0, 2, 1)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(1, 1), Monomial: newMonomial(0, 3)},
			),
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			z := NewPolynomial(NewRat(0, 1), Deglex)
			z.Pow(test.x, test.y)
			if !z.Equal(test.z) {
				t.Errorf("%v", z)
			}
		})
	}
}

func TestMulScalar(t *testing.T) {
	tests := []struct {
		x      *Polynomial[*Rat]
		scalar *Rat
		z      *Polynomial[*Rat]
	}{
		{
			x: NewPolynomial(
				NewRat(0, 1), Deglex,
				PolynomialTerm[*Rat]{Coefficient: NewRat(-2, 3), Monomial: newMonomial(0, 1, 2)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(5, 1), Monomial: newMonomial(0, 0, 1)},
			),
			scalar: NewRat(-6, 1),
			z: NewPolynomial(
				NewRat(0, 1), Deglex,
				PolynomialTerm[*Rat]{Coefficient: NewRat(4, 1), Monomial: newMonomial(0, 1, 2)},
				PolynomialTerm[*Rat]{Coefficient: NewRat(-30, 1), Monomial: newMonomial(0, 0, 1)},
			),
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			z := NewPolynomial(test.z.field, test.z.order)
			z.mulScalar(test.scalar, test.x)
			if !z.Equal(test.z) {
				t.Errorf("%v", z)
			}

			z = NewPolynomial(NewRat(0, 1), Deglex).Set(test.x)
			z.mulScalar(test.scalar, z)
			if !z.Equal(test.z) {
				t.Errorf("%v", z)
			}
		})
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	log.SetFlags(log.Lmicroseconds | log.Llongfile | log.LstdFlags)

	m.Run()
}

func mul[K Field[K]](x *Polynomial[K], y ...*Polynomial[K]) *Polynomial[K] {
	z := x
	for i := range y {
		z = NewPolynomial(z.field, z.order).Mul(z, y[i])
	}
	return z
}

func termEq[K Field[K]](a, b PolynomialTerm[K]) bool {
	if eq := monomialEq(a.Monomial, b.Monomial); !eq {
		return false
	}
	if !a.Coefficient.Equal(b.Coefficient) {
		return false
	}
	return true
}
