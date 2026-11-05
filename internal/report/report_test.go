package report

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/TR1ST4N226/ConfigGuard/internal/rules"
)

func TestWriteJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "report.json")
	result := FromRules([]rules.Rule{{
		ID:          "ssh-required",
		Description: "SSH must be enabled",
		Severity:    rules.SeverityHigh,
	}})

	if err := WriteJSON(path, result); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"ssh-required"`) {
		t.Fatalf("JSON report = %s", data)
	}
}

func TestWriteSARIF(t *testing.T) {
	path := filepath.Join(t.TempDir(), "report.sarif")
	result := FromRules([]rules.Rule{{
		ID:          "telnet-disabled",
		Description: "Telnet must be disabled",
		Severity:    rules.SeverityHigh,
	}})

	if err := WriteSARIF(path, result); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"ruleId": "telnet-disabled"`) {
		t.Fatalf("SARIF report = %s", data)
	}
}
