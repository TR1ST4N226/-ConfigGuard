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

func TestFromAuditIncludesDrift(t *testing.T) {
	result := FromAudit(nil, []string{"ntp server 10.0.0.10"}, []string{"telnet enabled"})

	if result.Compliant {
		t.Fatal("FromAudit() marked drift as compliant")
	}
	if !result.Drift {
		t.Fatal("FromAudit() did not mark drift")
	}
	if len(result.DriftMissing) != 1 || len(result.DriftUnexpected) != 1 {
		t.Fatalf("FromAudit() drift = %#v/%#v", result.DriftMissing, result.DriftUnexpected)
	}
}

func TestWriteSARIFIncludesDrift(t *testing.T) {
	path := filepath.Join(t.TempDir(), "drift.sarif")
	result := FromAudit(nil, []string{"ntp server 10.0.0.10"}, nil)

	if err := WriteSARIF(path, result); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"ruleId": "config-drift"`) ||
		!strings.Contains(string(data), "Missing configuration: ntp server 10.0.0.10") {
		t.Fatalf("SARIF drift report = %s", data)
	}
}
