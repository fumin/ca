package ca

import (
	"fmt"
	"iter"
	"math/big"
	"reflect"
	"slices"
	"strings"

	"github.com/jba/omap/ordered"
)

// A Field is an element whose addition and multiplication operations satisfy the [field] axioms.
//
// [field]: https://en.wikipedia.org/wiki/Field_(mathematics)
type Field[T any] interface {
	// NewZero returns the additive identity of the field.
	NewZero() T
	// NewOne returns the multiplicative identity of the field.
	NewOne() T

	// Equal reports whether x and y are equal, where x is the method receiver.
	Equal(y T) bool
	// Add sets z to the sum x+y and returns z, where z is the method receiver.
	Add(x, y T) T
	// Sub sets z to the difference x-y and returns z, where z is the method receiver.
	Sub(x, y T) T
	// Mul sets z to the product x*y and returns z, where z is the method receiver.
	Mul(x, y T) T
	// Div sets z to the quotient x/y and returns z, where z is the method receiver.
	Div(x, y T) T
	// Inv sets z to 1/x and returns z, where z is the method receiver.
	Inv(x T) T

	// String returns the string representation.
	String() string
}

// A Rat represents a quotient of arbitrary precision.
type Rat struct{ *big.Rat }

// NewRat creates a new [Rat] with numerator a and denominator b.
func NewRat(a, b int64) *Rat { return &Rat{big.NewRat(a, b)} }

// NewZero returns the additive identity 0.
func (x *Rat) NewZero() *Rat {
	return &Rat{big.NewRat(0, 1)}
}

// NewOne returns the multiplicative identity 1.
func (x *Rat) NewOne() *Rat {
	return &Rat{big.NewRat(1, 1)}
}

// Add sets z to the sum x+y and returns z.
func (z *Rat) Add(x, y *Rat) *Rat { return &Rat{z.Rat.Add(x.Rat, y.Rat)} }

// Sub sets z to the difference x-y and returns z.
func (z *Rat) Sub(x, y *Rat) *Rat { return &Rat{z.Rat.Sub(x.Rat, y.Rat)} }

// Mul sets z to the product x*y and returns z.
func (z *Rat) Mul(x, y *Rat) *Rat { return &Rat{z.Rat.Mul(x.Rat, y.Rat)} }

// Div sets z to the quotient x/y and returns z. If y == 0, Div panics.
func (z *Rat) Div(x, y *Rat) *Rat { return &Rat{z.Rat.Quo(x.Rat, y.Rat)} }

// Inv sets z to 1/x and returns z. If x == 0, Inv panics.
func (z *Rat) Inv(x *Rat) *Rat { return &Rat{z.Rat.Inv(x.Rat)} }

// Equal reports whether x and y are equal.
func (x *Rat) Equal(y *Rat) bool {
	return x.Rat.Cmp(y.Rat) == 0
}

// String returns a string representation of x in the form "a/b" if b != 1, and in the form "a" if b == 1.
func (x *Rat) String() string {
	return x.RatString()
}

// A [Monomial] is a product of variables chained by multiplication.
//
// [Monomial]: https://en.wikipedia.org/wiki/Monomial
type Monomial []*big.Int

func (m Monomial) deg() *big.Int {
	deg := big.NewInt(0)
	for _, e := range m {
		deg.Add(deg, e)
	}
	return deg
}

func (z Monomial) mul(x Monomial) Monomial {
	for range len(x) - len(z) {
		z = append(z, big.NewInt(0))
	}
	for s, xe := range x {
		z[s].Add(z[s], xe)
	}
	return z
}

func (z Monomial) divide(x, y Monomial) (Monomial, bool) {
	if len(y) > len(x) {
		return nil, false
	}
	for s, xe := range x {
		if s >= len(z) {
			z = append(z, big.NewInt(0))
		}
		var ye *big.Int
		if s < len(y) {
			ye = y[s]
		} else {
			ye = big.NewInt(0)
		}

		z[s].Sub(xe, ye)
		if z[s].Cmp(big.NewInt(0)) == -1 {
			return nil, false
		}
	}
	return z, true
}

