package rules

import "testing"

func TestDefaultRulesDetectViolations(t *testing.T) {
	failed := Evaluate([]string{
		"hostname switch-01",
		"telnet enabled",
	}, Default())

	if len(failed) != 3 {
		t.Fatalf("failed rules = %d, want 3", len(failed))
	}
	if failed[0].ID != "ssh-required" || failed[1].ID != "telnet-disabled" || failed[2].ID != "ntp-configured" {
		t.Fatalf("failed rules = %#v", failed)
	}
}

func TestDefaultRulesAcceptSecureConfiguration(t *testing.T) {
	failed := Evaluate([]string{
		"hostname switch-01",
		"ssh enabled",
		"ntp server 10.0.0.10",
	}, Default())

	if len(failed) != 0 {
		t.Fatalf("failed rules = %#v, want none", failed)
	}
}
