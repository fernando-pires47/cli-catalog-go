package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"command-cli/internal/debug"
	"command-cli/internal/domain"
	"command-cli/internal/execution"
	"command-cli/internal/storage"
)

type App struct {
	repo *storage.Repository
}

func NewApp() (*App, error) {
	repoPath, err := catalogPath()
	if err != nil {
		return nil, err
	}
	return &App{repo: storage.NewRepository(repoPath)}, nil
}

func (a *App) Run(args []string) int {
	debug.SetForced(false)
	args = normalizeGlobalFlags(args)

	if len(args) == 0 {
		printHelp()
		return 0
	}

	ctx := context.Background()
	switch args[0] {
	case "help", "--help", "-h":
		printHelp()
		return 0
	case "create":
		return a.runCreate(ctx, args[1:])
	case "list":
		return a.runList(ctx)
	case "delete":
		return a.runDelete(ctx, args[1:])
	default:
		return a.runExecute(ctx, args)
	}
}

func (a *App) runCreate(ctx context.Context, args []string) int {
	if len(args) != 2 {
		return printErrf(2, "%w: usage: cs create \"<key>\" \"<value>\"", domain.ErrValidation)
	}

	cmd, err := a.repo.Create(ctx, args[0], args[1])
	if err != nil {
		return printDomainError(err)
	}
	debug.Event("command_created", map[string]string{"id": cmd.ID, "key": cmd.Key})

	fmt.Printf("created %s\n", cmd.ID)
	return 0
}

func (a *App) runList(ctx context.Context) int {
	catalog, err := a.repo.Load(ctx)
	if err != nil {
		return printDomainError(err)
	}

	if len(catalog.Commands) == 0 {
		fmt.Println("no commands found")
		return 0
	}

	sort.Slice(catalog.Commands, func(i, j int) bool {
		if catalog.Commands[i].Key != catalog.Commands[j].Key {
			return catalog.Commands[i].Key < catalog.Commands[j].Key
		}
		return catalog.Commands[i].ID < catalog.Commands[j].ID
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "id\t|\tkey\t|\tvalue")
	fmt.Fprintln(w, "--\t|\t---\t|\t-----")
	for _, cmd := range catalog.Commands {
		fmt.Fprintf(w, "%s\t|\t%s\t|\t%s\n", cmd.ID, cmd.Key, cmd.Value)
	}
	_ = w.Flush()
	return 0
}

func (a *App) runDelete(ctx context.Context, args []string) int {
	if len(args) != 1 {
		return printErrf(2, "%w: usage: cs delete <id>", domain.ErrValidation)
	}

	if err := a.repo.DeleteByID(ctx, args[0]); err != nil {
		return printDomainError(err)
	}
	debug.Event("command_deleted", map[string]string{"id": args[0]})

	fmt.Printf("deleted %s\n", args[0])
	return 0
}

func (a *App) runExecute(ctx context.Context, input []string) int {
	catalog, err := a.repo.Load(ctx)
	if err != nil {
		return printDomainError(err)
	}

	cmd, keyParams, runtimeArgs, err := domain.ResolveBestMatch(input, catalog.Commands)
	if err != nil {
		return printDomainError(err)
	}
	debug.Event("match_resolved", map[string]string{"command_id": cmd.ID, "key": cmd.Key})

	resolved, err := execution.BindWithNamed(cmd.Value, keyParams, runtimeArgs)
	if err != nil {
		return printDomainError(err)
	}

	isDangerous, reasons := domain.ClassifyDanger(resolved)
	if isDangerous {
		prompt := fmt.Sprintf("dangerous command detected (%s). continue?", strings.Join(reasons, ", "))
		debug.Event("danger_confirmation_prompted", map[string]string{"command_id": cmd.ID, "reasons": strings.Join(reasons, ",")})
		if _, err := execution.ConfirmDanger(prompt); err != nil {
			return printDomainError(err)
		}
	}

	result, err := execution.Run(ctx, resolved)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}

	if result.Stdout != "" {
		fmt.Print(result.Stdout)
	}
	if result.Stderr != "" {
		fmt.Fprint(os.Stderr, result.Stderr)
	}

	debug.Event("command_executed", map[string]string{"command_id": cmd.ID, "exit_code": fmt.Sprintf("%d", result.ExitCode)})

	return result.ExitCode
}

func catalogPath() (string, error) {
	if fromEnv := strings.TrimSpace(os.Getenv("CS_CATALOG_PATH")); fromEnv != "" {
		return fromEnv, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cs", "catalog.json"), nil
}

func printHelp() {
	fmt.Println("cs - local command catalog")
	fmt.Println("usage:")
	fmt.Println("  cs [--debug] <command>")
	fmt.Println("  cs create \"<key>\" \"<value>\"")
	fmt.Println("  cs list")
	fmt.Println("  cs delete <id>")
	fmt.Println("  cs <key...> [args...]")
}

func normalizeGlobalFlags(args []string) []string {
	filtered := make([]string, 0, len(args))
	for _, arg := range args {
		if arg == "--debug" {
			debug.SetForced(true)
			continue
		}
		filtered = append(filtered, arg)
	}
	return filtered
}

func printErrf(code int, format string, a ...any) int {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	return code
}

func printDomainError(err error) int {
	code := 1
	switch {
	case errors.Is(err, domain.ErrValidation):
		code = 2
	case errors.Is(err, domain.ErrNotFound):
		code = 3
	case errors.Is(err, domain.ErrAmbiguous):
		code = 4
	case errors.Is(err, domain.ErrDangerDenied):
		code = 5
	}
	fmt.Fprintln(os.Stderr, err.Error())
	return code
}
