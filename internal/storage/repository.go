package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"command-cli/internal/debug"
	"command-cli/internal/domain"
)

type Repository struct {
	Path string
	now  func() time.Time
	id   func() string
}

func NewRepository(path string) *Repository {
	return &Repository{Path: path, now: time.Now, id: defaultID}
}

func (r *Repository) Load(_ context.Context) (domain.CatalogFile, error) {
	b, err := os.ReadFile(r.Path)
	if err != nil {
		if os.IsNotExist(err) {
			debug.Event("catalog_loaded", map[string]string{"path": r.Path, "source": "missing"})
			return domain.EmptyCatalog(), nil
		}
		return domain.CatalogFile{}, err
	}

	var catalog domain.CatalogFile
	if err := json.Unmarshal(b, &catalog); err != nil {
		return domain.CatalogFile{}, fmt.Errorf("%w: fix or remove %s", domain.ErrInvalidCatalog, r.Path)
	}

	if catalog.Version == "" {
		catalog.Version = "1"
	}
	if catalog.Commands == nil {
		catalog.Commands = []domain.CatalogCommand{}
	}

	debug.Event("catalog_loaded", map[string]string{"path": r.Path, "source": "file"})

	return catalog, nil
}

func (r *Repository) Save(_ context.Context, catalog domain.CatalogFile) error {
	if catalog.Version == "" {
		catalog.Version = "1"
	}
	if catalog.Commands == nil {
		catalog.Commands = []domain.CatalogCommand{}
	}

	if err := os.MkdirAll(filepath.Dir(r.Path), 0o755); err != nil {
		return err
	}

	b, err := json.MarshalIndent(catalog, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')

	tmp := r.Path + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}

	if _, err := f.Write(b); err != nil {
		_ = f.Close()
		_ = os.Remove(tmp)
		return err
	}

	if err := f.Sync(); err != nil {
		_ = f.Close()
		_ = os.Remove(tmp)
		return err
	}

	if err := f.Close(); err != nil {
		_ = os.Remove(tmp)
		return err
	}

	if err := os.Rename(tmp, r.Path); err != nil {
		_ = os.Remove(tmp)
		return err
	}

	return nil
}

func (r *Repository) Create(ctx context.Context, key, value string, dangerous bool) (domain.CatalogCommand, error) {
	if err := domain.ValidateCreateInput(key, value); err != nil {
		return domain.CatalogCommand{}, err
	}

	catalog, err := r.Load(ctx)
	if err != nil {
		return domain.CatalogCommand{}, err
	}

	now := r.now().UTC().Format(time.RFC3339)
	entry := domain.CatalogCommand{
		ID:        r.id(),
		Key:       key,
		Value:     value,
		Dangerous: dangerous,
		CreatedAt: now,
		UpdatedAt: now,
	}

	catalog.Commands = append(catalog.Commands, entry)
	if err := r.Save(ctx, catalog); err != nil {
		return domain.CatalogCommand{}, err
	}

	return entry, nil
}

func (r *Repository) DeleteByID(ctx context.Context, id string) error {
	if err := domain.ValidateDeleteInput(id); err != nil {
		return err
	}

	catalog, err := r.Load(ctx)
	if err != nil {
		return err
	}

	idx := -1
	for i, cmd := range catalog.Commands {
		if cmd.ID == id {
			idx = i
			break
		}
	}

	if idx == -1 {
		return fmt.Errorf("%w: id=%s", domain.ErrNotFound, id)
	}

	catalog.Commands = append(catalog.Commands[:idx], catalog.Commands[idx+1:]...)
	return r.Save(ctx, catalog)
}

func defaultID() string {
	return fmt.Sprintf("cmd_%d", time.Now().UnixNano())
}
