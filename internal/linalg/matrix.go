package linalg

import (
	"fmt"
	"strings"

	"nct/internal/field"
)

type Vector []int

type Matrix struct {
	F field.Field
	R int
	C int
	A []int
}

func NewMatrix(f field.Field, r, c int, values []int) Matrix {
	a := make([]int, len(values))
	copy(a, values)
	return Matrix{F: f, R: r, C: c, A: a}
}

func ZeroMatrix(f field.Field, r, c int) Matrix {
	return Matrix{F: f, R: r, C: c, A: make([]int, r*c)}
}

func IdentityMatrix(f field.Field, n int) Matrix {
	m := ZeroMatrix(f, n, n)
	for i := 0; i < n; i++ {
		m.Set(i, i, f.One())
	}
	return m
}

func MatrixFromColumns(f field.Field, cols []Vector) Matrix {
	r := len(cols[0])
	c := len(cols)
	m := ZeroMatrix(f, r, c)
	for j, col := range cols {
		for i, v := range col {
			m.Set(i, j, v)
		}
	}
	return m
}

func (m Matrix) At(r, c int) int {
	return m.A[r*m.C+c]
}

func (m Matrix) Set(r, c, v int) {
	m.A[r*m.C+c] = v
}

func (m Matrix) Clone() Matrix {
	return NewMatrix(m.F, m.R, m.C, m.A)
}

func (m Matrix) Mul(n Matrix) Matrix {
	if m.C != n.R {
		panic("matrix dimensions do not match")
	}
	out := ZeroMatrix(m.F, m.R, n.C)
	for i := 0; i < m.R; i++ {
		for j := 0; j < n.C; j++ {
			sum := m.F.Zero()
			for k := 0; k < m.C; k++ {
				sum = m.F.Add(sum, m.F.Mul(m.At(i, k), n.At(k, j)))
			}
			out.Set(i, j, sum)
		}
	}
	return out
}

func (m Matrix) MulVec(v Vector) Vector {
	if m.C != len(v) {
		panic("matrix and vector dimensions do not match")
	}
	out := make(Vector, m.R)
	for i := 0; i < m.R; i++ {
		sum := m.F.Zero()
		for j := 0; j < m.C; j++ {
			sum = m.F.Add(sum, m.F.Mul(m.At(i, j), v[j]))
		}
		out[i] = sum
	}
	return out
}

func (m Matrix) Transpose() Matrix {
	out := ZeroMatrix(m.F, m.C, m.R)
	for i := 0; i < m.R; i++ {
		for j := 0; j < m.C; j++ {
			out.Set(j, i, m.At(i, j))
		}
	}
	return out
}

func (m Matrix) Rank() int {
	a := make([][]int, m.R)
	for i := 0; i < m.R; i++ {
		a[i] = make([]int, m.C)
		for j := 0; j < m.C; j++ {
			a[i][j] = m.At(i, j)
		}
	}
	row := 0
	for col := 0; col < m.C && row < m.R; col++ {
		pivot := -1
		for r := row; r < m.R; r++ {
			if a[r][col] != m.F.Zero() {
				pivot = r
				break
			}
		}
		if pivot < 0 {
			continue
		}
		a[row], a[pivot] = a[pivot], a[row]
		inv, _ := m.F.Inv(a[row][col])
		for c := col; c < m.C; c++ {
			a[row][c] = m.F.Mul(a[row][c], inv)
		}
		for r := 0; r < m.R; r++ {
			if r == row || a[r][col] == m.F.Zero() {
				continue
			}
			factor := a[r][col]
			for c := col; c < m.C; c++ {
				a[r][c] = field.Sub(m.F, a[r][c], m.F.Mul(factor, a[row][c]))
			}
		}
		row++
	}
	return row
}

func (m Matrix) Inverse() (Matrix, bool) {
	if m.R != m.C {
		return Matrix{}, false
	}
	n := m.R
	a := make([][]int, n)
	for i := 0; i < n; i++ {
		a[i] = make([]int, 2*n)
		for j := 0; j < n; j++ {
			a[i][j] = m.At(i, j)
		}
		a[i][n+i] = m.F.One()
	}
	for col := 0; col < n; col++ {
		pivot := -1
		for r := col; r < n; r++ {
			if a[r][col] != m.F.Zero() {
				pivot = r
				break
			}
		}
		if pivot < 0 {
			return Matrix{}, false
		}
		a[col], a[pivot] = a[pivot], a[col]
		inv, _ := m.F.Inv(a[col][col])
		for c := col; c < 2*n; c++ {
			a[col][c] = m.F.Mul(a[col][c], inv)
		}
		for r := 0; r < n; r++ {
			if r == col || a[r][col] == m.F.Zero() {
				continue
			}
			factor := a[r][col]
			for c := col; c < 2*n; c++ {
				a[r][c] = field.Sub(m.F, a[r][c], m.F.Mul(factor, a[col][c]))
			}
		}
	}
	out := ZeroMatrix(m.F, n, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			out.Set(i, j, a[i][n+j])
		}
	}
	return out, true
}

func (m Matrix) IsIdentity() bool {
	if m.R != m.C {
		return false
	}
	for i := 0; i < m.R; i++ {
		for j := 0; j < m.C; j++ {
			want := m.F.Zero()
			if i == j {
				want = m.F.One()
			}
			if m.At(i, j) != want {
				return false
			}
		}
	}
	return true
}

func (m Matrix) Format() string {
	var b strings.Builder
	for i := 0; i < m.R; i++ {
		b.WriteByte('[')
		for j := 0; j < m.C; j++ {
			if j > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(m.F.ElementName(m.At(i, j)))
		}
		b.WriteByte(']')
		if i+1 < m.R {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

type elementNamer interface {
	ElementName(int) string
}

func FormatVector(f elementNamer, v Vector) string {
	parts := make([]string, len(v))
	for i, x := range v {
		parts[i] = f.ElementName(x)
	}
	return fmt.Sprintf("(%s)", strings.Join(parts, ", "))
}
