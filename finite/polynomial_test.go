package finite

import (
	"fmt"
	"slices"
	"testing"

	"github.com/fumin/ca"
	"github.com/pkg/errors"
)

func TestFactorize(t *testing.T) {
	tests := []struct {
		order  int
		a      string
		factor []factorStr
	}{
		{
			order: 2,
			a:     "x^2+x",
			factor: []factorStr{
				{a: "x", n: 1},
				{a: "x+1", n: 1},
			},
		},
		{
			order: 2,
			a:     "x^2+1",
			factor: []factorStr{
				{a: "x+1", n: 2},
			},
		},
		{
			order: 2,
			a:     "x^2+x+1",
			factor: []factorStr{
				{a: "x^2+x+1", n: 1},
			},
		},
		{
			order: 7,
			a:     "(x+2)^6(x+5)^6(x^5+x+4)^6(x^2+x+3)^4(x^2+2x+5)^4",
			factor: []factorStr{
				{a: "x+2", n: 6},
				{a: "x+5", n: 6},
				{a: "x^2+x+3", n: 4},
				{a: "x^2+2x+5", n: 4},
				{a: "x^5+x+4", n: 6},
			},
		},
	}

	for testI, test := range tests {
		t.Run(fmt.Sprintf("%d", testI), func(t *testing.T) {
			t.Parallel()
			a := parseMust(test.order, test.a)
			tFactors := make([]factor[*prime], len(test.factor))
			for i, f := range test.factor {
				tFactors[i] = factor[*prime]{A: parseMust(test.order, f.a), N: f.n}
			}

			factors := factorize(a)
			slices.SortFunc(factors, func(a, b factor[*prime]) int { return polynomialCmp(a.A, b.A) })
			if len(factors) != len(tFactors) {
				t.Fatalf("%v", factors)
			}
			for i, f := range factors {
				tf := tFactors[i]
				if !(f.A.Equal(tf.A) && f.N == tf.N) {
					t.Errorf("%d got %v want %v", i, f, tf)
				}
			}
		})
	}
}

func TestEqualDegree(t *testing.T) {
	tests := []struct {
		order  int
		a      factorStr
		factor []string
	}{
		{
			order:  2,
			a:      factorStr{a: "x(x+1)", n: 1},
			factor: []string{"x", "x+1"},
		},
		{
			order:  3,
			a:      factorStr{a: "x(x+2)", n: 1},
			factor: []string{"x", "x+2"},
		},
		{
			order:  2,
			a:      factorStr{a: "(x^5+x^2+1)(x^5+x^3+1)(x^5+x^4+x^2+x+1)", n: 5},
			factor: []string{"x^5+x^2+1", "x^5+x^3+1", "x^5+x^4+x^2+x+1"},
		},
	}

	for testI, test := range tests {
		t.Run(fmt.Sprintf("%d", testI), func(t *testing.T) {
			t.Parallel()
			a := factor[*prime]{A: parseMust(test.order, test.a.a), N: test.a.n}
			tfactors := make([]*ca.Polynomial[*prime], len(test.factor))
			for i, f := range test.factor {
				tfactors[i] = parseMust(test.order, f)
			}

			factors := equalDegree(a.A, a.N)
			slices.SortFunc(factors, polynomialCmp)
			if len(factors) != len(tfactors) {
				t.Fatalf("%v", factors)
			}
			for i, f := range factors {
				if !f.Equal(tfactors[i]) {
					t.Errorf("got %v want %v", f, tfactors[i])
				}
			}
		})
	}
}

