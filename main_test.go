package main

import "testing"

func TestAuditResultExitCode(t *testing.T) {
	tests := []struct {
		name         string
		result       auditResult
		expectedCode int
	}{
		{
			name:         "compliant",
			result:       auditResult{},
			expectedCode: 0,
		},
		{
			name:         "non-compliant",
			result:       auditResult{nonCompliant: true},
			expectedCode: 3,
		},
		{
			name:         "drift",
			result:       auditResult{drift: true},
			expectedCode: 4,
		},
		{
			name:         "drift takes precedence",
			result:       auditResult{drift: true, nonCompliant: true},
			expectedCode: 4,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if code := test.result.exitCode(); code != test.expectedCode {
				t.Fatalf("exitCode() = %d, want %d", code, test.expectedCode)
			}
		})
	}
}
