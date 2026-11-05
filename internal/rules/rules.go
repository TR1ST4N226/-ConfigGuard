package rules

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Severity describes the impact of a failed compliance rule.
type Severity string

const (
	SeverityHigh   Severity = "HIGH"
	SeverityMedium Severity = "MEDIUM"
)

// Rule evaluates one normalized configuration.
type Rule struct {
	ID          string
	Description string
	Severity    Severity
	Check       func([]string) bool
}

type fileRule struct {
	ID          string   `yaml:"id"`
	Description string   `yaml:"description"`
	Severity    Severity `yaml:"severity"`
	Contains    []string `yaml:"contains"`
	NotContains []string `yaml:"not_contains"`
	Prefixes    []string `yaml:"prefixes"`
}

type ruleFile struct {
	Rules []fileRule `yaml:"rules"`
}

// Default returns the initial read-only security baseline.
func Default() []Rule {
	return []Rule{
		{
			ID:          "ssh-required",
			Description: "SSH must be enabled",
			Severity:    SeverityHigh,
			Check: func(lines []string) bool {
				return containsAny(lines, "ssh enabled", "ip ssh", "transport input ssh")
			},
		},
		{
			ID:          "telnet-disabled",
			Description: "Telnet must be disabled",
			Severity:    SeverityHigh,
			Check: func(lines []string) bool {
				return !containsAny(lines, "telnet enabled", "transport input telnet")
			},
		},
		{
			ID:          "ntp-configured",
			Description: "At least one NTP server must be configured",
			Severity:    SeverityMedium,
			Check: func(lines []string) bool {
				return containsPrefix(lines, "ntp server ", "server ") || containsAny(lines, "chrony")
			},
		},
	}
}

// Evaluate returns the rules that failed.
func Evaluate(lines []string, configuredRules []Rule) []Rule {
	var failed []Rule
	for _, rule := range configuredRules {
		if rule.Check == nil || !rule.Check(lines) {
			failed = append(failed, rule)
		}
	}
	return failed
}

// LoadFile loads declarative compliance rules from a YAML file.
func LoadFile(path string) ([]Rule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var document ruleFile
	if err := yaml.Unmarshal(data, &document); err != nil {
		return nil, fmt.Errorf("parse rules file: %w", err)
	}
	if len(document.Rules) == 0 {
		return nil, fmt.Errorf("rules file contains no rules")
	}

	rules := make([]Rule, 0, len(document.Rules))
	for _, definition := range document.Rules {
		if definition.ID == "" || definition.Description == "" {
			return nil, fmt.Errorf("each rule requires id and description")
		}
		if definition.Severity == "" {
			definition.Severity = SeverityMedium
		}
		if len(definition.Contains) == 0 && len(definition.NotContains) == 0 && len(definition.Prefixes) == 0 {
			return nil, fmt.Errorf("rule %q has no condition", definition.ID)
		}

		current := definition
		rules = append(rules, Rule{
			ID:          current.ID,
			Description: current.Description,
			Severity:    current.Severity,
			Check: func(lines []string) bool {
				if len(current.Contains) > 0 && !containsAny(lines, current.Contains...) {
					return false
				}
				if len(current.NotContains) > 0 && containsAny(lines, current.NotContains...) {
					return false
				}
				return len(current.Prefixes) == 0 || containsPrefix(lines, current.Prefixes...)
			},
		})
	}
	return rules, nil
}

func containsAny(lines []string, values ...string) bool {
	for _, line := range lines {
		lowerLine := strings.ToLower(line)
		for _, value := range values {
			if strings.Contains(lowerLine, value) {
				return true
			}
		}
	}
	return false
}

func containsPrefix(lines []string, prefixes ...string) bool {
	for _, line := range lines {
		lowerLine := strings.ToLower(line)
		for _, prefix := range prefixes {
			if strings.HasPrefix(lowerLine, prefix) {
				return true
			}
		}
	}
	return false
}