func TestDistinctDegree(t *testing.T) {
	tests := []struct {
		order  int
		a      string
		factor []factorStr
	}{
		{
			order: 3,
			a:     "x(x+2)(x^2+x+2)",
			factor: []factorStr{
				{a: "x(x+2)", n: 1},
				{a: "x^2+x+2", n: 2},
			},
		},
	}

	for testI, test := range tests {
		t.Run(fmt.Sprintf("%d", testI), func(t *testing.T) {
			t.Parallel()
			a := parseMust(test.order, test.a)
			tfactors := make([]factor[*prime], len(test.factor))
			for i, f := range test.factor {
				tfactors[i] = factor[*prime]{A: parseMust(test.order, f.a), N: f.n}
			}

			factors := distinctDegree(a)
			if len(factors) != len(tfactors) {
				t.Fatalf("%v", factors)
			}
			for i, f := range factors {
				tf := tfactors[i]
				if !(f.A.Equal(tf.A) && f.N == tf.N) {
					t.Errorf("got %v want %v", f, tf)
				}
			}
		})
	}
}

func TestSquareFree(t *testing.T) {
	tests := []struct {
		order  int
		a      string
		factor []factorStr
	}{
		{
			order: 3,
			a:     "x^11+2x^9+2x^8+x^6+x^5+2x^3+2x^2+1",
			factor: []factorStr{
				{a: "x+1", n: 1},
				{a: "x+2", n: 4},
				{a: "x^2+1", n: 3},
			},
		},
		{
			order:  5,
			a:      "x^6+x^4+x^3-x^2-2x-1",
			factor: []factorStr{{a: "x^3+3x+3", n: 2}},
		},
		{
			order: 13,
			a:     "x^7+3x^6+5x^5+7x^4+7x^3+5x^2+3x+1",
			factor: []factorStr{
				{a: "x^2+1", n: 2},
				{a: "x+1", n: 3},
			},
		},
		{
			order: 7,
			a:     "(x+2)^6(x+5)^6",
			factor: []factorStr{
				{a: "x^2+3", n: 6},
			},
		},
	}

	for testI, test := range tests {
		t.Run(fmt.Sprintf("%d", testI), func(t *testing.T) {
			t.Parallel()
			a := parseMust(test.order, test.a)
			testFactors := make([]factor[*prime], len(test.factor))
			for i, f := range test.factor {
				testFactors[i] = factor[*prime]{A: parseMust(test.order, f.a), N: f.n}
			}

			factors := squareFree(a)
			if len(factors) != len(testFactors) {
				t.Fatalf("%v", factors)
			}
			for i, f := range factors {
				tf := testFactors[i]
				if !(f.A.Equal(tf.A) && f.N == tf.N) {
					t.Errorf("%d got %v want %v", i, f, tf)
				}
			}
			// Check that the reconstruction from factors equals a.
			af := ca.NewPolynomial(a.Field(), a.Order(), ca.PolynomialTerm[*prime]{Coefficient: a.Field().NewOne()})
			buf := ca.NewPolynomial(a.Field(), a.Order())
			for _, f := range factors {
				for range f.N {
					af.Set(buf.Mul(af, f.A))
				}
			}
			if !af.Equal(a) {
				t.Errorf("got %v want %v", af, a)
			}
		})
	}
}

func TestDifferentiate(t *testing.T) {
	tests := []struct {
		a  *ca.Polynomial[*prime]
		aP *ca.Polynomial[*prime]
	}{
		{
			a:  parseMust(13, "3x^7+8x^5-2x^2+5"),
			aP: parseMust(13, "-5x^6+x^4-4x"),
		},
	}

	for testI, test := range tests {
		t.Run(fmt.Sprintf("%d", testI), func(t *testing.T) {
			t.Parallel()
			aP := differentiate(test.a)
			if !aP.Equal(test.aP) {
				t.Errorf("got %v want %v", aP, test.aP)
			}
		})
	}
}

