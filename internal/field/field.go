package field

import "fmt"

type Field interface {
	Order() int
	Zero() int
	One() int
	Add(a, b int) int
	Neg(a int) int
	Mul(a, b int) int
	Inv(a int) (int, bool)
	Name() string
	ElementName(a int) string
}

func New(order int) Field {
	switch order {
	case 2, 3:
		return primeField{p: order}
	case 4:
		return f4{}
	default:
		panic(fmt.Sprintf("field of order %d is not available", order))
	}
}

func Sub(f Field, a, b int) int {
	return f.Add(a, f.Neg(b))
}

func Div(f Field, a, b int) (int, bool) {
	inv, ok := f.Inv(b)
	if !ok {
		return 0, false
	}
	return f.Mul(a, inv), true
}

type primeField struct {
	p int
}

func (f primeField) Order() int {
	return f.p
}

func (f primeField) Zero() int {
	return 0
}

func (f primeField) One() int {
	return 1
}

func (f primeField) Add(a, b int) int {
	return (a + b) % f.p
}

func (f primeField) Neg(a int) int {
	if a == 0 {
		return 0
	}
	return f.p - a
}

func (f primeField) Mul(a, b int) int {
	return (a * b) % f.p
}

func (f primeField) Inv(a int) (int, bool) {
	if a == 0 {
		return 0, false
	}
	for x := 1; x < f.p; x++ {
		if f.Mul(a, x) == 1 {
			return x, true
		}
	}
	return 0, false
}

func (f primeField) Name() string {
	return fmt.Sprintf("F%d", f.p)
}

func (f primeField) ElementName(a int) string {
	return fmt.Sprintf("%d", a)
}

type f4 struct{}

func (f f4) Order() int {
	return 4
}

func (f f4) Zero() int {
	return 0
}

func (f f4) One() int {
	return 1
}

func (f f4) Add(a, b int) int {
	return a ^ b
}

func (f f4) Neg(a int) int {
	return a
}

func (f f4) Mul(a, b int) int {
	constPart := 0
	alphaPart := 0
	alphaSquarePart := 0
	if a&1 != 0 && b&1 != 0 {
		constPart ^= 1
	}
	if a&1 != 0 && b&2 != 0 {
		alphaPart ^= 1
	}
	if a&2 != 0 && b&1 != 0 {
		alphaPart ^= 1
	}
	if a&2 != 0 && b&2 != 0 {
		alphaSquarePart ^= 1
	}
	if alphaSquarePart != 0 {
		constPart ^= 1
		alphaPart ^= 1
	}
	return constPart | alphaPart<<1
}

func (f f4) Inv(a int) (int, bool) {
	if a == 0 {
		return 0, false
	}
	for x := 1; x < 4; x++ {
		if f.Mul(a, x) == 1 {
			return x, true
		}
	}
	return 0, false
}

func (f f4) Name() string {
	return "F4"
}

func (f f4) ElementName(a int) string {
	switch a {
	case 0:
		return "0"
	case 1:
		return "1"
	case 2:
		return "a"
	case 3:
		return "a+1"
	default:
		return fmt.Sprintf("%d", a)
	}
}
