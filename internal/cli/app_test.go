package cli

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAppCreateListDeleteFlow(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "catalog.json")
	t.Setenv("CS_CATALOG_PATH", path)

	app, err := NewApp()
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	code, _, stderr := runAndCapture(t, app, []string{"create", "kill port", "echo $port"})
	if code != 0 {
		t.Fatalf("create exit=%d", code)
	}
	if stderr != "" {
		t.Fatalf("create stderr should be empty, got=%q", stderr)
	}

	code, stdout, stderr := runAndCapture(t, app, []string{"list"})
	if code != 0 {
		t.Fatalf("list exit=%d", code)
	}
	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	if len(lines) == 0 {
		t.Fatalf("missing list output: %q", stdout)
	}
	headerCols := parsePipedColumns(lines[0])
	if len(headerCols) != 3 || headerCols[0] != "id" || headerCols[1] != "key" || headerCols[2] != "value" {
		t.Fatalf("missing list header: %q", stdout)
	}
	if stderr != "" {
		t.Fatalf("list stderr should be empty, got=%q", stderr)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read catalog: %v", err)
	}
	if len(b) == 0 {
		t.Fatal("catalog should not be empty")
	}
}

func TestCLIErrorContracts(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "catalog.json")
	t.Setenv("CS_CATALOG_PATH", path)

	app, err := NewApp()
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	code, _, stderr := runAndCapture(t, app, []string{"create"})
	if code != 2 || !strings.Contains(stderr, "validation error") {
		t.Fatalf("expected validation contract, code=%d stderr=%q", code, stderr)
	}

	code, _, stderr = runAndCapture(t, app, []string{"unknown", "cmd"})
	if code != 3 || !strings.Contains(stderr, "command not found") {
		t.Fatalf("expected not-found contract, code=%d stderr=%q", code, stderr)
	}

	code, _, _ = runAndCapture(t, app, []string{"create", "k", "echo one"})
	if code != 0 {
		t.Fatalf("create 1 exit=%d", code)
	}
	code, _, _ = runAndCapture(t, app, []string{"create", "k", "echo two"})
	if code != 0 {
		t.Fatalf("create 2 exit=%d", code)
	}

	code, _, stderr = runAndCapture(t, app, []string{"k"})
	if code != 4 || !strings.Contains(stderr, "ambiguous command match") {
		t.Fatalf("expected ambiguous contract, code=%d stderr=%q", code, stderr)
	}
}

func TestCLIDangerDeniedContract(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "catalog.json")
	t.Setenv("CS_CATALOG_PATH", path)

	app, err := NewApp()
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	code, _, _ := runAndCapture(t, app, []string{"create", "wipe", "rm -rf /tmp/safe-test"})
	if code != 0 {
		t.Fatalf("create exit=%d", code)
	}

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	defer r.Close()
	defer w.Close()

	oldStdin := os.Stdin
	os.Stdin = r
	t.Cleanup(func() { os.Stdin = oldStdin })

	code, _, stderr := runAndCapture(t, app, []string{"wipe"})
	if code != 5 || !strings.Contains(stderr, "dangerous command denied") {
		t.Fatalf("expected danger-denied contract, code=%d stderr=%q", code, stderr)
	}
}

