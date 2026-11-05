package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
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

	desired, err := readConfig(*desiredPath)
	if err != nil {
		return fmt.Errorf("read desired configuration: %w", err)
	}
	actual, err := readConfig(*actualPath)
	if err != nil {
		return fmt.Errorf("read actual configuration: %w", err)
	}

	missing, unexpected := compare(desired, actual)
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

func readConfig(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var lines []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		lines = append(lines, line)
	}
	return lines, nil
}

func compare(desired, actual []string) (missing, unexpected []string) {
	actualSet := make(map[string]struct{}, len(actual))
	for _, line := range actual {
		actualSet[line] = struct{}{}
	}
	for _, line := range desired {
		if _, ok := actualSet[line]; !ok {
			missing = append(missing, line)
		}
	}

	desiredSet := make(map[string]struct{}, len(desired))
	for _, line := range desired {
		desiredSet[line] = struct{}{}
	}
	for _, line := range actual {
		if _, ok := desiredSet[line]; !ok {
			unexpected = append(unexpected, line)
		}
	}
	return missing, unexpected
}

func printUsage() {
	fmt.Println(`ConfigGuard detects drift between expected and current network configuration.

Usage:
  configguard version
  configguard diff --desired expected.conf --actual current.conf`)
}
