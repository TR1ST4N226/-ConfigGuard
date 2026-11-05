package rules

import "strings"

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
