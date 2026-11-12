package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/TR1ST4N226/ConfigGuard/internal/collector"
	"github.com/TR1ST4N226/ConfigGuard/internal/config"
	"github.com/TR1ST4N226/ConfigGuard/internal/inventory"
	"github.com/TR1ST4N226/ConfigGuard/internal/report"
	"github.com/TR1ST4N226/ConfigGuard/internal/rules"
)

const version = "0.2.2"

const banner = `░██████                            ░████ ░██             ░██████                                        ░██ 
 ░██   ░██                          ░██                   ░██   ░██                                       ░██ 
░██         ░███████  ░████████  ░████████ ░██ ░████████ ░██        ░██    ░██  ░██████   ░██░████  ░████████ 
░██        ░██    ░██ ░██    ░██    ░██    ░██░██    ░██ ░██  █████ ░██    ░██       ░██  ░███     ░██    ░██ 
░██        ░██    ░██ ░██    ░██    ░██    ░██░██    ░██ ░██     ██ ░██    ░██  ░███████  ░██      ░██    ░██ 
 ░██   ░██ ░██    ░██ ░██    ░██    ░██    ░██░██   ░███  ░██  ░███ ░██   ░███ ░██   ░██  ░██      ░██   ░███ 
  ░██████   ░███████  ░██    ░██    ░██    ░██ ░█████░██   ░█████░█  ░█████░██  ░█████░██ ░██       ░█████░██ 
                                                     ░██                                                      
                                               ░███████`

type auditResult struct {
	drift        bool
	nonCompliant bool
}

func (r auditResult) exitCode() int {
	if r.drift {
		return 4
	}
	if r.nonCompliant {
		return 3
	}
	return 0
}

