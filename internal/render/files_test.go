package render

import (
	"bytes"
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"

	"nct/internal/algebra"
)

func TestWriteCountsCSVForExcel(t *testing.T) {
	path := filepath.Join(t.TempDir(), "counts.csv")
	reports := []algebra.OrderReport{algebra.BuildOrderReport(2)}

	if err := WriteCountsCSV(path, reports); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) < 3 || data[0] != 0xef || data[1] != 0xbb || data[2] != 0xbf {
		t.Fatal("csv has no utf-8 bom")
	}

	reader := csv.NewReader(bytes.NewReader(data[3:]))
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 2 {
		t.Fatalf("rows: got %d want 2", len(records))
	}
	if len(records[0]) != 24 || len(records[1]) != 24 {
		t.Fatalf("columns: got %d/%d want 24/24", len(records[0]), len(records[1]))
	}
	if records[0][0] != "порядок" || records[1][20] != "да" {
		t.Fatalf("unexpected csv values: %q, %q", records[0][0], records[1][20])
	}
}
