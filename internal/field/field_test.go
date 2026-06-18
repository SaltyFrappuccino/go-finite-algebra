package field

import "testing"

func TestFieldAxioms(t *testing.T) {
	for _, order := range []int{2, 3, 4} {
		f := New(order)
		for a := 0; a < order; a++ {
			if f.Add(a, f.Zero()) != a || f.Add(f.Zero(), a) != a {
				t.Fatalf("%s additive identity failed for %d", f.Name(), a)
			}
			if f.Mul(a, f.One()) != a || f.Mul(f.One(), a) != a {
				t.Fatalf("%s multiplicative identity failed for %d", f.Name(), a)
			}
			if f.Add(a, f.Neg(a)) != f.Zero() {
				t.Fatalf("%s additive inverse failed for %d", f.Name(), a)
			}
			if a != f.Zero() {
				inv, ok := f.Inv(a)
				if !ok || f.Mul(a, inv) != f.One() || f.Mul(inv, a) != f.One() {
					t.Fatalf("%s multiplicative inverse failed for %d", f.Name(), a)
				}
			}
			for b := 0; b < order; b++ {
				if f.Add(a, b) != f.Add(b, a) {
					t.Fatalf("%s addition is not commutative", f.Name())
				}
				if f.Mul(a, b) != f.Mul(b, a) {
					t.Fatalf("%s multiplication is not commutative", f.Name())
				}
				for c := 0; c < order; c++ {
					if f.Add(f.Add(a, b), c) != f.Add(a, f.Add(b, c)) {
						t.Fatalf("%s addition is not associative", f.Name())
					}
					if f.Mul(f.Mul(a, b), c) != f.Mul(a, f.Mul(b, c)) {
						t.Fatalf("%s multiplication is not associative", f.Name())
					}
					if f.Mul(a, f.Add(b, c)) != f.Add(f.Mul(a, b), f.Mul(a, c)) {
						t.Fatalf("%s distributivity failed", f.Name())
					}
				}
			}
		}
	}
}

func TestF4PolynomialModel(t *testing.T) {
	f := New(4)
	alpha := 2
	if f.Mul(alpha, alpha) != f.Add(alpha, f.One()) {
		t.Fatalf("F4 model must satisfy a^2 = a + 1")
	}
}
