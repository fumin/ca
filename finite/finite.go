package finite

import "math/big"

type prime struct {
	// Order is the number of elements in the finite field which must be a prime number.
	order *big.Int
	// I is the integer representation of the element.
	i *big.Int
}

func newPrime(p, i int) *prime {
	return &prime{order: big.NewInt(int64(p)), i: big.NewInt(int64(i))}
}

func (x *prime) NewZero() *prime {
	return &prime{order: big.NewInt(0).Set(x.order), i: big.NewInt(0)}
}

func (x *prime) NewOne() *prime {
	return &prime{order: big.NewInt(0).Set(x.order), i: big.NewInt(1)}
}

func (x *prime) Equal(y *prime) bool {
	return (x.order.Cmp(y.order) == 0) && (x.i.Cmp(y.i) == 0)
}

func (z *prime) Add(x, y *prime) *prime {
	z.order.Set(x.order)
	z.i.Add(x.i, y.i)
	z.i.Mod(z.i, z.order)
	return z
}

func (z *prime) Sub(x, y *prime) *prime {
	z.order.Set(x.order)
	z.i.Sub(x.i, y.i)
	z.i.Mod(z.i, z.order)
	return z
}

func (z *prime) Mul(x, y *prime) *prime {
	z.order.Set(x.order)
	z.i.Mul(x.i, y.i)
	z.i.Mod(z.i, z.order)
	return z
}

func (z *prime) Div(x, y *prime) *prime {
	z.Inv(y)
	z.i.Mul(x.i, z.i)
	z.i.Mod(z.i, z.order)
	return z
}

func (z *prime) Inv(x *prime) *prime {
	if x.i.Sign() == 0 {
		panic("division by zero")
	}
	z.order.Set(x.order)
	z.i.ModInverse(x.i, z.order)
	return z
}

func (x *prime) String() string {
	return x.i.String()
}

func (x *prime) characteristic() int {
	return int(x.order.Int64())
}

func (x *prime) primePower() int { return 1 }
