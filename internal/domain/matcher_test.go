package domain

import (
	"testing"
)

func TestResolveBestMatch(t *testing.T) {
	commands := []CatalogCommand{
		{ID: "1", Key: "kill port", Value: "echo $port"},
		{ID: "2", Key: "list pods", Value: "kubectl get pods"},
	}

	cmd, args, err := ResolveBestMatch([]string{"kill", "port", "3040"}, commands)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.ID != "1" {
		t.Fatalf("expected id=1 got=%s", cmd.ID)
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
	_, _, err := ResolveBestMatch([]string{"k"}, commands)
	if err == nil {
		t.Fatal("expected ambiguous error")
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
