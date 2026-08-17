package finite

import (
	"flag"
	"fmt"
	"log"
	"testing"
)

func TestPrime(t *testing.T) {
	tests := []struct {
		order int
		mul   [][]int
	}{
		{
			order: 13,
			mul: [][]int{
				{0},
				{0, 1},
				{0, 2, 4},
				{0, 3, 6, 9},
				{0, 4, 8, 12, 3},
				{0, 5, 10, 2, 7, 12},
				{0, 6, 12, 5, 11, 4, 10},
				{0, 7, 1, 8, 2, 9, 3, 10},
				{0, 8, 3, 11, 6, 1, 9, 4, 12},
				{0, 9, 5, 1, 10, 6, 2, 11, 7, 3},
				{0, 10, 7, 4, 1, 11, 8, 5, 2, 12, 9},
				{0, 11, 9, 7, 5, 3, 1, 12, 10, 8, 6, 4},
				{0, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			},
		},
		{
			order: 31,
			mul: [][]int{
				{0},
				{0, 1},
				{0, 2, 4},
				{0, 3, 6, 9},
				{0, 4, 8, 12, 16},
				{0, 5, 10, 15, 20, 25},
				{0, 6, 12, 18, 24, 30, 5},
				{0, 7, 14, 21, 28, 4, 11, 18},
				{0, 8, 16, 24, 1, 9, 17, 25, 2},
				{0, 9, 18, 27, 5, 14, 23, 1, 10, 19},
				{0, 10, 20, 30, 9, 19, 29, 8, 18, 28, 7},
				{0, 11, 22, 2, 13, 24, 4, 15, 26, 6, 17, 28},
				{0, 12, 24, 5, 17, 29, 10, 22, 3, 15, 27, 8, 20},
				{0, 13, 26, 8, 21, 3, 16, 29, 11, 24, 6, 19, 1, 14},
				{0, 14, 28, 11, 25, 8, 22, 5, 19, 2, 16, 30, 13, 27, 10},
				{0, 15, 30, 14, 29, 13, 28, 12, 27, 11, 26, 10, 25, 9, 24, 8},
				{0, 16, 1, 17, 2, 18, 3, 19, 4, 20, 5, 21, 6, 22, 7, 23, 8},
				{0, 17, 3, 20, 6, 23, 9, 26, 12, 29, 15, 1, 18, 4, 21, 7, 24, 10},
				{0, 18, 5, 23, 10, 28, 15, 2, 20, 7, 25, 12, 30, 17, 4, 22, 9, 27, 14},
				{0, 19, 7, 26, 14, 2, 21, 9, 28, 16, 4, 23, 11, 30, 18, 6, 25, 13, 1, 20},
				{0, 20, 9, 29, 18, 7, 27, 16, 5, 25, 14, 3, 23, 12, 1, 21, 10, 30, 19, 8, 28},
				{0, 21, 11, 1, 22, 12, 2, 23, 13, 3, 24, 14, 4, 25, 15, 5, 26, 16, 6, 27, 17, 7},
				{0, 22, 13, 4, 26, 17, 8, 30, 21, 12, 3, 25, 16, 7, 29, 20, 11, 2, 24, 15, 6, 28, 19},
				{0, 23, 15, 7, 30, 22, 14, 6, 29, 21, 13, 5, 28, 20, 12, 4, 27, 19, 11, 3, 26, 18, 10, 2},
				{0, 24, 17, 10, 3, 27, 20, 13, 6, 30, 23, 16, 9, 2, 26, 19, 12, 5, 29, 22, 15, 8, 1, 25, 18},
				{0, 25, 19, 13, 7, 1, 26, 20, 14, 8, 2, 27, 21, 15, 9, 3, 28, 22, 16, 10, 4, 29, 23, 17, 11, 5},
				{0, 26, 21, 16, 11, 6, 1, 27, 22, 17, 12, 7, 2, 28, 23, 18, 13, 8, 3, 29, 24, 19, 14, 9, 4, 30, 25},
				{0, 27, 23, 19, 15, 11, 7, 3, 30, 26, 22, 18, 14, 10, 6, 2, 29, 25, 21, 17, 13, 9, 5, 1, 28, 24, 20, 16},
				{0, 28, 25, 22, 19, 16, 13, 10, 7, 4, 1, 29, 26, 23, 20, 17, 14, 11, 8, 5, 2, 30, 27, 24, 21, 18, 15, 12, 9},
				{0, 29, 27, 25, 23, 21, 19, 17, 15, 13, 11, 9, 7, 5, 3, 1, 30, 28, 26, 24, 22, 20, 18, 16, 14, 12, 10, 8, 6, 4},
				{0, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			},
		},
	}

	for testI, test := range tests {
		t.Run(fmt.Sprintf("%d", testI), func(t *testing.T) {
			t.Parallel()
			add := make([][]int, test.order)
			sub := make([][]int, test.order)
			mul := make([][]int, test.order)
			div := make([][]int, test.order)
			for i := range test.order {
				add[i] = make([]int, test.order)
				sub[i] = make([]int, test.order)
				mul[i] = make([]int, test.order)
				div[i] = make([]int, test.order)
			}
			// Addition and subtraction.
			for j := range test.order {
				add[0][j] = j
				sub[add[0][j]][0] = j
			}
			for i := 1; i < test.order; i++ {
				for j := range test.order - 1 {
					add[i][j] = add[i-1][j+1]
					sub[add[i][j]][i] = j
				}
				add[i][test.order-1] = add[i-1][0]
				sub[add[i][test.order-1]][i] = test.order - 1
			}
			// Multiplication and division.
			for i := range test.order {
				for j := range test.order {
					if j <= i {
						mul[i][j] = test.mul[i][j]
					} else {
						mul[i][j] = test.mul[j][i]
					}
					div[mul[i][j]][i] = j
				}
			}

			// Check the arithmetics.
			for i := range test.order {
				for j := range test.order {
					x, y := newPrime(test.order, i), newPrime(test.order, j)
					if z := newPrime(0, 0).Add(x, y); !z.Equal(newPrime(test.order, add[i][j])) {
						t.Errorf("Add(%d %d): got %d want %d", x, y, z, add[i][j])
					}
					if z := newPrime(0, 0).Sub(x, y); !z.Equal(newPrime(test.order, sub[i][j])) {
						t.Errorf("Sub(%d %d): got %d want %d", x, y, z, sub[i][j])
					}
					if z := newPrime(0, 0).Mul(x, y); !z.Equal(newPrime(test.order, mul[i][j])) {
						t.Errorf("Mul(%d %d): got %d want %d", x, y, z, mul[i][j])
					}
					if j != 0 {
						if z := newPrime(0, 0).Div(x, y); !z.Equal(newPrime(test.order, div[i][j])) {
							t.Errorf("Div(%d %d): got %d want %d", x, y, z, div[i][j])
						}
					}
				}
			}
			for i := 1; i < test.order; i++ {
				x := newPrime(test.order, i)
				if z := newPrime(0, 0).Inv(x); !z.Equal(newPrime(test.order, div[1][i])) {
					t.Errorf("Inv(%d): got %d want %d", x, z, div[1][i])
				}
			}
		})
	}
}

func TestPrimeEqual(t *testing.T) {
	tests := []struct {
		x  *prime
		y  *prime
		eq bool
	}{
		{
			x:  newPrime(5, 2),
			y:  newPrime(5, 3),
			eq: false,
		},
		{
			x:  newPrime(11, 5),
			y:  newPrime(7, 5),
			eq: false,
		},
		{
			x:  newPrime(5, 3),
			y:  newPrime(5, 3),
			eq: true,
		},
	}

	for testI, test := range tests {
		t.Run(fmt.Sprintf("%d", testI), func(t *testing.T) {
			t.Parallel()
			if eq := test.x.Equal(test.y); eq != test.eq {
				t.Errorf("%d.Equal(%d): got %v want %v", test.x, test.y, eq, test.eq)
			}
		})
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	log.SetFlags(log.Lmicroseconds | log.Llongfile | log.LstdFlags)

	m.Run()
}
