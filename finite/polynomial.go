package finite

import (
	"cmp"
	"math/big"
	"math/rand"
	"slices"

	"github.com/fumin/ca"
)

// A finite is an element of a finite field GF(q) with order q = p^n where p is a prime number.
type finite[K ca.Field[K]] interface {
	ca.Field[K]
	// characteristic returns p in the finite field GF(p^n).
	characteristic() int
	// primePower returns n in the finite field GF(p^n).
	primePower() int
}

type factor[K finite[K]] struct {
	A *ca.Polynomial[K]
	N int
}

func factorize[K finite[K]](f *ca.Polynomial[K]) []factor[K] {
	factors := make([]factor[K], 0)
	sf := squareFree(f)
	for _, sfi := range sf {
		dd := distinctDegree(sfi.A)
		for _, ddi := range dd {
			ed := equalDegree(ddi.A, ddi.N)
			for _, edi := range ed {
				factors = append(factors, factor[K]{A: edi, N: sfi.N})
			}
		}
	}
	return factors
}

func equalDegree[K finite[K]](f *ca.Polynomial[K], d int) []*ca.Polynomial[K] {
	n := int(f.LeadingTerm().Monomial[0].Int64())
	r := n / d
	k := f.Field()
	// q is the order of the field k.
	q := expi(k.characteristic(), k.primePower())
	// numFq is the number of GF(q) polynomials with degree less than n.
	numFq := expi(q, n)

	factors := []*ca.Polynomial[K]{f}
	for len(factors) < r {
		hi := rand.Intn(numFq-q) + q
		h := ithPoly(poly0(f), hi, q)
		g, _, _, _, _ := gcd(h, f)
		if isOne(g) {
			qd := (expi(q, d) - 1) / 2
			hqd := sub(pow(h, qd), poly1(h))
			_, g = divide(hqd, f)
		}

		checked := make(map[string]struct{})
		for {
			// Find u in factors where deg(u) > d.
			var ui int = -1
			for i, f := range factors {
				if _, ok := checked[f.String()]; ok {
					continue
				}
				if int(f.LeadingTerm().Monomial[0].Int64()) > d {
					ui = i
					break
				}
			}
			if ui == -1 {
				break
			}
			u := factors[ui]

			gcdu, _, _, _, _ := gcd(g, u)
			if !isOne(gcdu) && !gcdu.Equal(u) {
				factors = slices.Delete(factors, ui, ui+1)
				factors = append(factors, gcdu)
				ugcd, _ := divide(u, gcdu)
				factors = append(factors, ugcd)
			}
			checked[u.String()] = struct{}{}
		}
	}

	return factors
}

func distinctDegree[K finite[K]](f *ca.Polynomial[K]) []factor[K] {
	degF := int(f.LeadingTerm().Monomial[0].Int64())
	k := f.Field()
	// q is the order of the field k.
	q := expi(k.characteristic(), k.primePower())
	neg1 := k.NewZero()
	neg1.Sub(k.NewZero(), k.NewOne())

	fs := poly0(f).Set(f)
	s := make([]factor[K], 0)
	for i := 1; i <= degF/2+1; i++ {
		qi := expi(q, i)
		xqi := ca.NewPolynomial(k, fs.Order(),
			ca.PolynomialTerm[K]{Coefficient: k.NewOne(), Monomial: ca.Monomial{big.NewInt(int64(qi))}},
			ca.PolynomialTerm[K]{Coefficient: neg1, Monomial: ca.Monomial{big.NewInt(1)}})
		_, xqi = divide(xqi, fs)
		g, _, _, _, _ := gcd(fs, xqi)
		if !isOne(g) {
			s = append(s, factor[K]{A: g, N: i})
			fs, _ = divide(fs, g)
		}
	}
	if !isOne(fs) {
		s = append(s, factor[K]{A: fs, N: degF})
	}
	if len(s) == 0 {
		s = append(s, factor[K]{A: f, N: 1})
	}
	return s
}

