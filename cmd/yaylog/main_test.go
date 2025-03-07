package main

import (
	"bytes"
	"os"
	"testing"
	"yaylog/internal/config"
)

type MockConfigProvider struct {
	mockConfig config.Config
}

func (m *MockConfigProvider) GetConfig() config.Config {
	return m.mockConfig
}

// TODO: more testing, this is just validating if the depenendency injection works for testing
func TestMainWithConfig(t *testing.T) {
	mockCfg := config.Config{
		Count:       5,
		SortBy:      "size:desc",
		OutputJson:  true,
		ColumnNames: []string{"name", "size"},
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
	if mockCfg.OutputJson && !contains(output, expectedSubstring) {
		t.Errorf("Expected JSON output but did not find JSON structure")
	}
}

func contains(str, substr string) bool {
	return bytes.Contains([]byte(str), []byte(substr))
}
