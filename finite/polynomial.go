package finite

import (
	"cmp"
	"crypto/rand"
	"fmt"
	"math/big"
	"slices"

	"github.com/fumin/ca"
)

// A finite is an element of a finite field GF(q) with order q = p^n where p is a prime number.
type finite[K ca.Field[K]] interface {
	ca.Field[K]
	// characteristic returns p in the finite field GF(p^n).
	characteristic() *big.Int
	// primePower returns n in the finite field GF(p^n).
	primePower() *big.Int
	// ith returns the i'th element in the field.
	ith(i *big.Int) K
}

type factor[K finite[K]] struct {
	A *ca.Polynomial[K]
	N *big.Int
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

// Factoring Polynomials Over Finite Fields: A Survey, Joachim von zur Gathen, Daniel Panario.
func equalDegree[K finite[K]](f *ca.Polynomial[K], d *big.Int) []*ca.Polynomial[K] {
	n := f.LeadingTerm().Monomial[0]
	r := big.NewInt(0).Div(n, d)
	k := f.Field()
	// q is the order of the field k.
	q := pow(ca.NewRat(0, 1).SetFrac(k.characteristic(), big.NewInt(1)), k.primePower()).Num()

	factors := []*ca.Polynomial[K]{f}
	for uint64(len(factors)) < r.Uint64() {
		h := randPoly(poly0(f), n)
		g, _, _, _, _ := gcd(h, f)
		if isOne(g) {
			qd := pow(ca.NewRat(0, 1).SetFrac(q, big.NewInt(1)), d).Num()
			qdExp := new(big.Int).Div(new(big.Int).Sub(qd, big.NewInt(1)), big.NewInt(2))
			hqd := sub(pow(newPolynomialMod(h, f), qdExp).p, poly1(h))
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
				if deg(f).Cmp(d) > 0 {
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
	k := f.Field()
	// q is the order of the field k.
	q := pow(ca.NewRat(0, 1).SetFrac(k.characteristic(), big.NewInt(1)), k.primePower()).Num()
	neg1 := k.NewZero()
	neg1.Sub(k.NewZero(), k.NewOne())

	fs := poly0(f).Set(f)
	s := make([]factor[K], 0)
	for i := big.NewInt(1); deg(fs).Cmp(new(big.Int).Mul(big.NewInt(2), i)) >= 0; i = new(big.Int).Add(i, big.NewInt(1)) {
		qi := pow(ca.NewRat(0, 1).SetFrac(q, big.NewInt(1)), i).Num()
		x := newPolynomialMod(ca.NewPolynomial(k, fs.Order(), ca.PolynomialTerm[K]{Coefficient: k.NewOne(), Monomial: ca.Monomial{big.NewInt(1)}}), fs)
		xqi := pow(x, qi).p
		xqi.Add(xqi, ca.NewPolynomial(k, fs.Order(), ca.PolynomialTerm[K]{Coefficient: neg1, Monomial: ca.Monomial{big.NewInt(1)}}))

		g, _, _, _, _ := gcd(fs, xqi)
		if !isOne(g) {
			s = append(s, factor[K]{A: g, N: new(big.Int).Set(i)})
			fs, _ = divide(fs, g)
		}
	}
	if !isOne(fs) {
		s = append(s, factor[K]{A: fs, N: deg(fs)})
	}
	if len(s) == 0 {
		s = append(s, factor[K]{A: f, N: big.NewInt(1)})
	}
	return s
}

func squareFree[K finite[K]](f *ca.Polynomial[K]) []factor[K] {
	c, _, _, _, _ := gcd(f, differentiate(f))
	w, _ := divide(f, c)

	r := make([]factor[K], 0)
	i := big.NewInt(1)
	for !isOne(w) {
		y, _, _, _, _ := gcd(w, c)
		fac, _ := divide(w, y)
		if !isOne(fac) {
			r = append(r, factor[K]{A: fac, N: new(big.Int).Set(i)})
		}
		w = y
		c, _ = divide(c, y)
		i.Add(i, big.NewInt(1))
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
				w[i] = big.NewInt(0).Div(cw[i], p)
			}

			// c -> c^{p^(e-1)}
			e1 := new(big.Int).Sub(e, big.NewInt(1))
			pe1 := pow(ca.NewRat(0, 1).SetFrac(p, big.NewInt(1)), e1).Num()
			rc := pow(cc, pe1)

			cRoot.Add(cRoot, ca.NewPolynomial(k, c.Order(), ca.PolynomialTerm[K]{Coefficient: rc, Monomial: w}))
		}

		// Do factorization on c^{1/p}.
		sf := squareFree(cRoot)
		for i := range sf {
			sf[i].N.Mul(sf[i].N, p)
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
		for i := big.NewInt(0); i.Cmp(aw[0]) < 0; i.Add(i, big.NewInt(1)) {
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

// gcd returns the greatest common divisor of a and b.
// gcd(a, b) = g = u*a + v*b
// a = a1*g
// b = b1*g
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

func randPoly[K finite[K]](poly *ca.Polynomial[K], n *big.Int) *ca.Polynomial[K] {
	k := poly.Field()
	// q is the order of the field k.
	q := pow(ca.NewRat(0, 1).SetFrac(k.characteristic(), big.NewInt(1)), k.primePower()).Num()

	for i := big.NewInt(0); i.Cmp(n) < 0; i.Add(i, big.NewInt(1)) {
		cint, _ := rand.Int(rand.Reader, q)
		if cint.Cmp(big.NewInt(0)) == 0 {
			continue
		}

		c := k.ith(cint)
		w := ca.Monomial{new(big.Int).Set(i)}
		poly.Add(poly, ca.NewPolynomial(k, poly.Order(), ca.PolynomialTerm[K]{Coefficient: c, Monomial: w}))
	}

	return poly
}

func deg[K ca.Field[K]](a *ca.Polynomial[K]) *big.Int {
	w := a.LeadingTerm().Monomial
	if len(w) == 0 {
		return big.NewInt(0)
	}
	return new(big.Int).Set(w[0])
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

type polynomialMod[K ca.Field[K]] struct {
	p   *ca.Polynomial[K]
	mod *ca.Polynomial[K]
}

func newPolynomialMod[K ca.Field[K]](p, mod *ca.Polynomial[K]) *polynomialMod[K] {
	z := &polynomialMod[K]{p: ca.NewPolynomial(p.Field(), p.Order()).Set(p), mod: ca.NewPolynomial(mod.Field(), mod.Order()).Set(mod)}
	_, z.p = ca.Divide(nil, z.p, []*ca.Polynomial[K]{z.mod})
	return z
}

func (x *polynomialMod[K]) NewOne() *polynomialMod[K] {
	z := &polynomialMod[K]{p: ca.NewPolynomial(x.p.Field(), x.p.Order(), ca.PolynomialTerm[K]{Coefficient: x.p.Field().NewOne()}), mod: ca.NewPolynomial(x.mod.Field(), x.mod.Order()).Set(x.mod)}
	z.p.SymbolStringer = x.p.SymbolStringer
	return z
}

func (z *polynomialMod[K]) Equal(x *polynomialMod[K]) bool {
	return z.p.Equal(x.p) && z.mod.Equal(x.mod)
}

func (z *polynomialMod[K]) Mul(x, y *polynomialMod[K]) *polynomialMod[K] {
	z.p.Mul(x.p, y.p)
	z.mod.Set(x.mod)
	_, z.p = ca.Divide(nil, z.p, []*ca.Polynomial[K]{z.mod})
	return z
}

func (z *polynomialMod[K]) Inv(x *polynomialMod[K]) *polynomialMod[K] {
	g, u, _, _, _ := gcd(x.p, x.mod)
	one := ca.NewPolynomial(x.p.Field(), x.p.Order(), ca.PolynomialTerm[K]{Coefficient: x.p.Field().NewOne()})
	if !g.Equal(one) {
		panic("inverse does not exist")
	}
	z.mod.Set(x.mod)
	z.p.Set(u)
	return z
}

func (x *polynomialMod[K]) String() string {
	return fmt.Sprintf("%s/%s", x.p, x.mod)
}

func pow[K ca.Group[K]](x K, n *big.Int) K {
	switch {
	case n.Sign() < 0:
		return pow(x.NewOne().Inv(x), new(big.Int).Neg(n))
	case n.Sign() == 0:
		return x.NewOne()
	case new(big.Int).Rem(n, big.NewInt(2)).Sign() == 0:
		return pow(x.NewOne().Mul(x, x), new(big.Int).Div(n, big.NewInt(2)))
	default:
		half := new(big.Int).Div(new(big.Int).Sub(n, big.NewInt(1)), big.NewInt(2))
		return x.NewOne().Mul(x, pow(x.NewOne().Mul(x, x), half))
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
