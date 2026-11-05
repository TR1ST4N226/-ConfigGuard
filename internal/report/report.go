package report

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/TR1ST4N226/ConfigGuard/internal/rules"
)

// Finding est la forme sérialisable d'une règle de conformité en échec.
// Finding is the serializable form of a failed compliance rule.
type Finding struct {
	ID          string         `json:"id"`
	Description string         `json:"description"`
	Severity    rules.Severity `json:"severity"`
}

// Result contient le résultat d'un contrôle de conformité.
// Result contains the outcome of one compliance check.
type Result struct {
	Compliant bool      `json:"compliant"`
	Findings  []Finding `json:"findings"`
}

// FromRules convertit les règles internes en résultat exploitable par un rapport.
// FromRules converts internal rules into a report-safe result.
func FromRules(failed []rules.Rule) Result {
	findings := make([]Finding, 0, len(failed))
	for _, rule := range failed {
		findings = append(findings, Finding{
			ID:          rule.ID,
			Description: rule.Description,
			Severity:    rule.Severity,
		})
	}
	return Result{Compliant: len(findings) == 0, Findings: findings}
}

// WriteJSON écrit un rapport de conformité lisible par une machine.
// WriteJSON writes a machine-readable compliance report.
func WriteJSON(path string, result Result) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("encode JSON report: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write JSON report: %w", err)
	}
	return nil
}

// WriteSARIF écrit un rapport SARIF utilisable par GitHub Code Scanning.
// WriteSARIF writes a SARIF report consumable by GitHub Code Scanning.
func WriteSARIF(path string, result Result) error {
	runs := []sarifRun{{Tool: sarifTool{Driver: sarifDriver{Name: "ConfigGuard", Version: "0.1.0"}}}}
	for _, finding := range result.Findings {
		runs[0].Results = append(runs[0].Results, sarifResult{
			RuleID:  finding.ID,
			Level:   sarifLevel(finding.Severity),
			Message: sarifMessage{Text: finding.Description},
		})
	}

	document := sarifDocument{
		Version: "2.1.0",
		Schema:  "https://json.schemastore.org/sarif-2.1.0.json",
		Runs:    runs,
	}
	data, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		return fmt.Errorf("encode SARIF report: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write SARIF report: %w", err)
	}
	return nil
}

type sarifDocument struct {
	Version string     `json:"version"`
	Schema  string     `json:"$schema"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool    sarifTool     `json:"tool"`
	Results []sarifResult `json:"results,omitempty"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type sarifResult struct {
	RuleID  string       `json:"ruleId"`
	Level   string       `json:"level"`
	Message sarifMessage `json:"message"`
}

type sarifMessage struct {
	Text string `json:"text"`
}

func sarifLevel(severity rules.Severity) string {
	if severity == rules.SeverityHigh {
		return "error"
	}
	return "warning"
}
