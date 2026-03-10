package debug

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

var forced bool

func SetForced(enabled bool) {
	forced = enabled
}

func Enabled() bool {
	if forced {
		return true
	}
	value := strings.ToLower(strings.TrimSpace(os.Getenv("CS_DEBUG")))
	return value == "1" || value == "true" || value == "yes"
}

func Event(name string, fields map[string]string) {
	if !Enabled() {
		return
	}

	parts := []string{fmt.Sprintf("event=%s", name)}
	if len(fields) > 0 {
		keys := make([]string, 0, len(fields))
		for key := range fields {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			parts = append(parts, fmt.Sprintf("%s=%q", key, fields[key]))
		}
	}

	fmt.Fprintf(os.Stderr, "debug %s\n", strings.Join(parts, " "))
}
