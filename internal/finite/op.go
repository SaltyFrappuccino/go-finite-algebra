package finite

import (
	"fmt"
	"strings"
)

type Op struct {
	N int
	T []int
}

func NewOp(n int, values []int) Op {
	t := make([]int, len(values))
	copy(t, values)
	return Op{N: n, T: t}
}

func (o Op) At(a, b int) int {
	return o.T[a*o.N+b]
}

func (o Op) Rows() [][]int {
	rows := make([][]int, o.N)
	for i := 0; i < o.N; i++ {
		rows[i] = make([]int, o.N)
		copy(rows[i], o.T[i*o.N:(i+1)*o.N])
	}
	return rows
}

func (o Op) IsAssociative() bool {
	for a := 0; a < o.N; a++ {
		for b := 0; b < o.N; b++ {
			for c := 0; c < o.N; c++ {
				if o.At(o.At(a, b), c) != o.At(a, o.At(b, c)) {
					return false
				}
			}
		}
	}
	return true
}

func (o Op) IsCommutative() bool {
	for a := 0; a < o.N; a++ {
		for b := a + 1; b < o.N; b++ {
			if o.At(a, b) != o.At(b, a) {
				return false
			}
		}
	}
	return true
}

func (o Op) Identity() (int, bool) {
	for e := 0; e < o.N; e++ {
		ok := true
		for a := 0; a < o.N; a++ {
			if o.At(e, a) != a || o.At(a, e) != a {
				ok = false
				break
			}
		}
		if ok {
			return e, true
		}
	}
	return 0, false
}

func (o Op) Absorbs(z int) bool {
	for a := 0; a < o.N; a++ {
		if o.At(z, a) != z || o.At(a, z) != z {
			return false
		}
	}
	return true
}

func (o Op) IsGroup() bool {
	e, ok := o.Identity()
	if !ok {
		return false
	}
	for a := 0; a < o.N; a++ {
		found := false
		for b := 0; b < o.N; b++ {
			if o.At(a, b) == e && o.At(b, a) == e {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (o Op) IsAbelianGroup() bool {
	return o.IsGroup() && o.IsCommutative()
}

func (o Op) FormatTable() string {
	var b strings.Builder
	for i := 0; i < o.N; i++ {
		for j := 0; j < o.N; j++ {
			if j > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(fmt.Sprintf("%d", o.At(i, j)))
		}
		if i+1 < o.N {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func PowUint64(base, exp int) uint64 {
	result := uint64(1)
	for i := 0; i < exp; i++ {
		result *= uint64(base)
	}
	return result
}
