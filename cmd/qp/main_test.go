package main

import (
	"bytes"
	"os"
	"qp/internal/config"
	"qp/internal/consts"
	"qp/internal/syntax"
	"strings"
	"testing"
)

type MockConfigProvider struct {
	mockConfig config.Config
}

func (m *MockConfigProvider) GetConfig() (*config.Config, error) {
	return &m.mockConfig, nil
}

// TODO: more testing, this is just validating if the depenendency injection works for testing
func TestMainWithConfig(t *testing.T) {
	mockCfg := config.Config{
		Limit:      5,
		SortOption: syntax.SortOption{Field: consts.FieldSize, Asc: false},
		OutputJSON: true,
		Fields:     []consts.FieldType{consts.FieldName, consts.FieldSize},
	}

	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	mainWithConfig(&MockConfigProvider{mockConfig: mockCfg})

	w.Close()
	os.Stdout = oldStdout
	buf.ReadFrom(r)
	output := buf.String()

	if output == "" {
		t.Errorf("Expected output, but got empty string")
	}

	expectedSubstring := "{"
	if mockCfg.OutputJSON && !strings.Contains(output, expectedSubstring) {
		t.Errorf("Expected JSON output but did not find JSON structure")
	}
}
