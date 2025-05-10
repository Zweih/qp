package brew

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"qp/internal/pkgdata"
	"strings"
)

func fetchPackages(origin string, prefix string) ([]*pkgdata.PkgInfo, error) {
	formulaMeta, err := loadFormulaMetadata()
	if err != nil {
		return nil, err
	}

	binRoot := filepath.Join(prefix, binSubPath)
	cellarRoot := filepath.Join(prefix, cellarSubPath)
	installedNames, err := getInstalledPkgNames(cellarRoot)
	if err != nil {
		return nil, err
	}

	var pkgs []*pkgdata.PkgInfo
	for _, name := range installedNames {
		version, err := getLinkedVersion(name, cellarRoot, binRoot)
		if err != nil {
			return nil, err
		}

		receiptPath := filepath.Join(cellarRoot, name, version, "INSTALL_RECEIPT.json")
		pkg, err := parseInstallReceipt(receiptPath)
		if err != nil {
			return nil, err
		}

		if meta, ok := formulaMeta[name]; ok {
			mergeFormulaMetadata(pkg, meta)
		}

		versionPath := filepath.Join(prefix, cellarSubPath, pkg.Name, pkg.Version)
		size, err := getInstallSize(versionPath)
		if err == nil {
			pkg.Size = size
		}

		pkg.Origin = origin
		pkgs = append(pkgs, pkg)
	}

	return pkgs, nil
}

func getInstalledPkgNames(cellarRoot string) ([]string, error) {
	entries, err := os.ReadDir(cellarRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to read Cellar directory: %w", err)
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}

	return names, nil
}

func getLinkedVersion(pkgName string, cellarRoot string, binRoot string) (string, error) {
	pkgPath := filepath.Join(cellarRoot, pkgName)

	entries, err := os.ReadDir(pkgPath)
	if err != nil {
		return "", fmt.Errorf("failed to read Cellar/%s: %w", pkgName, err)
	}

	// only check symlinks if there are multiple versions
	if len(entries) == 1 {
		return entries[0].Name(), nil
	}

	binPath := filepath.Join(binRoot, pkgName)
	target, err := os.Readlink(binPath)
	if err != nil {
		return "", fmt.Errorf("no symlink found in /bin for %s", pkgName)
	}

	absPath, err := filepath.Abs(filepath.Join(filepath.Dir(binPath), target))
	if err != nil {
		return "", err
	}

	parts := strings.Split(filepath.Clean(absPath), string(os.PathSeparator))
	if len(parts) < 3 {
		return "", fmt.Errorf("unexpected symlink path: %s", absPath)
	}

	return parts[len(parts)-3], nil
}

func loadFormulaMetadata() (map[string]*FormulaMetadata, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	fullPath := filepath.Join(homeDir, formulaCachePath)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read formula cache: %w", err)
	}

	var container struct {
		Payload string `json:"payload"`
	}

	if err := json.Unmarshal(data, &container); err != nil {
		return nil, fmt.Errorf("failed to parse formula jws: %w", err)
	}

	var formulae []*FormulaMetadata
	if err := json.Unmarshal([]byte(container.Payload), &formulae); err != nil {
		return nil, fmt.Errorf("failed to parse formula payload: %w", err)
	}

	result := make(map[string]*FormulaMetadata)
	for _, formula := range formulae {
		result[formula.Name] = formula
	}

	return result, nil
}

func getInstallSize(dir string) (int64, error) {
	var total int64

	err := filepath.WalkDir(dir, func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			total += info.Size()
		}

		return nil
	})

	return total, err
}
