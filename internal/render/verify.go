package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"nct/internal/algebra"
	"nct/internal/linalg"
)

type Verification struct {
	Name   string `json:"name"`
	Status bool   `json:"status"`
	Detail string `json:"detail"`
}

func BuildVerifications(reports []algebra.OrderReport, demos []linalg.Demo) []Verification {
	var out []Verification
	for _, r := range reports {
		out = append(out,
			Verification{
				Name:   verificationName(r.Single.Order, "полугруппы"),
				Status: r.Single.Semigroups >= r.Single.Monoids && r.Single.Monoids >= r.Single.Groups,
				Detail: "полугруппы >= моноиды >= группы",
			},
			Verification{
				Name:   verificationName(r.Single.Order, "разбиение полугрупп"),
				Status: r.Single.SemigroupsNotMonoids+r.Single.Monoids == r.Single.Semigroups,
				Detail: "не моноиды + моноиды = полугруппы",
			},
			Verification{
				Name:   verificationName(r.Single.Order, "разбиение моноидов"),
				Status: r.Single.MonoidsNotGroups+r.Single.Groups == r.Single.Monoids,
				Detail: "не группы + группы = моноиды",
			},
			Verification{
				Name:   verificationName(r.Single.Order, "разбиение групп"),
				Status: r.Single.NonAbelianGroups+r.Single.AbelianGroups == r.Single.Groups,
				Detail: "неабелевы + абелевы = группы",
			},
			Verification{
				Name:   verificationName(r.Single.Order, "разбиение полуколец"),
				Status: r.Pairs.SemiringsNotRings+r.Pairs.Rings == r.Pairs.Semirings,
				Detail: "не кольца + кольца = полукольца",
			},
			Verification{
				Name:   verificationName(r.Single.Order, "разбиение полуколец с 1"),
				Status: r.Pairs.SemiringsWithOneNotRings+r.Pairs.UnitalRings == r.Pairs.SemiringsWithOne,
				Detail: "полукольца с 1, не кольца + кольца с 1 = полукольца с 1",
			},
			Verification{
				Name:   verificationName(r.Single.Order, "разбиение колец"),
				Status: r.Pairs.RingsWithoutOne+r.Pairs.UnitalRings == r.Pairs.Rings,
				Detail: "без 1 + с 1 = кольца",
			},
			Verification{
				Name:   verificationName(r.Single.Order, "разбиение колец с единицей"),
				Status: r.Pairs.UnitalRingsNotDivisions+r.Pairs.DivisionRings == r.Pairs.UnitalRings,
				Detail: "не тела + тела = кольца с 1",
			},
			Verification{
				Name:   verificationName(r.Single.Order, "разбиение тел"),
				Status: r.Pairs.DivisionRingsNotFields+r.Pairs.Fields == r.Pairs.DivisionRings,
				Detail: "тела не поля + поля = тела",
			},
			Verification{
				Name:   verificationName(r.Single.Order, "изоморфизм полей"),
				Status: r.FieldsPairwiseIsomorphic && len(r.FieldClasses) == 1,
				Detail: "все найденные поля одного порядка лежат в одном классе",
			},
			Verification{
				Name:   verificationName(r.Single.Order, "размеченные поля"),
				Status: len(r.FieldClasses) == 1 && r.FieldClasses[0].LabelledByFormula == r.Pairs.Fields,
				Detail: "число полей совпадает с n!/|Aut(F)|",
			},
		)
	}
	for _, demo := range demos {
		prefix := fmt.Sprintf("%s dim %d: ", demo.Space.F.Name(), demo.Space.Dim)
		out = append(out,
			Verification{
				Name:   prefix + "базисы V",
				Status: demo.BasisAValid && demo.BasisBValid,
				Detail: "оба выбранных базиса имеют полный ранг",
			},
			Verification{
				Name:   prefix + "кодирование и декодирование",
				Status: demo.VectorDecoded,
				Detail: "координаты в A и B собираются обратно в один и тот же вектор",
			},
			Verification{
				Name:   prefix + "матрицы перехода базисов",
				Status: demo.BasisTransitionsInvert,
				Detail: "C(A<-B) и C(B<-A) взаимно обратны",
			},
			Verification{
				Name:   prefix + "матрицы пересчета координат",
				Status: demo.CoordinateTransitionsInvert,
				Detail: "матрицы пересчета A->B и B->A взаимно обратны",
			},
			Verification{
				Name:   prefix + "координаты вектора",
				Status: demo.VectorContravariant,
				Detail: "[v]_A = C[v]_B и [v]_B = C^-1[v]_A",
			},
			Verification{
				Name:   prefix + "сопряженные базисы",
				Status: demo.DualAConjugate && demo.DualBConjugate,
				Detail: "dual(A) и dual(B) удовлетворяют условию сопряженности",
			},
			Verification{
				Name:   prefix + "координаты ковектора",
				Status: demo.CovectorCovariant,
				Detail: "[phi]_B = C^T[phi]_A",
			},
			Verification{
				Name:   prefix + "значение ковектора",
				Status: demo.CovectorEvaluationInvariant,
				Detail: "phi(v) не зависит от выбранного базиса",
			},
		)
	}
	return out
}

func WriteVerificationJSON(path string, checks []Verification) error {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(checks); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0644)
}

func AllVerificationsPass(checks []Verification) bool {
	for _, check := range checks {
		if !check.Status {
			return false
		}
	}
	return true
}

func verificationName(order int, name string) string {
	return "порядок " + strconv.Itoa(order) + ": " + name
}
