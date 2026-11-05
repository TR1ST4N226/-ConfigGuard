package inventory

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/TR1ST4N226/ConfigGuard/internal/collector"
)

func TestLoadFileAndFindDevice(t *testing.T) {
	path := filepath.Join(t.TempDir(), "inventory.yaml")
	content := `devices:
  - name: linux-test
    source: ssh
    platform: linux
    ssh:
      host: 192.0.2.10:22
      user: admin
      key: ./.ssh/id_ed25519
      known_hosts: ./.ssh/known_hosts
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	devices, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	device, err := Find(devices, "linux-test")
	if err != nil {
		t.Fatal(err)
	}
	expectedKey := filepath.Join(filepath.Dir(path), ".ssh", "id_ed25519")
	expectedKnownHosts := filepath.Join(filepath.Dir(path), ".ssh", "known_hosts")
	if device.Source != collector.SourceSSH ||
		device.SSH.Host != "192.0.2.10:22" ||
		device.SSH.KeyPath != expectedKey ||
		device.SSH.KnownHosts != expectedKnownHosts {
		t.Fatalf("device = %#v", device)
	}
}

func TestLoadFileRejectsDuplicateNames(t *testing.T) {
	path := filepath.Join(t.TempDir(), "inventory.yaml")
	content := `devices:
  - name: duplicate
    source: local
    platform: linux
  - name: duplicate
    source: local
    platform: linux
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	if _, err := LoadFile(path); err == nil {
		t.Fatal("LoadFile() error = nil, want duplicate name error")
	}
}
