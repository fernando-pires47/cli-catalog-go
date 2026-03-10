package domain

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type Match struct {
	Command      CatalogCommand
	KeyLen       int
	LiteralCount int
	KeyParams    map[string]string
}

var keyPlaceholderRe = regexp.MustCompile(`^\$[a-zA-Z_][a-zA-Z0-9_]*$`)

func ResolveBestMatch(inputTokens []string, commands []CatalogCommand) (CatalogCommand, map[string]string, []string, error) {
	best := make([]Match, 0)
	bestLiteralCount := -1
	bestLen := -1

	for _, cmd := range commands {
		keyTokens := strings.Fields(cmd.Key)
		if len(keyTokens) == 0 || len(inputTokens) < len(keyTokens) {
			continue
		}

		captured := map[string]string{}
		literalCount := 0
		ok := true
		for i := range keyTokens {
			if name, isPlaceholder := keyPlaceholderName(keyTokens[i]); isPlaceholder {
				captured[name] = inputTokens[i]
				continue
			}
			literalCount++
			if inputTokens[i] != keyTokens[i] {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}

		if literalCount > bestLiteralCount || (literalCount == bestLiteralCount && len(keyTokens) > bestLen) {
			bestLiteralCount = literalCount
			bestLen = len(keyTokens)
			best = []Match{{Command: cmd, KeyLen: len(keyTokens), LiteralCount: literalCount, KeyParams: captured}}
			continue
		}

		if literalCount == bestLiteralCount && len(keyTokens) == bestLen {
			best = append(best, Match{Command: cmd, KeyLen: len(keyTokens), LiteralCount: literalCount, KeyParams: captured})
		}
	}

	if len(best) == 0 {
		return CatalogCommand{}, nil, nil, fmt.Errorf("%w: %s", ErrNotFound, strings.Join(inputTokens, " "))
	}

	if len(best) > 1 {
		sort.Slice(best, func(i, j int) bool {
			if best[i].Command.Key != best[j].Command.Key {
				return best[i].Command.Key < best[j].Command.Key
			}
			return best[i].Command.ID < best[j].Command.ID
		})
		candidates := make([]string, 0, len(best))
		for _, m := range best {
			candidates = append(candidates, m.Command.Key)
		}
		return CatalogCommand{}, nil, nil, fmt.Errorf("%w: candidates=%s", ErrAmbiguous, strings.Join(candidates, ", "))
	}

	matched := best[0]
	return matched.Command, matched.KeyParams, inputTokens[matched.KeyLen:], nil
}

func keyPlaceholderName(token string) (string, bool) {
	if !keyPlaceholderRe.MatchString(token) {
		return "", false
	}
	return strings.TrimPrefix(token, "$"), true
}

func Suggest(prefixTokens []string, commands []CatalogCommand) []string {
	prefix := strings.Join(prefixTokens, " ")
	results := make([]string, 0)
	for _, cmd := range commands {
		if prefix == "" || strings.HasPrefix(cmd.Key, prefix) {
			results = append(results, cmd.Key)
		}
	}
	sort.Strings(results)
	return results
}

func SuggestBest(prefixTokens []string, commands []CatalogCommand) (string, bool) {
	suggestions := Suggest(prefixTokens, commands)
	if len(suggestions) == 0 {
		return "", false
	}
	return suggestions[0], true
}
