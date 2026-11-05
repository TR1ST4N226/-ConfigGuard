package collector

import (
	"os"
	"strings"
)

// Write stores a collected configuration as one normalized line per entry.
func Write(path string, lines []string) error {
	content := strings.Join(lines, "\n")
	if len(lines) > 0 {
		content += "\n"
	}
	return os.WriteFile(path, []byte(content), 0600)
}