func squareFree[K finite[K]](f *ca.Polynomial[K]) []factor[K] {
	c, _, _, _, _ := gcd(f, differentiate(f))
	w, _ := divide(f, c)

	r := make([]factor[K], 0)
	var i int = 1
	for !isOne(w) {
		y, _, _, _, _ := gcd(w, c)
		fac, _ := divide(w, y)
		if !isOne(fac) {
			r = append(r, factor[K]{A: fac, N: i})
		}
		w = y
		c, _ = divide(c, y)
		i++
	}

	if !isOne(c) {
		// Let the field be GF(q), where q = p^e.
		// Compute cRoot = c^{1/p}.
		k := f.Field()
		p, e := k.characteristic(), k.primePower()
		cRoot := poly0(c)
		for cc, cw := range c.Terms() {
			// x^{ap} -> x^a
			w := make(ca.Monomial, len(cw))
			for i := range w {
				w[i] = big.NewInt(0).Div(cw[i], big.NewInt(int64(p)))
			}

			// c -> c^{p^(e-1)}
			pe1 := expi(p, e-1)
			rc := exp(cc, pe1)

			cRoot.Add(cRoot, ca.NewPolynomial(k, c.Order(), ca.PolynomialTerm[K]{Coefficient: rc, Monomial: w}))
		}

		// Do factorization on c^{1/p}.
		sf := squareFree(cRoot)
		for i := range sf {
			sf[i].N *= p
		}
		r = append(r, sf...)
	}
	return r
}

func differentiate[K ca.Field[K]](a *ca.Polynomial[K]) *ca.Polynomial[K] {
	k := a.Field()
	aP := poly0(a)
	for ac, aw := range a.Terms() {
		if len(aw) == 0 {
			continue
		}
		w := make(ca.Monomial, len(aw))
		for i := range w {
			w[i] = big.NewInt(0).Sub(aw[i], big.NewInt(1))
		}

		c := k.NewZero()
		for exp := big.NewInt(0); exp.Cmp(aw[0]) < 0; exp.Add(exp, big.NewInt(1)) {
			c.Add(c, ac)
		}

		aP.Add(aP, ca.NewPolynomial(k, a.Order(), ca.PolynomialTerm[K]{Coefficient: c, Monomial: w}))
	}
	return aP
}

func inverse[K ca.Field[K]](a, p *ca.Polynomial[K]) *ca.Polynomial[K] {
	r, _, v, _, _ := gcd(p, a)
	if !isOne(r) {
		return nil
	}
	return v
}

func gcd[K ca.Field[K]](a, b *ca.Polynomial[K]) (g, u, v, a1, b1 *ca.Polynomial[K]) {
	r0, r1 := a, b
	s0, s1 := poly1(a), poly0(a)
	t0, t1 := poly0(a), poly1(a)

	var n1i int = 1
	for r1.Len() != 0 {
		n1i *= -1
		q, _ := divide(r0, r1)
		r2 := sub(r0, mul(q, r1))
		s2 := sub(s0, mul(q, s1))
		t2 := sub(t0, mul(q, t1))

		r0, r1 = r1, r2
		s0, s1 = s1, s2
		t0, t1 = t1, t2
	}

	k := a.Field()
	var n1, nn1 K
	if n1i == 1 {
		n1 = k.NewOne()
		nn1 = k.NewZero().Sub(k.NewZero(), k.NewOne())
	} else {
		n1 = k.NewZero().Sub(k.NewZero(), k.NewOne())
		nn1 = k.NewOne()
	}
	a1 = mulScalar(t1, n1)
	b1 = mulScalar(s1, nn1)

	// Make g monic.
	c := r0.LeadingTerm().Coefficient
	g = mulScalar(r0, k.NewZero().Inv(c))
	u = mulScalar(s0, k.NewZero().Inv(c))
	v = mulScalar(t0, k.NewZero().Inv(c))
	a1 = mulScalar(a1, c)
	b1 = mulScalar(b1, c)
	return g, u, v, a1, b1
}

func divide[K ca.Field[K]](a, b *ca.Polynomial[K]) (*ca.Polynomial[K], *ca.Polynomial[K]) {
	k := a.Field()
	a2 := poly0(a).Set(a)
	quotient := make([][]ca.Quotient[K], 0)
	quotient, r := ca.Divide(quotient, a2, []*ca.Polynomial[K]{b})

	q := poly0(a)
	for i := range quotient {
		for j := range quotient[i] {
			c := ca.NewPolynomial(k, q.Order(), ca.PolynomialTerm[K]{Coefficient: quotient[i][j].Coefficient})
			t := ca.NewPolynomial(k, q.Order(), ca.PolynomialTerm[K]{Coefficient: k.NewOne(), Monomial: quotient[i][j].Monomial})
			cwgw := mul[K](c, t)
			q.Add(q, cwgw)
		}
	}
	return q, r
}

