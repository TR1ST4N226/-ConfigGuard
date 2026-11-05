//go:build linux

package collector

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// Collect reads Linux network state using read-only system commands.
func Collect(ctx context.Context) ([]string, error) {
	// Ces commandes sont explicites et volontairement limitées à la lecture.
	// These commands are explicit and intentionally read-only.
	commands := []struct {
		name string
		args []string
	}{
		{name: "ip", args: []string{"-brief", "address"}},
		{name: "ip", args: []string{"route"}},
	}

	var lines []string
	for _, command := range commands {
		// CommandContext permettra d'annuler une collecte lente.
		// CommandContext allows a caller to cancel a slow collection.
		output, err := exec.CommandContext(ctx, command.name, command.args...).CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("collect with %s: %w: %s", command.name, err, strings.TrimSpace(string(output)))
		}
		for _, line := range strings.Split(string(output), "\n") {
			line = strings.TrimSpace(line)
			if line != "" {
				lines = append(lines, command.name+" "+line)
			}
		}
	}
	return lines, nil
}
