package domain

import "testing"

func TestClassifyDanger(t *testing.T) {
	t.Setenv("CS_DANGER_PATTERNS", "")

	dangerous, reasons := ClassifyDanger("rm -rf /tmp/test")
	if !dangerous {
		t.Fatal("expected dangerous")
	}
	if len(reasons) == 0 {
		t.Fatal("expected reasons")
	}

	dangerous, _ = ClassifyDanger("echo hello")
	if dangerous {
		t.Fatal("did not expect dangerous")
	}
}

func TestClassifyDangerWithConfigPatterns(t *testing.T) {
	t.Setenv("CS_DANGER_PATTERNS", "terraform destroy, :(){")

	dangerous, reasons := ClassifyDanger("terraform destroy -auto-approve")
	if !dangerous {
		t.Fatal("expected dangerous")
	}
	if len(reasons) != 1 || reasons[0] != "terraform destroy" {
		t.Fatalf("unexpected reasons: %v", reasons)
	}
}
