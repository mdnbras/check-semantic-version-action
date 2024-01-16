package app_test

import (
	"checkVersionApplication/app"
	"testing"
)

func TestCheck(t *testing.T) {
	regexPattern := "(feat|chore|refactor|style|fix|docs|build|perf|ci|revert)([\\(])([\\#0-9]+)([\\)\\: ]+)(\\W|\\w)+"

	tests := []struct {
		name     string
		result   string
		regex    string
		expected bool
		wantErr  bool
	}{
		{"TestMatch", "chore(#124): implementação de alguns testes de integração", regexPattern, true, false},
		{"TestNoMatch", "chore: add discord notification", regexPattern, false, false},
		//{"TestInvalidRegex", "", "[a-z+", false, true}, // Regex inválido
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := app.Check(tt.result, tt.regex)

			if (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.expected {
				t.Errorf("Check() = %v, want %v", got, tt.expected)
			}
		})
	}
}
