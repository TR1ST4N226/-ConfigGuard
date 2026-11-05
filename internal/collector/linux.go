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
	commands := []struct {
		name string
		args []string
	}{
		{name: "ip", args: []string{"-brief", "address"}},
		{name: "ip", args: []string{"route"}},
	}

	var lines []string
	for _, command := range commands {
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
