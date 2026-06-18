package linalg

import (
	"nct/internal/field"
)

type Space struct {
	F   field.Field
	Dim int
}

func NewSpace(f field.Field, dim int) Space {
	if dim < 1 || dim > 4 {
		panic("dimension must be between 1 and 4")
	}
	return Space{F: f, Dim: dim}
}

func (s Space) AllVectors() []Vector {
	total := 1
	for i := 0; i < s.Dim; i++ {
		total *= s.F.Order()
	}
	out := make([]Vector, 0, total)
	for code := 0; code < total; code++ {
		x := code
		v := make(Vector, s.Dim)
		for i := 0; i < s.Dim; i++ {
			v[i] = x % s.F.Order()
			x /= s.F.Order()
		}
		out = append(out, v)
	}
	return out
}

func (s Space) StandardBasis() []Vector {
	basis := make([]Vector, s.Dim)
	for i := 0; i < s.Dim; i++ {
		v := make(Vector, s.Dim)
		v[i] = s.F.One()
		basis[i] = v
	}
	return basis
}

func (s Space) TriangularBasis() []Vector {
	basis := s.StandardBasis()
	for i := 0; i+1 < s.Dim; i++ {
		basis[i][i+1] = s.F.One()
	}
	return basis
}

func (s Space) IsBasis(basis []Vector) bool {
	if len(basis) != s.Dim {
		return false
	}
	for _, v := range basis {
		if len(v) != s.Dim {
			return false
		}
	}
	return MatrixFromColumns(s.F, basis).Rank() == s.Dim
}

func (s Space) Coordinates(v Vector, basis []Vector) (Vector, bool) {
	if !s.IsBasis(basis) || len(v) != s.Dim {
		return nil, false
	}
	b := MatrixFromColumns(s.F, basis)
	inv, ok := b.Inverse()
	if !ok {
		return nil, false
	}
	return inv.MulVec(v), true
}

func (s Space) Transition(from, to []Vector) (Matrix, bool) {
	if !s.IsBasis(from) || !s.IsBasis(to) {
		return Matrix{}, false
	}
	out := ZeroMatrix(s.F, s.Dim, s.Dim)
	for j, v := range from {
		coords, ok := s.Coordinates(v, to)
		if !ok {
			return Matrix{}, false
		}
		for i, x := range coords {
			out.Set(i, j, x)
		}
	}
	return out, true
}

func (s Space) DualBasis(basis []Vector) ([]Vector, bool) {
	if !s.IsBasis(basis) {
		return nil, false
	}
	b := MatrixFromColumns(s.F, basis)
	inv, ok := b.Inverse()
	if !ok {
		return nil, false
	}
	dual := make([]Vector, s.Dim)
	for i := 0; i < s.Dim; i++ {
		row := make(Vector, s.Dim)
		for j := 0; j < s.Dim; j++ {
			row[j] = inv.At(i, j)
		}
		dual[i] = row
	}
	return dual, true
}

func (s Space) CovectorCoordinates(cov Vector, dual []Vector) (Vector, bool) {
	if len(cov) != s.Dim || len(dual) != s.Dim {
		return nil, false
	}
	for _, d := range dual {
		if len(d) != s.Dim {
			return nil, false
		}
	}
	return s.Coordinates(cov, dual)
}

func (s Space) Eval(cov Vector, v Vector) int {
	sum := s.F.Zero()
	for i := 0; i < s.Dim; i++ {
		sum = s.F.Add(sum, s.F.Mul(cov[i], v[i]))
	}
	return sum
}

func (s Space) IdentityGram() Matrix {
	return IdentityMatrix(s.F, s.Dim)
}

func (s Space) IsotropicVectors(gram Matrix) []Vector {
	var out []Vector
	for _, v := range s.AllVectors() {
		if zeroVector(s.F, v) {
			continue
		}
		gv := gram.MulVec(v)
		if s.Eval(v, gv) == s.F.Zero() {
			out = append(out, v)
		}
	}
	return out
}

func zeroVector(f field.Field, v Vector) bool {
	for _, x := range v {
		if x != f.Zero() {
			return false
		}
	}
	return true
}

func EqualVector(a, b Vector) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
