package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/TR1ST4N226/ConfigGuard/internal/collector"
	"github.com/TR1ST4N226/ConfigGuard/internal/config"
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
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "configguard: unknown command %q\n\n", os.Args[1])
		printUsage()
		os.Exit(2)
	}
}

func runSnapshot(args []string) error {
	flags := flag.NewFlagSet("snapshot", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	outputPath := flags.String("output", "", "path for the collected configuration")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *outputPath == "" {
		return fmt.Errorf("--output is required")
	}

	lines, err := collector.Collect(context.Background())
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
  configguard diff --desired expected.conf --actual current.conf`)
}
