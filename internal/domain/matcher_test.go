package domain

import (
	"testing"
)

func TestResolveBestMatch(t *testing.T) {
	commands := []CatalogCommand{
		{ID: "1", Key: "kill port", Value: "echo $port"},
		{ID: "2", Key: "list pods", Value: "kubectl get pods"},
	}

	cmd, params, args, err := ResolveBestMatch([]string{"kill", "port", "3040"}, commands)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.ID != "1" {
		t.Fatalf("expected id=1 got=%s", cmd.ID)
	}
	if len(params) != 0 {
		t.Fatalf("unexpected params: %v", params)
	}
	if len(args) != 1 || args[0] != "3040" {
		t.Fatalf("unexpected args: %v", args)
	}
}

func TestResolveBestMatchAmbiguous(t *testing.T) {
	commands := []CatalogCommand{
		{ID: "1", Key: "k", Value: "echo 1"},
		{ID: "2", Key: "k", Value: "echo 2"},
	}
	_, _, _, err := ResolveBestMatch([]string{"k"}, commands)
	if err == nil {
		t.Fatal("expected ambiguous error")
	}
}

func TestResolveBestMatchKeyParams(t *testing.T) {
	commands := []CatalogCommand{
		{ID: "1", Key: "kp $port", Value: "echo $port"},
	}

	cmd, params, args, err := ResolveBestMatch([]string{"kp", "3040"}, commands)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.ID != "1" {
		t.Fatalf("expected id=1 got=%s", cmd.ID)
	}
	if got := params["port"]; got != "3040" {
		t.Fatalf("expected port param 3040 got=%q", got)
	}
	if len(args) != 0 {
		t.Fatalf("unexpected args: %v", args)
	}
}

func TestResolveBestMatchKeyParamsMultiple(t *testing.T) {
	commands := []CatalogCommand{
		{ID: "1", Key: "deploy $env $service", Value: "echo $env $service"},
	}

	_, params, args, err := ResolveBestMatch([]string{"deploy", "prod", "api"}, commands)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := params["env"]; got != "prod" {
		t.Fatalf("expected env param prod got=%q", got)
	}
	if got := params["service"]; got != "api" {
		t.Fatalf("expected service param api got=%q", got)
	}
	if len(args) != 0 {
		t.Fatalf("unexpected args: %v", args)
	}
}

func TestResolveBestMatchPrefersLiteralOverPlaceholder(t *testing.T) {
	commands := []CatalogCommand{
		{ID: "1", Key: "kp $target", Value: "echo dynamic"},
		{ID: "2", Key: "kp now", Value: "echo literal"},
	}

	cmd, _, _, err := ResolveBestMatch([]string{"kp", "now"}, commands)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.ID != "2" {
		t.Fatalf("expected literal command id=2 got=%s", cmd.ID)
	}
}

func TestSuggestBestDeterministic(t *testing.T) {
	commands := []CatalogCommand{
		{ID: "3", Key: "kill process", Value: "echo"},
		{ID: "2", Key: "kill pod", Value: "echo"},
		{ID: "1", Key: "kill port", Value: "echo"},
	}

	best, ok := SuggestBest([]string{"kill", "p"}, commands)
	if !ok {
		t.Fatal("expected suggestion")
	}
	if best != "kill pod" {
		t.Fatalf("expected deterministic suggestion 'kill pod', got %q", best)
	}
}
