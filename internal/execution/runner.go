package execution

import (
	"bytes"
	"context"
	"os/exec"

	"command-cli/internal/domain"
)

func Run(ctx context.Context, command string) (domain.ExecutionResult, error) {
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	res := domain.ExecutionResult{Stdout: stdout.String(), Stderr: stderr.String()}
	if err == nil {
		res.ExitCode = 0
		return res, nil
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		res.ExitCode = exitErr.ExitCode()
		return res, nil
	}

	res.ExitCode = 1
	return res, err
}
