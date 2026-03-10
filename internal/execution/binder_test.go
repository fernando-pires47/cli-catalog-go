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

func TestBindWithNamedSuccess(t *testing.T) {
	out, err := BindWithNamed("echo $port", map[string]string{"port": "3040"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "echo 3040" {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestBindWithNamedAndPositionalSuccess(t *testing.T) {
	out, err := BindWithNamed("echo $env $service $lines", map[string]string{"env": "prod", "service": "api"}, []string{"200"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "echo prod api 200" {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestBindWithNamedRejectsExtraPositionalArgs(t *testing.T) {
	_, err := BindWithNamed("echo $port", map[string]string{"port": "3040"}, []string{"extra"})
	if err == nil {
		t.Fatal("expected error")
	}
}