// An Order is a [monomial order] for comparing monomials.
// The meaning of the return value is the same as [cmp.Compare].
//
// [monomial order]: https://en.wikipedia.org/wiki/Monomial_order
type Order func(x, y Monomial) int

// Lex compares x and y lexicographically.
func Lex(x, y Monomial) int {
	for i := range x {
		if i >= len(y) {
			return 1
		}
		if c := x[i].Cmp(y[i]); c != 0 {
			return -1 * c
		}
	}
	if len(x) == len(y) {
		return 0
	}
	return -1
}

// [Deglex] compares x, y by first comparing their degrees, and in case of a tie applies the lexicographic order.
//
// [Deglex]: https://en.wikipedia.org/wiki/Monomial_order#Graded_lexicographic_order
func Deglex(x, y Monomial) int {
	if c := x.deg().Cmp(y.deg()); c != 0 {
		return c
	}
	return Lex(x, y)
}

// A PolynomialTerm is a term in a polynomial.
type PolynomialTerm[K Field[K]] struct {
	Coefficient K
	Monomial    Monomial
}

type Symbol int

// A Polynomial is a polynomial of noncommutative variables.
type Polynomial[K Field[K]] struct {
	// SymbolStringer specifies how a symbol in a monomial is formated when the polynomial is printed out.
	SymbolStringer func(s Symbol) string

	field K
	order Order
	m     *ordered.Map[Monomial, K]
}

// NewPolynomial returns a new polynomial containing the given terms.
func NewPolynomial[K Field[K]](field K, order Order, terms ...PolynomialTerm[K]) *Polynomial[K] {
	x := &Polynomial[K]{
		SymbolStringer: englishSymbolStringer,
		field:          field,
		order:          order,
		m:              ordered.NewMap[Monomial, K](order),
	}
	for _, term := range terms {
		x.addTerm(1, term)
	}
	return x
}

// Field returns the field of the coefficients in x.
func (x *Polynomial[K]) Field() K { return x.field }

// Order returns the monomial order employed by x.
func (x *Polynomial[K]) Order() Order { return x.order }

// Len reports the number of terms in x.
func (x *Polynomial[K]) Len() int { return x.m.Len() }

// Terms iterates the terms in a polynomial.
func (x *Polynomial[K]) Terms() iter.Seq2[K, Monomial] {
	return func(yield func(K, Monomial) bool) {
		for w, c := range x.m.Backward() {
			if !yield(c, w) {
				return
			}
		}
	}
}

// Equal reports whether x and y have the same coefficients and monomials.
func (x *Polynomial[K]) Equal(y *Polynomial[K]) bool {
	if x.m.Len() != y.m.Len() {
		return false
	}
	for i := range x.m.Len() {
		xw, xc := x.m.Nth(x.m.Len() - 1 - i)
		yw, yc := y.m.Nth(y.m.Len() - 1 - i)
		if !monomialEq(xw, yw) {
			return false
		}
		if !xc.Equal(yc) {
			return false
		}
	}
	return true
}

// Set sets z to x and returns z.
func (z *Polynomial[K]) Set(x *Polynomial[K]) *Polynomial[K] {
	if z == x {
		return z
	}
	z.SymbolStringer = x.SymbolStringer
	z.field = x.field
	z.order = x.order
	z.m = ordered.NewMap[Monomial, K](z.order)
	for xw, xc := range x.m.All() {
		w := makeMonomial(len(xw))
		copyMonomial(w, xw)
		z.addTerm(1, PolynomialTerm[K]{Coefficient: xc, Monomial: w})
	}
	return z
}

