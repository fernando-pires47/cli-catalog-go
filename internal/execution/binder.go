package execution

import (
	"fmt"
	"regexp"
	"strings"

	"command-cli/internal/domain"
)

var placeholderRe = regexp.MustCompile(`\$[a-zA-Z_][a-zA-Z0-9_]*`)

func ExtractPlaceholders(template string) []string {
	matches := placeholderRe.FindAllString(template, -1)
	seen := map[string]struct{}{}
	result := make([]string, 0, len(matches))
	for _, m := range matches {
		name := strings.TrimPrefix(m, "$")
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		result = append(result, name)
	}
	return result
}

func Bind(template string, args []string) (string, error) {
	placeholders := ExtractPlaceholders(template)
	if len(placeholders) != len(args) {
		return "", fmt.Errorf("%w: expected %d args for placeholders %v, got %d", domain.ErrValidation, len(placeholders), placeholders, len(args))
	}

	resolved := template
	for i, p := range placeholders {
		resolved = strings.ReplaceAll(resolved, "$"+p, args[i])
	}

	return resolved, nil
}
