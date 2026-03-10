package execution

import "testing"

func TestBindSuccess(t *testing.T) {
	out, err := Bind("echo $name from $city", []string{"fer", "rio"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "echo fer from rio" {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestBindMissingArgs(t *testing.T) {
	_, err := Bind("echo $name from $city", []string{"fer"})
	if err == nil {
		t.Fatal("expected error")
	}
}
