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
			if !d.TransitionsInvert {
				t.Fatalf("transition inverse check failed for F%d dim %d", order, dim)
			}
			if !d.VectorContravariant {
				t.Fatalf("vector contravariance failed for F%d dim %d", order, dim)
			}
			if !d.CovectorCovariant {
				t.Fatalf("covector covariance failed for F%d dim %d", order, dim)
			}
			if !d.DualAConjugate || !d.DualBConjugate {
				t.Fatalf("dual basis check failed for F%d dim %d", order, dim)
			}
		}
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