func TestDebugLoggingHooks(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "catalog.json")
	t.Setenv("CS_CATALOG_PATH", path)
	t.Setenv("CS_DEBUG", "1")

	app, err := NewApp()
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	code, _, stderr := runAndCapture(t, app, []string{"create", "hello", "printf hi"})
	if code != 0 {
		t.Fatalf("create exit=%d", code)
	}
	if !strings.Contains(stderr, "event=catalog_loaded") || !strings.Contains(stderr, "event=command_created") {
		t.Fatalf("missing create debug hooks: %q", stderr)
	}

	code, stdout, stderr := runAndCapture(t, app, []string{"hello"})
	if code != 0 || strings.TrimSpace(stdout) != "hi" {
		t.Fatalf("execute failed, code=%d stdout=%q", code, stdout)
	}
	if !strings.Contains(stderr, "event=match_resolved") || !strings.Contains(stderr, "event=command_executed") {
		t.Fatalf("missing execute debug hooks: %q", stderr)
	}

	code, _, stderr = runAndCapture(t, app, []string{"list"})
	if code != 0 {
		t.Fatalf("list exit=%d", code)
	}
	if !strings.Contains(stderr, "event=catalog_loaded") {
		t.Fatalf("missing list debug hook: %q", stderr)
	}

	code, stdout, _ = runAndCapture(t, app, []string{"list"})
	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	if len(lines) < 3 {
		t.Fatalf("unexpected list output: %q", stdout)
	}
	rowLine := firstDataRow(lines)
	if rowLine == "" {
		t.Fatalf("could not find list data row: %q", stdout)
	}
	rowCols := parsePipedColumns(rowLine)
	if len(rowCols) == 0 {
		t.Fatalf("unexpected list row: %q", rowLine)
	}
	id := rowCols[0]

	code, _, stderr = runAndCapture(t, app, []string{"delete", id})
	if code != 0 {
		t.Fatalf("delete exit=%d", code)
	}
	if !strings.Contains(stderr, "event=command_deleted") {
		t.Fatalf("missing delete debug hook: %q", stderr)
	}

	code, _, _ = runAndCapture(t, app, []string{"create", "wipe", "rm -rf /tmp/safe-test"})
	if code != 0 {
		t.Fatalf("create dangerous exit=%d", code)
	}

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	defer r.Close()
	defer w.Close()

	oldStdin := os.Stdin
	os.Stdin = r
	t.Cleanup(func() { os.Stdin = oldStdin })

	code, _, stderr = runAndCapture(t, app, []string{"wipe"})
	if code != 5 {
		t.Fatalf("dangerous execute exit=%d", code)
	}
	if !strings.Contains(stderr, "event=danger_confirmation_prompted") {
		t.Fatalf("missing danger debug hook: %q", stderr)
	}
}

func TestDebugFlagEnablesHooks(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "catalog.json")
	t.Setenv("CS_CATALOG_PATH", path)
	t.Setenv("CS_DEBUG", "")

	app, err := NewApp()
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	code, _, stderr := runAndCapture(t, app, []string{"--debug", "create", "hello", "printf hi"})
	if code != 0 {
		t.Fatalf("create exit=%d", code)
	}
	if !strings.Contains(stderr, "event=command_created") {
		t.Fatalf("expected debug event via --debug flag, got=%q", stderr)
	}
}

func TestExecuteBindsKeyParamsIntoValue(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "catalog.json")
	t.Setenv("CS_CATALOG_PATH", path)

	app, err := NewApp()
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	code, _, stderr := runAndCapture(t, app, []string{"create", "kp $port", "echo $port"})
	if code != 0 {
		t.Fatalf("create exit=%d stderr=%q", code, stderr)
	}

	code, stdout, stderr := runAndCapture(t, app, []string{"kp", "3040"})
	if code != 0 {
		t.Fatalf("execute exit=%d stderr=%q", code, stderr)
	}
	if strings.TrimSpace(stdout) != "3040" {
		t.Fatalf("unexpected output: %q", stdout)
	}
}

func TestExecuteBindsMultipleKeyParams(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "catalog.json")
	t.Setenv("CS_CATALOG_PATH", path)

	app, err := NewApp()
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	code, _, stderr := runAndCapture(t, app, []string{"create", "logs $ns $lines", "echo $ns:$lines"})
	if code != 0 {
		t.Fatalf("create exit=%d stderr=%q", code, stderr)
	}

	code, stdout, stderr := runAndCapture(t, app, []string{"logs", "prod", "200"})
	if code != 0 {
		t.Fatalf("execute exit=%d stderr=%q", code, stderr)
	}
	if strings.TrimSpace(stdout) != "prod:200" {
		t.Fatalf("unexpected output: %q", stdout)
	}
}

func runAndCapture(t *testing.T, app *App, args []string) (int, string, string) {
	t.Helper()

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	outR, outW, err := os.Pipe()
	if err != nil {
		t.Fatalf("stdout pipe: %v", err)
	}
	errR, errW, err := os.Pipe()
	if err != nil {
		t.Fatalf("stderr pipe: %v", err)
	}

	os.Stdout = outW
	os.Stderr = errW

	code := app.Run(args)

	_ = outW.Close()
	_ = errW.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	stdoutBytes, _ := io.ReadAll(outR)
	stderrBytes, _ := io.ReadAll(errR)
	_ = outR.Close()
	_ = errR.Close()

	return code, string(stdoutBytes), string(stderrBytes)
}

func parsePipedColumns(line string) []string {
	parts := strings.Split(line, "|")
	cols := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		cols = append(cols, trimmed)
	}
	return cols
}

func firstDataRow(lines []string) string {
	for i, line := range lines {
		if i < 2 {
			continue
		}
		if strings.TrimSpace(line) == "" {
			continue
		}
		return line
	}
	return ""
}
