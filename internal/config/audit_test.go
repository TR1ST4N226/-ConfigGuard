package config

import "testing"

func TestCompareCompliantConfiguration(t *testing.T) {
	missing, unexpected := Compare(
		[]string{"hostname switch-01", "ssh enabled"},
		[]string{"ssh enabled", "hostname switch-01"},
	)
	if len(missing) != 0 || len(unexpected) != 0 {
		t.Fatalf("configuration reported as drifted: missing=%v unexpected=%v", missing, unexpected)
	}
}
