package render

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"nct/internal/algebra"
)

type FieldJSON struct {
	Order int     `json:"order"`
	Zero  int     `json:"zero"`
	One   int     `json:"one"`
	Add   [][]int `json:"add"`
	Mul   [][]int `json:"mul"`
}

func WriteArtifacts(dir string, reports []algebra.OrderReport) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if err := WriteCountsCSV(filepath.Join(dir, "algebra_counts.csv"), reports); err != nil {
		return err
	}
	for _, r := range reports {
		if err := WriteFieldsJSON(filepath.Join(dir, fmt.Sprintf("fields_order_%d.json", r.Single.Order)), r.Fields); err != nil {
			return err
		}
	}
	return nil
}

func WriteFieldsJSON(path string, fields []algebra.FieldStructure) error {
	values := make([]FieldJSON, len(fields))
	for i, f := range fields {
		values[i] = FieldJSON{
			Order: f.Order,
			Zero:  f.Zero,
			One:   f.One,
			Add:   f.Add.Rows(),
			Mul:   f.Mul.Rows(),
		}
	}
	data, err := json.MarshalIndent(values, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func WriteCountsCSV(path string, reports []algebra.OrderReport) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.WriteString("\ufeff"); err != nil {
		return err
	}
	w := csv.NewWriter(file)
	w.Comma = ';'
	w.UseCRLF = true
	defer w.Flush()
	header := []string{
		"порядок",
		"таблиц всего",
		"полугрупп",
		"полугрупп не моноидов",
		"моноидов",
		"моноидов не групп",
		"групп",
		"неабелевых групп",
		"абелевых групп",
		"полуколец",
		"полуколец не колец",
		"полуколец с 1",
		"полуколец с 1 не колец",
		"колец",
		"колец без 1",
		"колец с 1",
		"колец с 1 не тел",
		"тел",
		"тел не полей",
		"полей",
		"поля изоморфны",
		"классов изоморфизма полей",
		"автоморфизмов представителя",
		"n! делить на автоморфизмы",
	}
	if err := w.Write(header); err != nil {
		return err
	}
	for _, r := range reports {
		row := []string{
			strconv.Itoa(r.Single.Order),
			strconv.FormatUint(r.Single.BinaryTables, 10),
			strconv.Itoa(r.Single.Semigroups),
			strconv.Itoa(r.Single.SemigroupsNotMonoids),
			strconv.Itoa(r.Single.Monoids),
			strconv.Itoa(r.Single.MonoidsNotGroups),
			strconv.Itoa(r.Single.Groups),
			strconv.Itoa(r.Single.NonAbelianGroups),
			strconv.Itoa(r.Single.AbelianGroups),
			strconv.Itoa(r.Pairs.Semirings),
			strconv.Itoa(r.Pairs.SemiringsNotRings),
			strconv.Itoa(r.Pairs.SemiringsWithOne),
			strconv.Itoa(r.Pairs.SemiringsWithOneNotRings),
			strconv.Itoa(r.Pairs.Rings),
			strconv.Itoa(r.Pairs.RingsWithoutOne),
			strconv.Itoa(r.Pairs.UnitalRings),
			strconv.Itoa(r.Pairs.UnitalRingsNotDivisions),
			strconv.Itoa(r.Pairs.DivisionRings),
			strconv.Itoa(r.Pairs.DivisionRingsNotFields),
			strconv.Itoa(r.Pairs.Fields),
			yesNo(r.FieldsPairwiseIsomorphic),
			strconv.Itoa(len(r.FieldClasses)),
			strconv.Itoa(firstClassAutomorphisms(r)),
			strconv.Itoa(firstClassFormula(r)),
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}
	return w.Error()
}

func yesNo(v bool) string {
	if v {
		return "да"
	}
	return "нет"
}

func firstClassAutomorphisms(r algebra.OrderReport) int {
	if len(r.FieldClasses) == 0 {
		return 0
	}
	return r.FieldClasses[0].Automorphisms
}

func firstClassFormula(r algebra.OrderReport) int {
	if len(r.FieldClasses) == 0 {
		return 0
	}
	return r.FieldClasses[0].LabelledByFormula
}
