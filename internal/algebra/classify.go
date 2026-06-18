package algebra

import "nct/internal/finite"

type SingleCounts struct {
	Order                int
	BinaryTables         uint64
	Semigroups           int
	SemigroupsNotMonoids int
	Monoids              int
	MonoidsNotGroups     int
	Groups               int
	NonAbelianGroups     int
	AbelianGroups        int
}

type PairCounts struct {
	Order                    int
	Semirings                int
	SemiringsNotRings        int
	SemiringsWithOne         int
	SemiringsWithOneNotRings int
	Rings                    int
	RingsWithoutOne          int
	UnitalRings              int
	UnitalRingsNotDivisions  int
	DivisionRings            int
	DivisionRingsNotFields   int
	Fields                   int
}

type FieldStructure struct {
	Order int
	Zero  int
	One   int
	Add   finite.Op
	Mul   finite.Op
}

type OrderReport struct {
	Single                   SingleCounts
	Pairs                    PairCounts
	Fields                   []FieldStructure
	FieldClasses             []FieldClass
	FieldsPairwiseIsomorphic bool
}

func BuildOrderReport(n int) OrderReport {
	semigroups := finite.Semigroups(n)
	additiveMonoids := finite.CommutativeMonoids(n)
	single := CountSingle(n, semigroups)
	pairs, fields := CountPairs(n, additiveMonoids, semigroups)
	return OrderReport{
		Single:                   single,
		Pairs:                    pairs,
		Fields:                   fields,
		FieldClasses:             FieldClasses(fields),
		FieldsPairwiseIsomorphic: FieldsPairwiseIsomorphic(fields),
	}
}

func CountSingle(n int, semigroups []finite.Op) SingleCounts {
	counts := SingleCounts{
		Order:        n,
		BinaryTables: finite.PowUint64(n, n*n),
		Semigroups:   len(semigroups),
	}
	for _, op := range semigroups {
		if _, ok := op.Identity(); ok {
			counts.Monoids++
		}
		if op.IsGroup() {
			counts.Groups++
			if op.IsCommutative() {
				counts.AbelianGroups++
			} else {
				counts.NonAbelianGroups++
			}
		}
	}
	counts.SemigroupsNotMonoids = counts.Semigroups - counts.Monoids
	counts.MonoidsNotGroups = counts.Monoids - counts.Groups
	return counts
}

func CountPairs(n int, additiveMonoids []finite.Op, multiplications []finite.Op) (PairCounts, []FieldStructure) {
	counts := PairCounts{Order: n}
	var fields []FieldStructure
	for _, add := range additiveMonoids {
		zero, _ := add.Identity()
		additiveGroup := add.IsGroup()
		for _, mul := range multiplications {
			if !mul.Absorbs(zero) {
				continue
			}
			if !Distributive(add, mul) {
				continue
			}
			counts.Semirings++
			one, hasOne := mul.Identity()
			if hasOne {
				counts.SemiringsWithOne++
			}
			if !additiveGroup {
				counts.SemiringsNotRings++
				if hasOne {
					counts.SemiringsWithOneNotRings++
				}
				continue
			}
			counts.Rings++
			if !hasOne {
				counts.RingsWithoutOne++
				continue
			}
			counts.UnitalRings++
			division := IsDivisionRing(add, mul, zero, one)
			if !division {
				counts.UnitalRingsNotDivisions++
				continue
			}
			counts.DivisionRings++
			if mul.IsCommutative() {
				counts.Fields++
				fields = append(fields, FieldStructure{
					Order: n,
					Zero:  zero,
					One:   one,
					Add:   add,
					Mul:   mul,
				})
			} else {
				counts.DivisionRingsNotFields++
			}
		}
	}
	return counts, fields
}

func Distributive(add, mul finite.Op) bool {
	n := add.N
	for a := 0; a < n; a++ {
		for b := 0; b < n; b++ {
			for c := 0; c < n; c++ {
				if mul.At(a, add.At(b, c)) != add.At(mul.At(a, b), mul.At(a, c)) {
					return false
				}
				if mul.At(add.At(a, b), c) != add.At(mul.At(a, c), mul.At(b, c)) {
					return false
				}
			}
		}
	}
	return true
}

func IsDivisionRing(add, mul finite.Op, zero, one int) bool {
	if zero == one {
		return false
	}
	n := add.N
	for a := 0; a < n; a++ {
		if a == zero {
			continue
		}
		for b := 0; b < n; b++ {
			if b != zero && mul.At(a, b) == zero {
				return false
			}
		}
		found := false
		for b := 0; b < n; b++ {
			if b != zero && mul.At(a, b) == one && mul.At(b, a) == one {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
