package algebra

type FieldClass struct {
	Representative    int
	Size              int
	Automorphisms     int
	LabelledByFormula int
}

func FieldsPairwiseIsomorphic(fields []FieldStructure) bool {
	for i := 1; i < len(fields); i++ {
		if !IsIsomorphic(fields[0], fields[i]) {
			return false
		}
	}
	return true
}

func FieldClasses(fields []FieldStructure) []FieldClass {
	used := make([]bool, len(fields))
	var classes []FieldClass
	for i, f := range fields {
		if used[i] {
			continue
		}
		size := 0
		for j, g := range fields {
			if IsIsomorphic(f, g) {
				used[j] = true
				size++
			}
		}
		automorphisms := AutomorphismCount(f)
		classes = append(classes, FieldClass{
			Representative:    i,
			Size:              size,
			Automorphisms:     automorphisms,
			LabelledByFormula: factorial(f.Order) / automorphisms,
		})
	}
	return classes
}

func AutomorphismCount(f FieldStructure) int {
	count := 0
	for _, p := range permutations(f.Order) {
		if preserves(f, f, p) {
			count++
		}
	}
	return count
}

func IsIsomorphic(a, b FieldStructure) bool {
	if a.Order != b.Order {
		return false
	}
	for _, p := range permutations(a.Order) {
		if preserves(a, b, p) {
			return true
		}
	}
	return false
}

func preserves(a, b FieldStructure, p []int) bool {
	n := a.Order
	for x := 0; x < n; x++ {
		for y := 0; y < n; y++ {
			if p[a.Add.At(x, y)] != b.Add.At(p[x], p[y]) {
				return false
			}
			if p[a.Mul.At(x, y)] != b.Mul.At(p[x], p[y]) {
				return false
			}
		}
	}
	return true
}

func permutations(n int) [][]int {
	base := make([]int, n)
	for i := range base {
		base[i] = i
	}
	var out [][]int
	var walk func(int)
	walk = func(pos int) {
		if pos == n {
			p := make([]int, n)
			copy(p, base)
			out = append(out, p)
			return
		}
		for i := pos; i < n; i++ {
			base[pos], base[i] = base[i], base[pos]
			walk(pos + 1)
			base[pos], base[i] = base[i], base[pos]
		}
	}
	walk(0)
	return out
}

func factorial(n int) int {
	x := 1
	for i := 2; i <= n; i++ {
		x *= i
	}
	return x
}
