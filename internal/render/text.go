package render

import (
	"fmt"
	"strings"

	"nct/internal/algebra"
	"nct/internal/linalg"
)

func PrintAlgebraReport(reports []algebra.OrderReport, outDir string) {
	fmt.Println("I. Конечные алгебраические структуры")
	fmt.Println()
	fmt.Println("Одна бинарная операция")
	fmt.Println("| порядок | таблиц всего | полугрупп | не моноидов | моноидов | не групп | групп | неабелевых групп | абелевых групп |")
	fmt.Println("|---:|---:|---:|---:|---:|---:|---:|---:|---:|")
	for _, r := range reports {
		s := r.Single
		fmt.Printf("| %d | %d | %d | %d | %d | %d | %d | %d | %d |\n",
			s.Order,
			s.BinaryTables,
			s.Semigroups,
			s.SemigroupsNotMonoids,
			s.Monoids,
			s.MonoidsNotGroups,
			s.Groups,
			s.NonAbelianGroups,
			s.AbelianGroups,
		)
	}
	fmt.Println()
	fmt.Println("Пары операций")
	fmt.Println("Полукольца показаны в двух трактовках: широкой (без обязательной 1) и строгой (умножение с двусторонней 1).")
	fmt.Println("| порядок | полуколец | не колец | полуколец с 1 | с 1, не кольца | колец | без 1 | с 1 | с 1, не тел | тел | тел, не полей | полей | поля изоморфны |")
	fmt.Println("|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|:---:|")
	for _, r := range reports {
		p := r.Pairs
		fmt.Printf("| %d | %d | %d | %d | %d | %d | %d | %d | %d | %d | %d | %d | %s |\n",
			p.Order,
			p.Semirings,
			p.SemiringsNotRings,
			p.SemiringsWithOne,
			p.SemiringsWithOneNotRings,
			p.Rings,
			p.RingsWithoutOne,
			p.UnitalRings,
			p.UnitalRingsNotDivisions,
			p.DivisionRings,
			p.DivisionRingsNotFields,
			p.Fields,
			boolText(r.FieldsPairwiseIsomorphic),
		)
	}
	fmt.Println()
	for _, r := range reports {
		fmt.Printf("Поля порядка %d: %d пар таблиц; все таблицы сохранены в %s/fields_order_%d.json\n", r.Single.Order, len(r.Fields), outDir, r.Single.Order)
		if len(r.FieldClasses) > 0 {
			c := r.FieldClasses[0]
			fmt.Printf("Классов изоморфизма: %d; размер класса: %d; |Aut(F)|=%d; n!/|Aut(F)|=%d\n", len(r.FieldClasses), c.Size, c.Automorphisms, c.LabelledByFormula)
		}
		if len(r.Fields) > 0 {
			f := displayField(r.Fields)
			fmt.Printf("Показан представитель с zero=%d, one=%d\n", f.Zero, f.One)
			fmt.Println("Пример сложения:")
			fmt.Println(f.Add.FormatTable())
			fmt.Println("Пример умножения:")
			fmt.Println(f.Mul.FormatTable())
		}
		fmt.Println()
	}
}

func displayField(fields []algebra.FieldStructure) algebra.FieldStructure {
	for _, f := range fields {
		if f.Zero == 0 && f.One == 1 {
			return f
		}
	}
	return fields[0]
}

func PrintVectorDemo(d linalg.Demo) {
	f := d.Space.F
	fmt.Println("II. Векторное пространство и сопряженное пространство")
	fmt.Println()
	fmt.Printf("Поле: %s, dim V = %d\n", f.Name(), d.Space.Dim)
	fmt.Printf("Базис A корректен: %s\n", boolText(d.BasisAValid))
	fmt.Printf("Базис B корректен: %s\n", boolText(d.BasisBValid))
	fmt.Printf("Матрицы перехода взаимно обратны: %s\n", boolText(d.TransitionsInvert))
	fmt.Println()
	fmt.Println("Базис A:")
	printVectors(f, d.BasisA)
	fmt.Println("Базис B:")
	printVectors(f, d.BasisB)
	fmt.Println()
	fmt.Printf("v = %s\n", linalg.FormatVector(f, d.Vector))
	fmt.Printf("[v]_A = %s\n", linalg.FormatVector(f, d.VectorCoordsA))
	fmt.Printf("[v]_B = %s\n", linalg.FormatVector(f, d.VectorCoordsB))
	fmt.Println()
	fmt.Println("T A->B:")
	fmt.Println(d.TransitionAToB.Format())
	fmt.Println("T B->A:")
	fmt.Println(d.TransitionBToA.Format())
	fmt.Printf("Контравариантность координат вектора: %s\n", boolText(d.VectorContravariant))
	fmt.Println()
	fmt.Println("Базис V*, сопряженный к A:")
	printVectors(f, d.DualA)
	fmt.Println("Базис V*, сопряженный к B:")
	printVectors(f, d.DualB)
	fmt.Printf("Сопряженность dual(A): %s\n", boolText(d.DualAConjugate))
	fmt.Printf("Сопряженность dual(B): %s\n", boolText(d.DualBConjugate))
	fmt.Printf("ковектор w = %s\n", linalg.FormatVector(f, d.Covector))
	fmt.Printf("[w]_A* = %s\n", linalg.FormatVector(f, d.CovectorCoordsA))
	fmt.Printf("[w]_B* = %s\n", linalg.FormatVector(f, d.CovectorCoordsB))
	fmt.Printf("Ковариантность координат ковектора: %s\n", boolText(d.CovectorCovariant))
	fmt.Println()
	fmt.Println("Изотропные ненулевые векторы для формы с единичной матрицей Грама:")
	fmt.Println(formatVectorList(f, d.Isotropic, 24))
}

func printVectors(f interface {
	ElementName(int) string
}, vectors []linalg.Vector) {
	for i, v := range vectors {
		fmt.Printf("  %d: %s\n", i, linalg.FormatVector(f, v))
	}
}

func boolText(v bool) string {
	if v {
		return "да"
	}
	return "нет"
}

func formatVectorList(f interface {
	ElementName(int) string
}, vectors []linalg.Vector, limit int) string {
	if len(vectors) == 0 {
		return "  нет"
	}
	var parts []string
	for i, v := range vectors {
		if i >= limit {
			parts = append(parts, fmt.Sprintf("... ещё %d", len(vectors)-limit))
			break
		}
		parts = append(parts, linalg.FormatVector(f, v))
	}
	return "  " + strings.Join(parts, ", ")
}
