package collector

import (
	"context"
	"fmt"
	"os"
	"strings"
)

// Source identifies how a configuration snapshot is collected.
type Source string

const (
	SourceLocal Source = "local"
	SourceSSH   Source = "ssh"
)

// Platform identifies the target system or network operating system.
type Platform string

const (
	PlatformLinux    Platform = "linux"
	PlatformCiscoIOS Platform = "cisco-ios"
)

// SSHOptions configures a read-only SSH snapshot.
type SSHOptions struct {
	Host       string
	User       string
	KeyPath    string
	KnownHosts string
	Command    string
	Platform   Platform
}

// Write stores a collected configuration as one normalized line per entry.
func Write(path string, lines []string) error {
	content := strings.Join(lines, "\n")
	if len(lines) > 0 {
		content += "\n"
	}
	return os.WriteFile(path, []byte(content), 0600)
}

// ValidateSource checks the supported snapshot sources.
func ValidateSource(source Source) error {
	switch source {
	case SourceLocal, SourceSSH:
		return nil
	default:
		return fmt.Errorf("unsupported snapshot source %q", source)
	}
}

// ValidatePlatform checks the supported target platforms.
func ValidatePlatform(platform Platform) error {
	switch platform {
	case PlatformLinux, PlatformCiscoIOS:
		return nil
	default:
		return fmt.Errorf("unsupported platform %q", platform)
	}
}

// CollectFromSource collects a snapshot using the selected source.
func CollectFromSource(ctx context.Context, source Source, options SSHOptions) ([]string, error) {
	// Centralise le choix de source pour isoler la CLI des détails de collecte.
	// Keep source selection here so the CLI stays independent from collector details.
	switch source {
	case SourceLocal:
		return Collect(ctx)
	case SourceSSH:
		return CollectSSH(ctx, options)
	default:
		return nil, fmt.Errorf("unsupported snapshot source %q", source)
	}
}
