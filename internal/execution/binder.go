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
	return BindWithNamed(template, map[string]string{}, args)
}

func BindWithNamed(template string, named map[string]string, args []string) (string, error) {
	placeholders := ExtractPlaceholders(template)
	missing := make([]string, 0)
	positionalIndex := 0
	values := map[string]string{}

	for _, p := range placeholders {
		if v, ok := named[p]; ok {
			values[p] = v
			continue
		}
		if positionalIndex >= len(args) {
			missing = append(missing, p)
			continue
		}
		values[p] = args[positionalIndex]
		positionalIndex++
	}

	if len(missing) > 0 {
		return "", fmt.Errorf("%w: missing args for placeholders %v", domain.ErrValidation, missing)
	}

	if positionalIndex != len(args) {
		return "", fmt.Errorf("%w: expected %d positional args for placeholders %v, got %d", domain.ErrValidation, positionalIndex, placeholders, len(args))
	}

	resolved := template
	for _, p := range placeholders {
		resolved = strings.ReplaceAll(resolved, "$"+p, values[p])
	}

	return resolved, nil
}