func TestInverse(t *testing.T) {
	tests := []struct {
		order int
		a     string
		p     string
		inv   string
	}{
		{
			order: 2,
			a:     "x^6+x^4+x+1",
			p:     "x^8+x^4+x^3+x+1",
			inv:   "x^7+x^6+x^3+x",
		},
	}

	for testI, test := range tests {
		t.Run(fmt.Sprintf("%d", testI), func(t *testing.T) {
			t.Parallel()
			a := parseMust(test.order, test.a)
			p := parseMust(test.order, test.p)
			tinv := parseMust(test.order, test.inv)
			inv := inverse(a, p)
			if !inv.Equal(tinv) {
				t.Errorf("got %v want %v", inv, tinv)
			}
		})
	}
}

func TestGcd(t *testing.T) {
	tests := []struct {
		order int
		a     string
		b     string
		g     string
		u     string
		v     string
		a1    string
		b1    string
	}{
		{
			order: 101,
			a:     "x^2+7x+6",
			b:     "x^2-5x-6",
			g:     "x+1",
			u:     "59",
			v:     "42",
			a1:    "x+6",
			b1:    "x-6",
		},
	}

	for testI, test := range tests {
		t.Run(fmt.Sprintf("%d", testI), func(t *testing.T) {
			t.Parallel()
			a := parseMust(test.order, test.a)
			b := parseMust(test.order, test.b)
			tg := parseMust(test.order, test.g)
			tu := parseMust(test.order, test.u)
			tv := parseMust(test.order, test.v)
			ta1 := parseMust(test.order, test.a1)
			tb1 := parseMust(test.order, test.b1)
			g, u, v, a1, b1 := gcd(a, b)
			if !g.Equal(tg) {
				t.Errorf("got %v want %v", g, tg)
			}
			if !u.Equal(tu) {
				t.Errorf("got %v want %v", u, tu)
			}
			if !v.Equal(tv) {
				t.Errorf("got %v want %v", v, tv)
			}
			if !a1.Equal(ta1) {
				t.Errorf("got %v want %v", a1, ta1)
			}
			if !b1.Equal(tb1) {
				t.Errorf("got %v want %v", b1, tb1)
			}
			// Check Bezout's identity.
			aubv := poly0(g).Add(mul(a, u), mul(b, v))
			if !aubv.Equal(tg) {
				t.Errorf("got %v want %v", aubv, tg)
			}
			if ga1 := mul(g, a1); !ga1.Equal(a) {
				t.Errorf("got %v want %v", ga1, a)
			}
			if gb1 := mul(g, b1); !gb1.Equal(b) {
				t.Errorf("got %v want %v", gb1, b)
			}
		})
	}
}

func TestIthPoly(t *testing.T) {
	tests := []struct {
		ith  int
		base int
		poly string
	}{
		{
			ith:  3,
			base: 2,
			poly: "x+1",
		},
	}

	for testI, test := range tests {
		t.Run(fmt.Sprintf("%d", testI), func(t *testing.T) {
			t.Parallel()
			tpoly := parseMust(test.base, test.poly)
			poly := ithPoly(poly0(tpoly), test.ith, test.base)
			if !poly.Equal(tpoly) {
				t.Errorf("%v %v", poly, tpoly)
			}
		})
	}
}

type factorStr struct {
	a string
	n int
}

func parse(order int, s string) (*ca.Polynomial[*prime], error) {
	variables := map[string]ca.Symbol{"x": 0}
	rp, err := ca.Parse(variables, ca.Deglex, s)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	// Cast coefficients from rationals to GF(order).
	field := newPrime(order, 0)
	p := ca.NewPolynomial[*prime](field, rp.Order())
	p.SymbolStringer = rp.SymbolStringer
	for rc, w := range rp.Terms() {
		c := newPrime(order, int(rc.Num().Int64()))
		p.Add(p, ca.NewPolynomial(p.Field(), p.Order(), ca.PolynomialTerm[*prime]{Coefficient: c, Monomial: w}))
	}
	return p, nil
}

func parseMust(order int, s string) *ca.Polynomial[*prime] {
	p, err := parse(order, s)
	if err != nil {
		panic(err)
	}
	return p
}
