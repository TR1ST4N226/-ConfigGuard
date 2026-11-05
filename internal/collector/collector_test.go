package collector

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteCreatesRestrictedSnapshot(t *testing.T) {
	path := filepath.Join(t.TempDir(), "snapshot.conf")
	lines := []string{"ip eth0 UP 192.0.2.10/24", "ip default via 192.0.2.1"}

	if err := Write(path, lines); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "ip eth0 UP 192.0.2.10/24\nip default via 192.0.2.1\n" {
		t.Fatalf("snapshot = %q", data)
	}
}
