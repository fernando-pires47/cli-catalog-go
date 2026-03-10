package execution

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"command-cli/internal/domain"
)

func IsInteractive() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}

func ConfirmDanger(prompt string) (bool, error) {
	if !IsInteractive() {
		return false, fmt.Errorf("%w: non-interactive mode blocks dangerous command", domain.ErrDangerDenied)
	}

	fmt.Fprint(os.Stderr, prompt+" [y/N]: ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	resp := strings.ToLower(strings.TrimSpace(line))
	if resp == "y" || resp == "yes" {
		return true, nil
	}

	return false, fmt.Errorf("%w: user rejected confirmation", domain.ErrDangerDenied)
}
