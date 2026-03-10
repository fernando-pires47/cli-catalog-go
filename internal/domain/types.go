package domain

type CatalogFile struct {
	Version  string           `json:"version"`
	Commands []CatalogCommand `json:"commands"`
}

type CatalogCommand struct {
	ID        string `json:"id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Dangerous bool   `json:"dangerous"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

type ExecutionPlan struct {
	CommandID            string
	ResolvedCommand      string
	IsDangerous          bool
	RequiresConfirmation bool
}

type ExecutionResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
}

func EmptyCatalog() CatalogFile {
	return CatalogFile{Version: "1", Commands: []CatalogCommand{}}
}
