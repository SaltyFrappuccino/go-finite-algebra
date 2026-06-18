package render

import (
	"testing"

	"nct/internal/algebra"
	"nct/internal/linalg"
)

func TestBuildVerifications(t *testing.T) {
	var reports []algebra.OrderReport
	for _, n := range []int{2, 3, 4} {
		reports = append(reports, algebra.BuildOrderReport(n))
	}
	checks := BuildVerifications(reports, linalg.BuildDemo(3, 3))
	if len(checks) != 38 {
		t.Fatalf("checks: got %d want 38", len(checks))
	}
	if !AllVerificationsPass(checks) {
		t.Fatal("expected all verification checks to pass")
	}
}
