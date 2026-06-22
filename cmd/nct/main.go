package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"nct/internal/algebra"
	"nct/internal/linalg"
	"nct/internal/render"
)

func main() {
	out := flag.String("out", "out", "artifact directory")
	fieldOrder := flag.Int("field", 3, "field order: 2, 3, or 4")
	dim := flag.Int("dim", 3, "vector space dimension: 1..4")
	flag.Parse()
	if *fieldOrder != 2 && *fieldOrder != 3 && *fieldOrder != 4 {
		fmt.Fprintln(os.Stderr, "field must be 2, 3, or 4")
		os.Exit(2)
	}
	if *dim < 1 || *dim > 4 {
		fmt.Fprintln(os.Stderr, "dim must be between 1 and 4")
		os.Exit(2)
	}
	var reports []algebra.OrderReport
	for _, n := range []int{2, 3, 4} {
		reports = append(reports, algebra.BuildOrderReport(n))
	}
	if err := render.WriteArtifacts(*out, reports); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	demos := linalg.BuildAllDemos()
	if err := render.WriteVectorSpacesCSV(filepath.Join(*out, "vector_spaces.csv"), demos); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	demo := selectDemo(demos, *fieldOrder, *dim)
	checks := render.BuildVerifications(reports, demos)
	if err := render.WriteVerificationJSON(filepath.Join(*out, "verification.json"), checks); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	render.PrintAlgebraReport(reports, *out)
	render.PrintVectorDemo(demo)
	fmt.Printf("\nПроверки корректности: %s\n", renderStatus(render.AllVerificationsPass(checks)))
	fmt.Printf("Данные: %s/algebra_counts.csv, %s/vector_spaces.csv, %s/verification.json, %s/fields_order_*.json\n", *out, *out, *out, *out)
	fmt.Println("Страница для показа: docs/index.html")
}

func selectDemo(demos []linalg.Demo, order, dim int) linalg.Demo {
	for _, demo := range demos {
		if demo.Space.F.Order() == order && demo.Space.Dim == dim {
			return demo
		}
	}
	return linalg.BuildDemo(order, dim)
}

func renderStatus(ok bool) string {
	if ok {
		return "все пройдены"
	}
	return "есть ошибки"
}