// Add sets z to the sum x+y and returns z.
func (z *Polynomial[K]) Add(x, y *Polynomial[K]) *Polynomial[K] {
	// Set z = x, while handling the case where x or y is z itself.
	if y == z {
		x, y = y, x
	}
	if z != x {
		z.m.Clear()
		for xw, c := range x.m.All() {
			w := makeMonomial(len(xw))
			copyMonomial(w, xw)
			z.addTerm(1, PolynomialTerm[K]{Coefficient: c, Monomial: w})
		}
	}

	// Compute z += y.
	for yw, c := range y.m.All() {
		w := makeMonomial(len(yw))
		copyMonomial(w, yw)
		z.addTerm(1, PolynomialTerm[K]{Coefficient: c, Monomial: w})
	}

	return z
}

// Mul sets z to the product x*y and returns z.
func (z *Polynomial[K]) Mul(x, y *Polynomial[K]) *Polynomial[K] {
	if z == x {
		panic(fmt.Sprintf("z == x"))
	}
	if z == y {
		panic(fmt.Sprintf("z == y"))
	}

	z.m.Clear()
	for xw, xc := range x.m.Backward() {
		for yw, yc := range y.m.Backward() {
			c := z.field.Mul(xc, yc)
			var w Monomial
			w = w.mul(xw).mul(yw)
			z.addTerm(1, PolynomialTerm[K]{Coefficient: c, Monomial: w})
		}
	}

	return z
}

// Pow sets z to the power x^y and returns z.
func (z *Polynomial[K]) Pow(x *Polynomial[K], y int) *Polynomial[K] {
	if z == x {
		panic("z == x")
	}

	z.Set(x)
	buf := NewPolynomial[K](z.field, z.order)
	for range y - 1 {
		buf.Mul(z, x)
		z, buf = buf, z
	}
	if y%2 == 0 {
		z, buf = buf, z
		z.Set(buf)
	}

	return z
}

// LeadingTerm returns the polynomial term of the leading monomial.
// Note that the leading term depends on the monomial order employed by the polynomial.
func (x *Polynomial[K]) LeadingTerm() PolynomialTerm[K] {
	w, _, ok := x.m.Max()
	if !ok {
		panic("zero polynomial has no terms")
	}
	c, _ := x.m.Get(w)
	return PolynomialTerm[K]{Coefficient: c, Monomial: w}
}

// String returns the string representation of x.
// Symbols in x are formatted using x.SymbolStringer.
func (x *Polynomial[K]) String() string {
	if x.Len() == 0 {
		return "0"
	}
	var b strings.Builder
	for i := range x.m.Len() {
		w, c := x.m.Nth(x.m.Len() - 1 - i)

		// Print c.
		s := c.String()
		if s[0] != '-' {
			s = "+" + s
		}
		switch {
		case i == 0 && s == "+1" && len(w) != 0:
			s = ""
		case i == 0 && s[0] == '+':
			s = s[1:]
		case s == "+1" && len(w) != 0:
			s = "+"
		case s == "-1" && len(w) != 0:
			s = "-"
		}
		fmt.Fprintf(&b, "%s", s)

		// Print w.
		printMonomial(&b, w, x.SymbolStringer)
	}
	return b.String()
}

func (x *Polynomial[K]) addTerm(sign int, term PolynomialTerm[K]) {
	term.Monomial = compactMonomial(term.Monomial)
	c, ok := x.m.Get(term.Monomial)
	if !ok {
		c = x.field.NewZero()
	}

	tc := term.Coefficient
	tcv := reflect.ValueOf(tc)
	kind := tcv.Kind()
	if (kind == reflect.Pointer || kind == reflect.Interface) && tcv.IsNil() {
		tc = x.field.NewOne()
	}
	if sign < 0 {
		c.Sub(c, tc)
	} else {
		c.Add(c, tc)
	}

	if c.Equal(x.field.NewZero()) {
		x.m.Delete(term.Monomial)
	} else {
		x.m.Set(term.Monomial, c)
	}
}

func (z *Polynomial[K]) add(sign int, c K, m Monomial, x *Polynomial[K]) {
	for xw, xc := range x.m.Backward() {
		c := z.field.Mul(c, xc)
		var w Monomial
		w = w.mul(m).mul(xw)
		z.addTerm(sign, PolynomialTerm[K]{Coefficient: c, Monomial: w})
	}
}

