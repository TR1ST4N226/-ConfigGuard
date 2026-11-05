package config

import (
	"os"
	"strings"
)

// ReadFile loads a configuration file and removes comments and blank lines.
func ReadFile(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var lines []string
	for _, line := range strings.Split(string(data), "\n") {
		// La normalisation stabilise les comparaisons malgré les différences de formatage.
		// Normalization keeps comparisons stable across formatting-only changes.
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		lines = append(lines, line)
	}
	return lines, nil
}

// Compare returns lines missing from the actual configuration and unexpected lines.
func Compare(desired, actual []string) (missing, unexpected []string) {
	// Les ensembles rendent la recherche indépendante de l'ordre des lignes.
	// Sets make membership checks independent of line order.
	actualSet := make(map[string]struct{}, len(actual))
	for _, line := range actual {
		actualSet[line] = struct{}{}
	}
	for _, line := range desired {
		if _, ok := actualSet[line]; !ok {
			missing = append(missing, line)
		}
	}

	desiredSet := make(map[string]struct{}, len(desired))
	for _, line := range desired {
		desiredSet[line] = struct{}{}
	}
	for _, line := range actual {
		if _, ok := desiredSet[line]; !ok {
			unexpected = append(unexpected, line)
		}
	}
	return missing, unexpected
}
