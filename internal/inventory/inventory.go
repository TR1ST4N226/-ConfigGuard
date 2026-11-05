package inventory

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/TR1ST4N226/ConfigGuard/internal/collector"
	"gopkg.in/yaml.v3"
)

// Device describes one read-only snapshot target.
type Device struct {
	Name     string               `yaml:"name"`
	Source   collector.Source     `yaml:"source"`
	Platform collector.Platform   `yaml:"platform"`
	SSH      collector.SSHOptions `yaml:"ssh"`
}

type file struct {
	Devices []Device `yaml:"devices"`
}

// LoadFile reads and validates a YAML inventory.
func LoadFile(path string) ([]Device, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var document file
	if err := yaml.Unmarshal(data, &document); err != nil {
		return nil, fmt.Errorf("parse inventory: %w", err)
	}
	if len(document.Devices) == 0 {
		return nil, fmt.Errorf("inventory contains no devices")
	}

	seen := make(map[string]struct{}, len(document.Devices))
	for index := range document.Devices {
		device := &document.Devices[index]
		if device.Name == "" {
			return nil, fmt.Errorf("device %d has no name", index+1)
		}
		if _, exists := seen[device.Name]; exists {
			return nil, fmt.Errorf("duplicate device name %q", device.Name)
		}
		seen[device.Name] = struct{}{}
		if err := collector.ValidateSource(device.Source); err != nil {
			return nil, fmt.Errorf("device %q: %w", device.Name, err)
		}
		if err := collector.ValidatePlatform(device.Platform); err != nil {
			return nil, fmt.Errorf("device %q: %w", device.Name, err)
		}
		if device.Source == collector.SourceSSH {
			if device.SSH.Host == "" || device.SSH.User == "" ||
				device.SSH.KeyPath == "" || device.SSH.KnownHosts == "" {
				return nil, fmt.Errorf("device %q: ssh host, user, key and known_hosts are required", device.Name)
			}
			device.SSH.KeyPath = resolvePath(path, device.SSH.KeyPath)
			device.SSH.KnownHosts = resolvePath(path, device.SSH.KnownHosts)
		}
		device.SSH.Platform = device.Platform
	}
	return document.Devices, nil
}

// Find returns one device by its unique inventory name.
func Find(devices []Device, name string) (Device, error) {
	for _, device := range devices {
		if device.Name == name {
			return device, nil
		}
	}
	return Device{}, fmt.Errorf("device %q not found in inventory", name)
}

func resolvePath(inventoryPath, configuredPath string) string {
	if filepath.IsAbs(configuredPath) {
		return filepath.Clean(configuredPath)
	}
	return filepath.Clean(filepath.Join(filepath.Dir(inventoryPath), configuredPath))
}
