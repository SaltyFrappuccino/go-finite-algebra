package linalg

import (
	"testing"

	"nct/internal/field"
)

func TestDemoInvariants(t *testing.T) {
	for _, order := range []int{2, 3, 4} {
		for dim := 1; dim <= 4; dim++ {
			d := BuildDemo(order, dim)
			if !d.BasisAValid || !d.BasisBValid {
				t.Fatalf("basis check failed for F%d dim %d", order, dim)
			}
			if !d.BasisTransitionsInvert || !d.CoordinateTransitionsInvert {
				t.Fatalf("transition inverse check failed for F%d dim %d", order, dim)
			}
			if !d.VectorDecoded {
				t.Fatalf("decode check failed for F%d dim %d", order, dim)
			}
			if !d.VectorContravariant {
				t.Fatalf("vector contravariance failed for F%d dim %d", order, dim)
			}
			if !d.CovectorCovariant {
				t.Fatalf("covector covariance failed for F%d dim %d", order, dim)
			}
			if !d.CovectorEvaluationInvariant {
				t.Fatalf("covector value invariant failed for F%d dim %d", order, dim)
			}
			if !d.DualAConjugate || !d.DualBConjugate {
				t.Fatalf("dual basis check failed for F%d dim %d", order, dim)
			}
		}
	}
}

func TestBasisTransitionOrientation(t *testing.T) {
	s := NewSpace(field.New(3), 3)
	a := s.StandardBasis()
	b := []Vector{{1, 1, 0}, {0, 1, 1}, {0, 0, 1}}
	c, ok := s.BasisTransition(a, b)
	if !ok {
		t.Fatal("basis transition not found")
	}
	want := [][]int{{1, 0, 0}, {1, 1, 0}, {0, 1, 1}}
	for i := range want {
		for j := range want[i] {
			if c.At(i, j) != want[i][j] {
				t.Fatalf("C[%d][%d] = %d, want %d", i, j, c.At(i, j), want[i][j])
			}
		}
	}
	coord, ok := s.CoordinateTransition(a, b)
	if !ok {
		t.Fatal("coordinate transition not found")
	}
	inv, ok := c.Inverse()
	if !ok {
		t.Fatal("basis transition should be invertible")
	}
	if coord.Format() != inv.Format() {
		t.Fatalf("coordinate transition A->B must be C^-1\ngot:\n%s\nwant:\n%s", coord.Format(), inv.Format())
	}
}

func TestCoordinates(t *testing.T) {
	s := NewSpace(field.New(3), 3)
	b := []Vector{{1, 1, 0}, {0, 1, 1}, {0, 0, 1}}
	v := Vector{2, 1, 2}
	if !s.IsBasis(b) {
		t.Fatal("expected basis")
	}
	c, ok := s.Coordinates(v, b)
	if !ok {
		t.Fatal("coordinates not found")
	}
	got := MatrixFromColumns(s.F, b).MulVec(c)
	if !EqualVector(got, v) {
		t.Fatalf("coordinates reconstruct %v, want %v", got, v)
	}
}
