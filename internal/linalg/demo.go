package linalg

import "nct/internal/field"

type Demo struct {
	Space               Space
	BasisA              []Vector
	BasisB              []Vector
	DualA               []Vector
	DualB               []Vector
	Vector              Vector
	VectorCoordsA       Vector
	VectorCoordsB       Vector
	Covector            Vector
	CovectorCoordsA     Vector
	CovectorCoordsB     Vector
	TransitionAToB      Matrix
	TransitionBToA      Matrix
	Isotropic           []Vector
	BasisAValid         bool
	BasisBValid         bool
	TransitionsInvert   bool
	VectorContravariant bool
	CovectorCovariant   bool
	DualAConjugate      bool
	DualBConjugate      bool
}

func BuildDemo(order, dim int) Demo {
	s := NewSpace(field.New(order), dim)
	a := s.StandardBasis()
	b := s.TriangularBasis()
	v := make(Vector, dim)
	cov := make(Vector, dim)
	for i := 0; i < dim; i++ {
		v[i] = (i + 1) % s.F.Order()
		cov[i] = s.F.One()
	}
	if zeroVector(s.F, v) {
		v[0] = s.F.One()
	}
	coordsA, _ := s.Coordinates(v, a)
	coordsB, _ := s.Coordinates(v, b)
	tAB, _ := s.Transition(a, b)
	tBA, _ := s.Transition(b, a)
	dualA, _ := s.DualBasis(a)
	dualB, _ := s.DualBasis(b)
	covA, _ := s.CovectorCoordinates(cov, dualA)
	covB, _ := s.CovectorCoordinates(cov, dualB)
	inv, _ := tBA.Inverse()
	return Demo{
		Space:               s,
		BasisA:              a,
		BasisB:              b,
		DualA:               dualA,
		DualB:               dualB,
		Vector:              v,
		VectorCoordsA:       coordsA,
		VectorCoordsB:       coordsB,
		Covector:            cov,
		CovectorCoordsA:     covA,
		CovectorCoordsB:     covB,
		TransitionAToB:      tAB,
		TransitionBToA:      tBA,
		Isotropic:           s.IsotropicVectors(s.IdentityGram()),
		BasisAValid:         s.IsBasis(a),
		BasisBValid:         s.IsBasis(b),
		TransitionsInvert:   tAB.Mul(tBA).IsIdentity() && tBA.Mul(tAB).IsIdentity(),
		VectorContravariant: EqualVector(coordsB, inv.MulVec(coordsA)),
		CovectorCovariant:   EqualVector(covB, tBA.Transpose().MulVec(covA)),
		DualAConjugate:      dualIsConjugate(s, a, dualA),
		DualBConjugate:      dualIsConjugate(s, b, dualB),
	}
}

func dualIsConjugate(s Space, basis []Vector, dual []Vector) bool {
	if len(basis) != s.Dim || len(dual) != s.Dim {
		return false
	}
	for i := 0; i < s.Dim; i++ {
		for j := 0; j < s.Dim; j++ {
			want := s.F.Zero()
			if i == j {
				want = s.F.One()
			}
			if s.Eval(dual[i], basis[j]) != want {
				return false
			}
		}
	}
	return true
}
