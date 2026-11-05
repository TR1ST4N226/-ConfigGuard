package config

import (
	"reflect"
	"testing"
)

func TestCompareReportsDrift(t *testing.T) {
	missing, unexpected := Compare(
		[]string{"hostname switch-01", "ssh enabled"},
		[]string{"hostname switch-01", "telnet enabled"},
	)

	if !reflect.DeepEqual(missing, []string{"ssh enabled"}) {
		t.Fatalf("missing = %v, want %v", missing, []string{"ssh enabled"})
	}
	if !reflect.DeepEqual(unexpected, []string{"telnet enabled"}) {
		t.Fatalf("unexpected = %v, want %v", unexpected, []string{"telnet enabled"})
	}
}

func TestCompareIgnoresLineOrder(t *testing.T) {
	missing, unexpected := Compare(
		[]string{"hostname switch-01", "ssh enabled"},
		[]string{"ssh enabled", "hostname switch-01"},
	)

	if len(missing) != 0 || len(unexpected) != 0 {
		t.Fatalf("order-only difference reported: missing=%v unexpected=%v", missing, unexpected)
	}
}
