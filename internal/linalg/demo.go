package linalg

import "nct/internal/field"

type Demo struct {
	Space                       Space
	BasisA                      []Vector
	BasisB                      []Vector
	DualA                       []Vector
	DualB                       []Vector
	Vector                      Vector
	VectorCoordsA               Vector
	VectorCoordsB               Vector
	VectorDecodedA              Vector
	VectorDecodedB              Vector
	Covector                    Vector
	CovectorCoordsA             Vector
	CovectorCoordsB             Vector
	CovectorValueA              int
	CovectorValueB              int
	BasisTransitionAToB         Matrix
	BasisTransitionBToA         Matrix
	CoordinateTransitionAToB    Matrix
	CoordinateTransitionBToA    Matrix
	Isotropic                   []Vector
	BasisAValid                 bool
	BasisBValid                 bool
	BasisTransitionsInvert      bool
	CoordinateTransitionsInvert bool
	VectorDecoded               bool
	VectorContravariant         bool
	CovectorCovariant           bool
	CovectorEvaluationInvariant bool
	DualAConjugate              bool
	DualBConjugate              bool
}

func BuildAllDemos() []Demo {
	var demos []Demo
	for _, order := range []int{2, 3, 4} {
		for dim := 1; dim <= 4; dim++ {
			demos = append(demos, BuildDemo(order, dim))
		}
	}
	return demos
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
	decodedA, _ := s.DecodeCoordinates(coordsA, a)
	decodedB, _ := s.DecodeCoordinates(coordsB, b)
	basisAB, _ := s.BasisTransition(a, b)
	basisBA, _ := s.BasisTransition(b, a)
	coordAB, _ := s.CoordinateTransition(a, b)
	coordBA, _ := s.CoordinateTransition(b, a)
	dualA, _ := s.DualBasis(a)
	dualB, _ := s.DualBasis(b)
	covA, _ := s.CovectorCoordinates(cov, dualA)
	covB, _ := s.CovectorCoordinates(cov, dualB)
	covValueA := evalCoordinates(s, covA, coordsA)
	covValueB := evalCoordinates(s, covB, coordsB)
	return Demo{
		Space:                       s,
		BasisA:                      a,
		BasisB:                      b,
		DualA:                       dualA,
		DualB:                       dualB,
		Vector:                      v,
		VectorCoordsA:               coordsA,
		VectorCoordsB:               coordsB,
		VectorDecodedA:              decodedA,
		VectorDecodedB:              decodedB,
		Covector:                    cov,
		CovectorCoordsA:             covA,
		CovectorCoordsB:             covB,
		CovectorValueA:              covValueA,
		CovectorValueB:              covValueB,
		BasisTransitionAToB:         basisAB,
		BasisTransitionBToA:         basisBA,
		CoordinateTransitionAToB:    coordAB,
		CoordinateTransitionBToA:    coordBA,
		Isotropic:                   s.IsotropicVectors(s.IdentityGram()),
		BasisAValid:                 s.IsBasis(a),
		BasisBValid:                 s.IsBasis(b),
		BasisTransitionsInvert:      basisAB.Mul(basisBA).IsIdentity() && basisBA.Mul(basisAB).IsIdentity(),
		CoordinateTransitionsInvert: coordAB.Mul(coordBA).IsIdentity() && coordBA.Mul(coordAB).IsIdentity(),
		VectorDecoded:               EqualVector(decodedA, v) && EqualVector(decodedB, v),
		VectorContravariant:         EqualVector(coordsA, basisAB.MulVec(coordsB)) && EqualVector(coordsB, coordAB.MulVec(coordsA)),
		CovectorCovariant:           EqualVector(covB, basisAB.Transpose().MulVec(covA)) && EqualVector(covA, basisBA.Transpose().MulVec(covB)),
		CovectorEvaluationInvariant: covValueA == covValueB && covValueA == s.Eval(cov, v),
		DualAConjugate:              dualIsConjugate(s, a, dualA),
		DualBConjugate:              dualIsConjugate(s, b, dualB),
	}
}

func evalCoordinates(s Space, covectorCoords, vectorCoords Vector) int {
	sum := s.F.Zero()
	for i := 0; i < s.Dim; i++ {
		sum = s.F.Add(sum, s.F.Mul(covectorCoords[i], vectorCoords[i]))
	}
	return sum
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