func ithPoly[K ca.Field[K]](poly *ca.Polynomial[K], ith, base int) *ca.Polynomial[K] {
	k := poly.Field()
	var pow int = -1
	for ith != 0 {
		pow++
		var r int
		ith, r = ith/base, ith%base

		c := k.NewZero()
		for range r {
			c.Add(c, k.NewOne())
		}
		w := ca.Monomial{big.NewInt(int64(pow))}
		poly.Add(poly, ca.NewPolynomial(k, poly.Order(), ca.PolynomialTerm[K]{Coefficient: c, Monomial: w}))
	}
	return poly
}

func isOne[K ca.Field[K]](a *ca.Polynomial[K]) bool {
	if a.Len() != 1 {
		return false
	}
	lt := a.LeadingTerm()
	if len(lt.Monomial) != 0 {
		return false
	}
	if !lt.Coefficient.Equal(a.Field().NewOne()) {
		return false
	}
	return true
}

func pow[K ca.Field[K]](a *ca.Polynomial[K], n int) *ca.Polynomial[K] {
	y := poly1(a)
	for range n {
		y = mul(y, a)
	}
	return y
}

func sub[K ca.Field[K]](x, y *ca.Polynomial[K]) *ca.Polynomial[K] {
	k := x.Field()
	neg1 := k.NewZero()
	neg1.Sub(neg1, k.NewOne())
	negY := mulScalar(y, neg1)
	return poly0(x).Add(x, negY)
}

func mulScalar[K ca.Field[K]](a *ca.Polynomial[K], b K) *ca.Polynomial[K] {
	bP := ca.NewPolynomial(a.Field(), a.Order(), ca.PolynomialTerm[K]{Coefficient: b})
	ab := poly0(a).Mul(a, bP)
	ab.SymbolStringer = a.SymbolStringer
	return ab
}

func mul[K ca.Field[K]](x *ca.Polynomial[K], y ...*ca.Polynomial[K]) *ca.Polynomial[K] {
	z := x
	for i := range y {
		z = poly0(z).Mul(z, y[i])
		z.SymbolStringer = x.SymbolStringer
	}
	return z
}

func expi(x, n int) int {
	return int(exp(ca.NewRat(int64(x), 1), n).Num().Int64())
}

func exp[K ca.Field[K]](x K, n int) K {
	switch {
	case n < 0:
		return exp(x.NewOne().Inv(x), -n)
	case n == 0:
		return x.NewOne()
	case n%2 == 0:
		return exp(x.NewOne().Mul(x, x), n/2)
	default:
		return x.NewOne().Mul(x, exp(x.NewOne().Mul(x, x), (n-1)/2))
	}
}

func poly1[K ca.Field[K]](x *ca.Polynomial[K]) *ca.Polynomial[K] {
	y := ca.NewPolynomial(x.Field(), x.Order(), ca.PolynomialTerm[K]{Coefficient: x.Field().NewOne()})
	y.SymbolStringer = x.SymbolStringer
	return y
}

func poly0[K ca.Field[K]](x *ca.Polynomial[K]) *ca.Polynomial[K] {
	y := ca.NewPolynomial(x.Field(), x.Order())
	y.SymbolStringer = x.SymbolStringer
	return y
}

func polynomialCmp[K ca.Field[K]](x, y *ca.Polynomial[K]) int {
	xTerms := make([]ca.PolynomialTerm[K], 0)
	for c, w := range x.Terms() {
		xTerms = append(xTerms, ca.PolynomialTerm[K]{Coefficient: c, Monomial: w})
	}
	yTerms := make([]ca.PolynomialTerm[K], 0)
	for c, w := range y.Terms() {
		yTerms = append(yTerms, ca.PolynomialTerm[K]{Coefficient: c, Monomial: w})
	}

	// Compare monomials.
	for i := range xTerms {
		if i >= len(yTerms) {
			return 1
		}
		xw := xTerms[i].Monomial
		yw := yTerms[i].Monomial
		if wo := x.Order()(xw, yw); wo != 0 {
			return wo
		}
	}
	if len(xTerms) < len(yTerms) {
		return -1
	}

	// Compare coefficients.
	for i := range xTerms {
		xc := xTerms[i].Coefficient
		yc := yTerms[i].Coefficient
		if co := cmp.Compare(xc.String(), yc.String()); co != 0 {
			return co
		}
	}
	return 0
}
