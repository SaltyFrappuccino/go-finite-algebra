package algebra

import (
	"testing"

	"nct/internal/finite"
)

func TestSingleCounts(t *testing.T) {
	wantSemigroups := map[int]int{2: 8, 3: 113, 4: 3492}
	wantGroups := map[int]int{2: 2, 3: 3, 4: 16}
	wantAbelian := map[int]int{2: 2, 3: 3, 4: 16}
	for n := 2; n <= 4; n++ {
		got := CountSingle(n, finite.Semigroups(n))
		if got.Semigroups != wantSemigroups[n] {
			t.Fatalf("order %d semigroups: got %d want %d", n, got.Semigroups, wantSemigroups[n])
		}
		if got.Groups != wantGroups[n] {
			t.Fatalf("order %d groups: got %d want %d", n, got.Groups, wantGroups[n])
		}
		if got.AbelianGroups != wantAbelian[n] || got.NonAbelianGroups != 0 {
			t.Fatalf("order %d abelian/nonabelian groups: got %d/%d", n, got.AbelianGroups, got.NonAbelianGroups)
		}
	}
}

func TestFieldCountsAndIsomorphism(t *testing.T) {
	wantFields := map[int]int{2: 2, 3: 6, 4: 12}
	wantSemirings := map[int]int{2: 8, 3: 129, 4: 6508}
	wantRings := map[int]int{2: 4, 3: 9, 4: 160}
	wantUnitalRings := map[int]int{2: 2, 3: 6, 4: 72}
	wantDivisionRings := map[int]int{2: 2, 3: 6, 4: 12}
	wantAutomorphisms := map[int]int{2: 1, 3: 1, 4: 2}
	for n := 2; n <= 4; n++ {
		report := BuildOrderReport(n)
		if report.Pairs.Semirings != wantSemirings[n] {
			t.Fatalf("order %d semirings: got %d want %d", n, report.Pairs.Semirings, wantSemirings[n])
		}
		if report.Pairs.Rings != wantRings[n] {
			t.Fatalf("order %d rings: got %d want %d", n, report.Pairs.Rings, wantRings[n])
		}
		if report.Pairs.UnitalRings != wantUnitalRings[n] {
			t.Fatalf("order %d unital rings: got %d want %d", n, report.Pairs.UnitalRings, wantUnitalRings[n])
		}
		if report.Pairs.DivisionRings != wantDivisionRings[n] {
			t.Fatalf("order %d division rings: got %d want %d", n, report.Pairs.DivisionRings, wantDivisionRings[n])
		}
		if report.Pairs.Fields != wantFields[n] {
			t.Fatalf("order %d fields: got %d want %d", n, report.Pairs.Fields, wantFields[n])
		}
		if !report.FieldsPairwiseIsomorphic {
			t.Fatalf("order %d fields are not pairwise isomorphic", n)
		}
		if report.Pairs.DivisionRingsNotFields != 0 {
			t.Fatalf("order %d non-field division rings: got %d", n, report.Pairs.DivisionRingsNotFields)
		}
		if len(report.FieldClasses) != 1 {
			t.Fatalf("order %d field classes: got %d want 1", n, len(report.FieldClasses))
		}
		class := report.FieldClasses[0]
		if class.Size != wantFields[n] {
			t.Fatalf("order %d field class size: got %d want %d", n, class.Size, wantFields[n])
		}
		if class.Automorphisms != wantAutomorphisms[n] {
			t.Fatalf("order %d automorphisms: got %d want %d", n, class.Automorphisms, wantAutomorphisms[n])
		}
		if class.LabelledByFormula != wantFields[n] {
			t.Fatalf("order %d labelled fields by formula: got %d want %d", n, class.LabelledByFormula, wantFields[n])
		}
	}
}

func TestSemiringConventions(t *testing.T) {
	wantWithOne := map[int]int{2: 4, 3: 36, 4: 924}
	wantWithOneNotRings := map[int]int{2: 2, 3: 30, 4: 852}
	for n := 2; n <= 4; n++ {
		p := BuildOrderReport(n).Pairs
		if p.SemiringsWithOne != wantWithOne[n] {
			t.Fatalf("order %d semirings with one: got %d want %d", n, p.SemiringsWithOne, wantWithOne[n])
		}
		if p.SemiringsWithOneNotRings != wantWithOneNotRings[n] {
			t.Fatalf("order %d semirings with one not rings: got %d want %d", n, p.SemiringsWithOneNotRings, wantWithOneNotRings[n])
		}
		if p.SemiringsWithOneNotRings+p.UnitalRings != p.SemiringsWithOne {
			t.Fatalf("order %d semiring-with-one partition broken: %d + %d != %d", n, p.SemiringsWithOneNotRings, p.UnitalRings, p.SemiringsWithOne)
		}
		if p.SemiringsWithOne > p.Semirings {
			t.Fatalf("order %d: semirings-with-one %d exceeds semirings %d", n, p.SemiringsWithOne, p.Semirings)
		}
	}
}

func TestFieldStructuresSatisfyAxioms(t *testing.T) {
	for n := 2; n <= 4; n++ {
		report := BuildOrderReport(n)
		for i, f := range report.Fields {
			if !f.Add.IsAbelianGroup() {
				t.Fatalf("order %d field %d additive table is not an abelian group", n, i)
			}
			if !f.Mul.IsAssociative() || !f.Mul.IsCommutative() {
				t.Fatalf("order %d field %d multiplication is not commutative associative", n, i)
			}
			if !Distributive(f.Add, f.Mul) {
				t.Fatalf("order %d field %d is not distributive", n, i)
			}
			if !IsDivisionRing(f.Add, f.Mul, f.Zero, f.One) {
				t.Fatalf("order %d field %d is not a division ring", n, i)
			}
		}
	}
}
