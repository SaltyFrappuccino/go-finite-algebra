package finite

func Semigroups(n int) []Op {
	var out []Op
	EnumerateSemigroups(n, func(op Op) bool {
		out = append(out, op)
		return true
	})
	return out
}

func EnumerateSemigroups(n int, yield func(Op) bool) {
	t := make([]int, n*n)
	for i := range t {
		t[i] = -1
	}
	var walk func([]int) bool
	walk = func(current []int) bool {
		idx, values, complete, ok := chooseCell(n, current)
		if !ok {
			return true
		}
		if complete {
			return yield(NewOp(n, current))
		}
		for _, v := range values {
			next := cloneTable(current)
			next[idx] = v
			if propagateAssociativity(n, next) {
				if !walk(next) {
					return false
				}
			}
		}
		return true
	}
	walk(t)
}

func CommutativeMonoids(n int) []Op {
	var out []Op
	EnumerateCommutativeMonoids(n, func(op Op) bool {
		out = append(out, op)
		return true
	})
	return out
}

func EnumerateCommutativeMonoids(n int, yield func(Op) bool) {
	var slots [][2]int
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			slots = append(slots, [2]int{i, j})
		}
	}
	total := PowUint64(n, len(slots))
	for code := uint64(0); code < total; code++ {
		x := code
		t := make([]int, n*n)
		for _, slot := range slots {
			v := int(x % uint64(n))
			x /= uint64(n)
			i, j := slot[0], slot[1]
			t[i*n+j] = v
			t[j*n+i] = v
		}
		op := NewOp(n, t)
		if !op.IsAssociative() {
			continue
		}
		if _, ok := op.Identity(); !ok {
			continue
		}
		if !yield(op) {
			return
		}
	}
}

func chooseCell(n int, t []int) (int, []int, bool, bool) {
	best := -1
	var bestValues []int
	for idx, value := range t {
		if value >= 0 {
			continue
		}
		values := make([]int, 0, n)
		for v := 0; v < n; v++ {
			next := cloneTable(t)
			next[idx] = v
			if propagateAssociativity(n, next) {
				values = append(values, v)
			}
		}
		if len(values) == 0 {
			return idx, nil, false, false
		}
		if best < 0 || len(values) < len(bestValues) {
			best = idx
			bestValues = values
			if len(values) == 1 {
				break
			}
		}
	}
	if best < 0 {
		return -1, nil, true, true
	}
	return best, bestValues, false, true
}

func propagateAssociativity(n int, t []int) bool {
	changed := true
	for changed {
		changed = false
		for a := 0; a < n; a++ {
			for b := 0; b < n; b++ {
				ab := t[a*n+b]
				if ab < 0 {
					continue
				}
				for c := 0; c < n; c++ {
					bc := t[b*n+c]
					if bc < 0 {
						continue
					}
					leftCell := ab*n + c
					rightCell := a*n + bc
					left := t[leftCell]
					right := t[rightCell]
					switch {
					case left >= 0 && right >= 0:
						if left != right {
							return false
						}
					case left >= 0:
						t[rightCell] = left
						changed = true
					case right >= 0:
						t[leftCell] = right
						changed = true
					}
				}
			}
		}
	}
	return true
}

func cloneTable(t []int) []int {
	out := make([]int, len(t))
	copy(out, t)
	return out
}
