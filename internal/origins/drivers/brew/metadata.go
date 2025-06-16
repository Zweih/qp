package brew

import (
	"fmt"
	"os"
	"path/filepath"
	"qp/internal/storage"

	json "github.com/goccy/go-json"
)

func getFormulaKey(formula *FormulaMetadata) string {
	return formula.Name
}

func getCaskKey(cask *CaskMetadata) string {
	return cask.Token
}

func loadMetadata[T any](
	filePath string,
	keyFunc func(*T) string,
	wanted map[string]struct{},
) (map[string]*T, error) {
	userCacheDir, err := storage.GetUserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("failed to read brew cache: %w", err)
	}

	fullPath := filepath.Join(userCacheDir, filePath)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s cache: %w", filePath, err)
	}

	var container struct {
		Payload string `json:"payload"`
	}

	if err := json.Unmarshal(data, &container); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	var items []*T
	if err := json.Unmarshal([]byte(container.Payload), &items); err != nil {
		return nil, fmt.Errorf("failed to parse %s payload: %w", filePath, err)
	}

	result := make(map[string]*T, len(wanted))
	for _, item := range items {
		if len(result) >= len(wanted) {
			break
		}

		key := keyFunc(item)
		if _, ok := wanted[key]; ok {
			result[key] = item
		}
	}

	return result, nil
}
