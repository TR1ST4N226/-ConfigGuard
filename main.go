package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/TR1ST4N226/ConfigGuard/internal/collector"
	"github.com/TR1ST4N226/ConfigGuard/internal/config"
	"github.com/TR1ST4N226/ConfigGuard/internal/rules"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "version":
		fmt.Printf("configguard %s\n", version)
	case "snapshot":
		if err := runSnapshot(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "configguard: %v\n", err)
			os.Exit(1)
		}
	case "diff":
		if err := runDiff(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "configguard: %v\n", err)
			os.Exit(1)
		}
	case "check":
		failed, err := runCheck(os.Args[2:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "configguard: %v\n", err)
			os.Exit(1)
		}
		if failed {
			os.Exit(3)
		}
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "configguard: unknown command %q\n\n", os.Args[1])
		printUsage()
		os.Exit(2)
	}
}

func runCheck(args []string) (bool, error) {
	flags := flag.NewFlagSet("check", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	configPath := flags.String("config", "", "path to the current configuration")
	rulesPath := flags.String("rules", "", "path to a custom YAML rules file")
	if err := flags.Parse(args); err != nil {
		return false, err
	}
	if *configPath == "" {
		return false, fmt.Errorf("--config is required")
	}

	lines, err := config.ReadFile(*configPath)
	if err != nil {
		return false, fmt.Errorf("read configuration: %w", err)
	}

	configuredRules := rules.Default()
	if *rulesPath != "" {
		configuredRules, err = rules.LoadFile(*rulesPath)
		if err != nil {
			return false, fmt.Errorf("load rules: %w", err)
		}
	}

	failed := rules.Evaluate(lines, configuredRules)
	if len(failed) == 0 {
		fmt.Println("COMPLIANT: all default rules passed")
		return false, nil
	}

	fmt.Printf("NON-COMPLIANT: %d rule(s) failed\n", len(failed))
	for _, rule := range failed {
		fmt.Printf("- %s [%s] %s\n", rule.ID, rule.Severity, rule.Description)
	}
	return true, nil
}

func runSnapshot(args []string) error {
	flags := flag.NewFlagSet("snapshot", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	source := flags.String("source", string(collector.SourceLocal), "snapshot source: local or ssh")
	outputPath := flags.String("output", "", "path for the collected configuration")
	sshHost := flags.String("ssh-host", "", "SSH host and port, for example 192.0.2.10:22")
	sshUser := flags.String("ssh-user", "", "SSH username")
	sshKey := flags.String("ssh-key", "", "path to an SSH private key")
	sshKnownHosts := flags.String("ssh-known-hosts", "", "path to the SSH known_hosts file")
	platform := flags.String("platform", string(collector.PlatformLinux), "target platform: linux or cisco-ios")
	sshCommand := flags.String("ssh-command", "", "read-only command executed over SSH; defaults to a platform profile")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *outputPath == "" {
		return fmt.Errorf("--output is required")
	}

	sourceType := collector.Source(*source)
	if err := collector.ValidateSource(sourceType); err != nil {
		return err
	}
	targetPlatform := collector.Platform(*platform)
	if err := collector.ValidatePlatform(targetPlatform); err != nil {
		return err
	}
	lines, err := collector.CollectFromSource(context.Background(), sourceType, collector.SSHOptions{
		Host:       *sshHost,
		User:       *sshUser,
		KeyPath:    *sshKey,
		KnownHosts: *sshKnownHosts,
		Command:    *sshCommand,
		Platform:   targetPlatform,
	})
	if err != nil {
		return err
	}
	if err := collector.Write(*outputPath, lines); err != nil {
		return fmt.Errorf("write snapshot: %w", err)
	}
	fmt.Printf("Snapshot saved: %s\n", *outputPath)
	return nil
}

func runDiff(args []string) error {
	flags := flag.NewFlagSet("diff", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	desiredPath := flags.String("desired", "", "path to the expected configuration")
	actualPath := flags.String("actual", "", "path to the current configuration")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *desiredPath == "" || *actualPath == "" {
		return fmt.Errorf("both --desired and --actual are required")
	}

	desired, err := config.ReadFile(*desiredPath)
	if err != nil {
		return fmt.Errorf("read desired configuration: %w", err)
	}
	actual, err := config.ReadFile(*actualPath)
	if err != nil {
		return fmt.Errorf("read actual configuration: %w", err)
	}

	missing, unexpected := config.Compare(desired, actual)
	if len(missing) == 0 && len(unexpected) == 0 {
		fmt.Println("OK: configuration is compliant")
		return nil
	}

	fmt.Println("DRIFT: configuration differs")
	for _, line := range missing {
		fmt.Printf("- %s\n", line)
	}
	for _, line := range unexpected {
		fmt.Printf("+ %s\n", line)
	}
	return nil
}

func printUsage() {
	fmt.Println(`ConfigGuard detects drift between expected and current network configuration.

Usage:
  configguard version
  configguard snapshot --output actual.conf
  configguard snapshot --source ssh --platform cisco-ios --ssh-host 192.0.2.10:22 --ssh-user admin --ssh-key .\id_ed25519 --ssh-known-hosts .\known_hosts --output actual.conf
  configguard diff --desired expected.conf --actual current.conf
  configguard check --config current.conf`)
}
