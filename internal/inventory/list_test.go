package inventory

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFileRejectsInvalidDevice(t *testing.T) {
	path := filepath.Join(t.TempDir(), "inventory.yaml")
	content := `devices:
  - name: invalid
    source: ssh
    platform: linux
    ssh:
      host: 192.0.2.10:22
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	if _, err := LoadFile(path); err == nil {
		t.Fatal("LoadFile() error = nil, want missing SSH fields error")
	}
}
