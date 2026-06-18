package finite

import "testing"

func TestSemigroupEnumeratorAgainstBruteForce(t *testing.T) {
	for _, n := range []int{2, 3} {
		got := len(Semigroups(n))
		want := bruteForceSemigroups(n)
		if got != want {
			t.Fatalf("order %d semigroups: got %d want %d", n, got, want)
		}
	}
}

func TestCommutativeMonoidEnumeratorAgainstBruteForce(t *testing.T) {
	for _, n := range []int{2, 3} {
		got := len(CommutativeMonoids(n))
		want := bruteForceCommutativeMonoids(n)
		if got != want {
			t.Fatalf("order %d commutative monoids: got %d want %d", n, got, want)
		}
	}
}

func bruteForceSemigroups(n int) int {
	total := PowUint64(n, n*n)
	count := 0
	for code := uint64(0); code < total; code++ {
		if opFromCode(n, code).IsAssociative() {
			count++
		}
	}
	return count
}

func bruteForceCommutativeMonoids(n int) int {
	total := PowUint64(n, n*n)
	count := 0
	for code := uint64(0); code < total; code++ {
		op := opFromCode(n, code)
		if op.IsAssociative() && op.IsCommutative() {
			if _, ok := op.Identity(); ok {
				count++
			}
		}
	}
	return count
}

func opFromCode(n int, code uint64) Op {
	values := make([]int, n*n)
	x := code
	for i := range values {
		values[i] = int(x % uint64(n))
		x /= uint64(n)
	}
	return NewOp(n, values)
}