func main() {
	// Garde le routage léger ; chaque commande possède son propre parsing.
	// Keep command dispatch small; each command owns its argument parsing.
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "version":
		fmt.Println(banner)
		fmt.Printf("configguard %s\n", version)
	case "snapshot":
		if err := runSnapshot(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "configguard: %v\n", err)
			os.Exit(1)
		}
	case "list":
		if err := runList(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "configguard: %v\n", err)
			os.Exit(1)
		}
	case "audit":
		result, err := runAudit(os.Args[2:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "configguard: %v\n", err)
			os.Exit(1)
		}
		if exitCode := result.exitCode(); exitCode != 0 {
			os.Exit(exitCode)
		}
	case "diff":
		drift, err := runDiff(os.Args[2:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "configguard: %v\n", err)
			os.Exit(1)
		}
		if drift {
			os.Exit(4)
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

func runList(args []string) error {
	flags := flag.NewFlagSet("list", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	inventoryPath := flags.String("inventory", "", "path to a YAML inventory")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *inventoryPath == "" {
		return fmt.Errorf("--inventory is required")
	}

	devices, err := inventory.LoadFile(*inventoryPath)
	if err != nil {
		return err
	}
	for _, device := range devices {
		fmt.Printf("%s\t%s\t%s\n", device.Name, device.Source, device.Platform)
	}
	return nil
}

func runAudit(args []string) (auditResult, error) {
	flags := flag.NewFlagSet("audit", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	inventoryPath := flags.String("inventory", "", "path to a YAML inventory")
	deviceName := flags.String("device", "", "device name from the inventory")
	allDevices := flags.Bool("all", false, "audit every device in the inventory")
	desiredPath := flags.String("desired", "", "path to the expected configuration")
	outputPath := flags.String("output", "", "path for the collected snapshot")
	outputDir := flags.String("output-dir", "", "directory for snapshots when using --all")
	timeout := flags.Duration("timeout", 30*time.Second, "maximum collection time")
	rulesPath := flags.String("rules", "", "path to a custom YAML rules file")
	reportFormat := flags.String("report-format", "text", "report format: text, json or sarif")
	reportPath := flags.String("report-output", "", "path for a JSON or SARIF report")
	if err := flags.Parse(args); err != nil {
		return auditResult{}, err
	}
	if *inventoryPath == "" || *desiredPath == "" {
		return auditResult{}, fmt.Errorf("--inventory and --desired are required")
	}
	if (*allDevices && *deviceName != "") || (!*allDevices && *deviceName == "") {
		return auditResult{}, fmt.Errorf("use exactly one of --device or --all")
	}
	if *allDevices {
		if *outputDir == "" || *outputPath != "" {
			return auditResult{}, fmt.Errorf("--all requires --output-dir and cannot use --output")
		}
	} else if *outputPath == "" || *outputDir != "" {
		return auditResult{}, fmt.Errorf("--device requires --output and cannot use --output-dir")
	}
	if *timeout <= 0 {
		return auditResult{}, fmt.Errorf("--timeout must be greater than zero")
	}

	devices, err := inventory.LoadFile(*inventoryPath)
	if err != nil {
		return auditResult{}, err
	}
	desired, err := config.ReadFile(*desiredPath)
	if err != nil {
		return auditResult{}, fmt.Errorf("read desired configuration: %w", err)
	}

	configuredRules := rules.Default()
	if *rulesPath != "" {
		configuredRules, err = rules.LoadFile(*rulesPath)
		if err != nil {
			return auditResult{}, fmt.Errorf("load rules: %w", err)
		}
	}
	if *reportFormat != "text" && *reportPath == "" {
		return auditResult{}, fmt.Errorf("--report-output is required for %s reports", *reportFormat)
	}
	if *allDevices {
		if err := os.MkdirAll(*outputDir, 0700); err != nil {
			return auditResult{}, fmt.Errorf("create output directory: %w", err)
		}
		if *reportFormat != "text" {
			if err := os.MkdirAll(*reportPath, 0700); err != nil {
				return auditResult{}, fmt.Errorf("create report directory: %w", err)
			}
		}
		var aggregate auditResult
		for _, device := range devices {
			result, err := auditDevice(device, desired, filepath.Join(*outputDir, device.Name+".conf"), *timeout, configuredRules, *reportFormat, *reportPath, true)
			if err != nil {
				return auditResult{}, fmt.Errorf("device %q: %w", device.Name, err)
			}
			aggregate.drift = aggregate.drift || result.drift
			aggregate.nonCompliant = aggregate.nonCompliant || result.nonCompliant
		}
		return aggregate, nil
	}
	device, err := inventory.Find(devices, *deviceName)
	if err != nil {
		return auditResult{}, err
	}
	return auditDevice(device, desired, *outputPath, *timeout, configuredRules, *reportFormat, *reportPath, false)
}

func auditDevice(device inventory.Device, desired []string, outputPath string, timeout time.Duration, configuredRules []rules.Rule, reportFormat, reportPath string, multi bool) (auditResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	lines, err := collector.CollectFromSource(ctx, device.Source, device.SSH)
	if err != nil {
		return auditResult{}, fmt.Errorf("collect snapshot: %w", err)
	}
	if err := collector.Write(outputPath, lines); err != nil {
		return auditResult{}, fmt.Errorf("write snapshot: %w", err)
	}
	actual, err := config.ReadFile(outputPath)
	if err != nil {
		return auditResult{}, fmt.Errorf("read snapshot: %w", err)
	}
	missing, unexpected := config.Compare(desired, actual)
	drift := len(missing) > 0 || len(unexpected) > 0
	if drift {
		fmt.Printf("[%s] DRIFT: configuration differs\n", device.Name)
		for _, line := range missing {
			fmt.Printf("[%s] - %s\n", device.Name, line)
		}
		for _, line := range unexpected {
			fmt.Printf("[%s] + %s\n", device.Name, line)
		}
	}
	failed := rules.Evaluate(actual, configuredRules)
	result := report.FromAudit(failed, missing, unexpected)
	if reportFormat == "json" {
		if err := report.WriteJSON(reportPathForDevice(reportPath, device.Name, multi, false), result); err != nil {
			return auditResult{}, err
		}
	} else if reportFormat == "sarif" {
		if err := report.WriteSARIF(reportPathForDevice(reportPath, device.Name, multi, true), result); err != nil {
			return auditResult{}, err
		}
	} else if reportFormat != "text" {
		return auditResult{}, fmt.Errorf("unsupported report format %q", reportFormat)
	}
	if len(failed) > 0 {
		fmt.Printf("[%s] NON-COMPLIANT: %d rule(s) failed\n", device.Name, len(failed))
		return auditResult{drift: drift, nonCompliant: true}, nil
	}
	if !drift {
		fmt.Printf("[%s] COMPLIANT: audit passed\n", device.Name)
	}
	return auditResult{drift: drift}, nil
}

func reportPathForDevice(path, name string, multi, sarif bool) string {
	if !multi {
		return path
	}
	extension := ".json"
	if sarif {
		extension = ".sarif"
	}
	return filepath.Join(path, name+extension)
}

func runCheck(args []string) (bool, error) {
	flags := flag.NewFlagSet("check", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	configPath := flags.String("config", "", "path to the current configuration")
	rulesPath := flags.String("rules", "", "path to a custom YAML rules file")
	reportFormat := flags.String("report-format", "text", "report format: text, json or sarif")
	reportPath := flags.String("report-output", "", "path for a JSON or SARIF report")
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
		// Le fichier personnalisé remplace la base de règles intégrée.
		// A custom rules file replaces the built-in baseline.
		configuredRules, err = rules.LoadFile(*rulesPath)
		if err != nil {
			return false, fmt.Errorf("load rules: %w", err)
		}
	}

	failed := rules.Evaluate(lines, configuredRules)
	if *reportFormat != "text" && *reportPath == "" {
		return false, fmt.Errorf("--report-output is required for %s reports", *reportFormat)
	}
	result := report.FromRules(failed)
	switch *reportFormat {
	case "text":
	case "json":
		if err := report.WriteJSON(*reportPath, result); err != nil {
			return false, err
		}
	case "sarif":
		if err := report.WriteSARIF(*reportPath, result); err != nil {
			return false, err
		}
	default:
		return false, fmt.Errorf("unsupported report format %q", *reportFormat)
	}
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
	timeout := flags.Duration("timeout", 30*time.Second, "maximum collection time")
	inventoryPath := flags.String("inventory", "", "path to a YAML inventory")
	deviceName := flags.String("device", "", "device name from the inventory")
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
	if *timeout <= 0 {
		return fmt.Errorf("--timeout must be greater than zero")
	}

	sourceType := collector.Source(*source)
	targetPlatform := collector.Platform(*platform)
	options := collector.SSHOptions{
		Host:       *sshHost,
		User:       *sshUser,
		KeyPath:    *sshKey,
		KnownHosts: *sshKnownHosts,
		Command:    *sshCommand,
		Platform:   targetPlatform,
	}
	if *inventoryPath != "" || *deviceName != "" {
		if *inventoryPath == "" || *deviceName == "" {
			return fmt.Errorf("--inventory and --device must be used together")
		}
		devices, err := inventory.LoadFile(*inventoryPath)
		if err != nil {
			return err
		}
		device, err := inventory.Find(devices, *deviceName)
		if err != nil {
			return err
		}
		sourceType = device.Source
		targetPlatform = device.Platform
		options = device.SSH
		options.Platform = targetPlatform
	}
	if err := collector.ValidateSource(sourceType); err != nil {
		return err
	}
	if err := collector.ValidatePlatform(targetPlatform); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	// La collecte est en lecture seule : elle récupère l'état et écrit un snapshot local.
	// Collection is read-only: it gathers state and writes a local snapshot.
	lines, err := collector.CollectFromSource(ctx, sourceType, options)
	if err != nil {
		return err
	}
	if err := collector.Write(*outputPath, lines); err != nil {
		return fmt.Errorf("write snapshot: %w", err)
	}
	fmt.Printf("Snapshot saved: %s\n", *outputPath)
	return nil
}

func runDiff(args []string) (bool, error) {
	flags := flag.NewFlagSet("diff", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	desiredPath := flags.String("desired", "", "path to the expected configuration")
	actualPath := flags.String("actual", "", "path to the current configuration")
	if err := flags.Parse(args); err != nil {
		return false, err
	}
	if *desiredPath == "" || *actualPath == "" {
		return false, fmt.Errorf("both --desired and --actual are required")
	}

	desired, err := config.ReadFile(*desiredPath)
	if err != nil {
		return false, fmt.Errorf("read desired configuration: %w", err)
	}
	actual, err := config.ReadFile(*actualPath)
	if err != nil {
		return false, fmt.Errorf("read actual configuration: %w", err)
	}

	// Compare les lignes normalisées pour éviter les faux écarts liés au formatage.
	// Compare normalized lines to avoid formatting-only drift.
	missing, unexpected := config.Compare(desired, actual)
	if len(missing) == 0 && len(unexpected) == 0 {
		fmt.Println("OK: configuration is compliant")
		return false, nil
	}

	fmt.Println("DRIFT: configuration differs")
	for _, line := range missing {
		fmt.Printf("- %s\n", line)
	}
	for _, line := range unexpected {
		fmt.Printf("+ %s\n", line)
	}
	return true, nil
}

func printUsage() {
	fmt.Println(banner)
	fmt.Println(`ConfigGuard detects drift between expected and current network configuration.

Usage:
  configguard version
  configguard list --inventory inventory.yaml
  configguard audit --inventory inventory.yaml --device switch-01 --desired expected.conf --output actual.conf
  configguard audit --inventory inventory.yaml --all --desired expected.conf --output-dir snapshots
  configguard snapshot --output actual.conf
  configguard snapshot --inventory inventory.yaml --device switch-01 --output actual.conf
  configguard snapshot --source ssh --platform cisco-ios --ssh-host 192.0.2.10:22 --ssh-user admin --ssh-key .\id_ed25519 --ssh-known-hosts .\known_hosts --output actual.conf
  configguard diff --desired expected.conf --actual current.conf
  configguard check --config current.conf`)
}
