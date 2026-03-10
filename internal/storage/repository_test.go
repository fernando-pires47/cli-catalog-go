package storage

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMissingReturnsEmpty(t *testing.T) {
	tmp := t.TempDir()
	repo := NewRepository(filepath.Join(tmp, "catalog.json"))
	c, err := repo.Load(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Version != "1" || len(c.Commands) != 0 {
		t.Fatalf("unexpected catalog: %+v", c)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "catalog.json")
	if err := os.WriteFile(path, []byte("{"), 0o600); err != nil {
		t.Fatal(err)
	}
	repo := NewRepository(path)
	_, err := repo.Load(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCreateAndDelete(t *testing.T) {
	tmp := t.TempDir()
	repo := NewRepository(filepath.Join(tmp, "catalog.json"))

	created, err := repo.Create(context.Background(), "k", "echo $name")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	catalog, err := repo.Load(context.Background())
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(catalog.Commands) != 1 {
		t.Fatalf("expected 1 command got=%d", len(catalog.Commands))
	}

	if err := repo.DeleteByID(context.Background(), created.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	catalog, err = repo.Load(context.Background())
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(catalog.Commands) != 0 {
		t.Fatalf("expected empty catalog got=%d", len(catalog.Commands))
	}

}
