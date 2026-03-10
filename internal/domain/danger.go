package domain

import (
	"os"
	"strings"
)

var defaultDangerousPatterns = []string{
	"rm -rf",
	"mkfs",
	"dd if=",
	"shutdown",
	"reboot",
	":(){",
}

func ClassifyDanger(command string) (bool, []string) {
	return classifyDangerWithPatterns(command, configuredDangerPatterns())
}

func classifyDangerWithPatterns(command string, patterns []string) (bool, []string) {
	cmd := strings.ToLower(command)
	reasons := make([]string, 0)
	for _, pattern := range patterns {
		if strings.Contains(cmd, pattern) {
			reasons = append(reasons, pattern)
		}
	}
	return len(reasons) > 0, reasons
}

func configuredDangerPatterns() []string {
	combined := append([]string{}, defaultDangerousPatterns...)
	raw := strings.TrimSpace(os.Getenv("CS_DANGER_PATTERNS"))
	if raw == "" {
		return combined
	}

	for _, token := range strings.Split(raw, ",") {
		pattern := strings.ToLower(strings.TrimSpace(token))
		if pattern == "" {
			continue
		}
		combined = append(combined, pattern)
	}

	return uniqueStrings(combined)
}

func uniqueStrings(values []string) []string {
	seen := map[string]struct{}{}
	result := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}
