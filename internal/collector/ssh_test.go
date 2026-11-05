package collector

import (
	"context"
	"strings"
	"testing"
)

func TestCollectSSHValidatesRequiredOptions(t *testing.T) {
	_, err := CollectSSH(context.Background(), SSHOptions{})
	if err == nil {
		t.Fatal("CollectSSH() error = nil, want missing options error")
	}
	if !strings.Contains(err.Error(), "ssh-host") {
		t.Fatalf("CollectSSH() error = %q, want ssh-host guidance", err)
	}
}

func TestValidateSource(t *testing.T) {
	if err := ValidateSource(SourceSSH); err != nil {
		t.Fatalf("ValidateSource(ssh) error = %v", err)
	}
	if err := ValidateSource(Source("unknown")); err == nil {
		t.Fatal("ValidateSource(unknown) error = nil, want error")
	}
}

func TestValidatePlatform(t *testing.T) {
	if err := ValidatePlatform(PlatformCiscoIOS); err != nil {
		t.Fatalf("ValidatePlatform(cisco-ios) error = %v", err)
	}
	if err := ValidatePlatform(Platform("unknown")); err == nil {
		t.Fatal("ValidatePlatform(unknown) error = nil, want error")
	}
}

func TestDefaultCommand(t *testing.T) {
	if command := defaultCommand(PlatformLinux); command != "ip -brief address; ip route" {
		t.Fatalf("Linux command = %q", command)
	}
	if command := defaultCommand(PlatformCiscoIOS); !strings.Contains(command, "show ip interface brief") {
		t.Fatalf("Cisco IOS command = %q", command)
	}
}
