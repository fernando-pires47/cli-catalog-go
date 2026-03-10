package domain

import (
	"fmt"
	"sort"
	"strings"
)

type Match struct {
	Command CatalogCommand
	KeyLen  int
}

func ResolveBestMatch(inputTokens []string, commands []CatalogCommand) (CatalogCommand, []string, error) {
	best := make([]Match, 0)
	bestLen := -1

	for _, cmd := range commands {
		keyTokens := strings.Fields(cmd.Key)
		if len(keyTokens) == 0 || len(inputTokens) < len(keyTokens) {
			continue
		}

		ok := true
		for i := range keyTokens {
			if inputTokens[i] != keyTokens[i] {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}

		if len(keyTokens) > bestLen {
			bestLen = len(keyTokens)
			best = []Match{{Command: cmd, KeyLen: len(keyTokens)}}
			continue
		}

		if len(keyTokens) == bestLen {
			best = append(best, Match{Command: cmd, KeyLen: len(keyTokens)})
		}
	}

	if len(best) == 0 {
		return CatalogCommand{}, nil, fmt.Errorf("%w: %s", ErrNotFound, strings.Join(inputTokens, " "))
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
		return CatalogCommand{}, nil, fmt.Errorf("%w: candidates=%s", ErrAmbiguous, strings.Join(candidates, ", "))
	}

	matched := best[0]
	return matched.Command, inputTokens[matched.KeyLen:], nil
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