func (z *Polynomial[K]) mulScalar(scalar K, x *Polynomial[K]) *Polynomial[K] {
	if z == x {
		for zw, zc := range z.m.All() {
			zc.Mul(scalar, zc)
			z.m.Set(zw, zc)
		}
		return z
	}

	z.m.Clear()
	for xw, xc := range x.m.All() {
		c := z.field.Mul(scalar, xc)
		w := makeMonomial(len(xw))
		copyMonomial(w, xw)
		z.addTerm(1, PolynomialTerm[K]{Coefficient: c, Monomial: w})
	}
	return z
}

// A Quotient is the resulting quotient of a polynomial division.
type Quotient[K Field[K]] struct {
	Coefficient K
	Monomial    Monomial
}

// Divide divides the polynomial f by the ideal g, and returns the quotient and remainder.
// The polynomial f is modified upon return.
func Divide[K Field[K]](quotient [][]Quotient[K], f *Polynomial[K], g []*Polynomial[K]) (outQuotient [][]Quotient[K], remainder *Polynomial[K]) {
	if quotient != nil {
		short := len(g) - len(quotient)
		if short > 0 {
			quotient = append(quotient, make([][]Quotient[K], short)...)
		}
		quotient = quotient[:len(g)]
		for i := range quotient {
			quotient[i] = quotient[i][:0]
		}
	}
	p := NewPolynomial[K](f.field, f.order)
	p.SymbolStringer = f.SymbolStringer
	v := f

	for v.m.Len() != 0 {
		lmv := v.LeadingTerm()
		ltv := lmv.Monomial

		// Find basis where ltv = ltq * ltg.
		basis, ltq, lmg := -1, Monomial{}, PolynomialTerm[K]{}
		for i, gi := range g {
			if gi == nil {
				continue
			}
			lmg = gi.LeadingTerm()
			ltg := lmg.Monomial

			var ok bool
			ltq, ok = ltq.divide(ltv, ltg)
			if ok {
				basis = i
				break
			}
		}

		if basis == -1 {
			p.addTerm(1, lmv)
			v.addTerm(-1, lmv)
		} else {
			q := Quotient[K]{
				Coefficient: f.field.NewZero().Div(lmv.Coefficient, lmg.Coefficient),
				Monomial:    ltq,
			}
			if quotient != nil {
				quotient[basis] = append(quotient[basis], q)
			}
			v.add(-1, q.Coefficient, q.Monomial, g[basis])
		}
	}

	return quotient, p
}

func newMonomial(es ...int) Monomial {
	m := make(Monomial, len(es))
	for i, e := range es {
		m[i] = big.NewInt(int64(e))
	}
	return m
}

func makeMonomial(size int) Monomial {
	m := make(Monomial, size)
	for i := range m {
		m[i] = big.NewInt(0)
	}
	return m
}

func copyMonomial(dst, src Monomial) {
	for i := range min(len(dst), len(src)) {
		dst[i].Set(src[i])
	}
}

func compactMonomial(x Monomial) Monomial {
	nonzero := -1
	for i := len(x) - 1; i >= 0; i-- {
		if x[i].Cmp(big.NewInt(0)) != 0 {
			nonzero = i
			break
		}
	}
	return x[:nonzero+1]
}

func monomialEq(x, y Monomial) bool {
	return slices.EqualFunc(x, y, func(a, b *big.Int) bool { return a.Cmp(b) == 0 })
}

func englishSymbolStringer(s Symbol) string {
	return string(byte(s) + 'a')
}

func printSymbol(b *strings.Builder, s Symbol, power *big.Int, symbolStringer func(Symbol) string) {
	v := symbolStringer(s)
	switch {
	case power.Cmp(big.NewInt(1)) == 0:
		fmt.Fprintf(b, "%s", v)
	default:
		fmt.Fprintf(b, "%s^%d", v, power)
	}
}

func printMonomial(b *strings.Builder, w Monomial, ss func(Symbol) string) {
	for s, e := range w {
		printSymbol(b, Symbol(s), e, ss)
	}
}
