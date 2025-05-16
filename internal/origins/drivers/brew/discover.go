package brew

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func getInstalledPkgs(cellarRoot, binRoot string) ([]installedPkg, error) {
	entries, err := os.ReadDir(cellarRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to read Cellar directory: %w", err)
	}

	var pkgs []installedPkg
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		version, err := resolveLinkedVersion(name, cellarRoot, binRoot)
		if err != nil {
			continue
		}

		pkgs = append(pkgs, installedPkg{
			Name:        name,
			Version:     version,
			ReceiptPath: filepath.Join(cellarRoot, name, version, receiptName),
			VersionPath: filepath.Join(cellarRoot, name, version),
		})
	}

	return pkgs, nil
}

func resolveLinkedVersion(pkgName string, cellarRoot string, binRoot string) (string, error) {
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
