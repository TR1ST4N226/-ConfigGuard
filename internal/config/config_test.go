package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestCompareReportsMissingAndUnexpectedLines(t *testing.T) {
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

func TestReadFileIgnoresCommentsAndBlankLines(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.conf")
	content := "# comment\n\n hostname switch-01 \nssh enabled\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	got, err := ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"hostname switch-01", "ssh enabled"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("config = %v, want %v", got, want)
	}
}

func TestCompareReportsCompliance(t *testing.T) {
	missing, unexpected := Compare(
		[]string{"hostname switch-01", "ssh enabled"},
		[]string{"hostname switch-01", "ssh enabled"},
	)

	if len(missing) != 0 || len(unexpected) != 0 {
		t.Fatalf("configuration reported as different: missing=%v unexpected=%v", missing, unexpected)
	}
}
