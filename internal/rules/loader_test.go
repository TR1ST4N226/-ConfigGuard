package rules

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFileCreatesRules(t *testing.T) {
	path := filepath.Join(t.TempDir(), "rules.yaml")
	content := `rules:
  - id: ssh-required
    description: SSH must be enabled
    severity: HIGH
    contains:
      - ssh enabled
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	loaded, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded) != 1 || len(Evaluate([]string{"ssh enabled"}, loaded)) != 0 {
		t.Fatalf("loaded rules did not accept expected configuration: %#v", loaded)
	}
	if len(Evaluate([]string{"telnet enabled"}, loaded)) != 1 {
		t.Fatal("loaded rule did not reject invalid configuration")
	}
}
