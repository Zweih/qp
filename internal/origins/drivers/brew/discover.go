package brew

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func getInstalledFormulae(cellarRoot, binRoot string) ([]*installedPkg, error) {
	entries, err := os.ReadDir(cellarRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to read Cellar directory: %w", err)
	}

	var iPkgs []*installedPkg
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		version, err := resolveLinkedVersion(name, cellarRoot, binRoot)
		if err != nil {
			continue
		}

		iPkgs = append(iPkgs, &installedPkg{
			Name:        name,
			Version:     version,
			VersionPath: filepath.Join(cellarRoot, name, version),
			IsTap:       true,
		})
	}

	return iPkgs, nil
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

func getTapPackageNames(prefix string, pkgType string) (map[string]struct{}, error) {
	tapsRoot := filepath.Join(prefix, "Homebrew/Library/Taps")
	result := make(map[string]struct{})

	users, err := os.ReadDir(tapsRoot)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		userPath := filepath.Join(tapsRoot, user.Name())
		repos, err := os.ReadDir(userPath)
		if err != nil {
			continue
		}

		for _, repo := range repos {
			if !repo.IsDir() {
				continue
			}

			repoPath := filepath.Join(userPath, repo.Name())
			var searchDirs []string

			switch pkgType {
			case typeFormula:
				searchDirs = []string{
					filepath.Join(repoPath, "Formula"),
					filepath.Join(repoPath, "HomebrewFormula"),
					repoPath,
				}
			case typeCask:
				searchDirs = []string{filepath.Join(repoPath, "Casks")}
			default:
				continue
			}

			for _, dir := range searchDirs {
				entries, err := os.ReadDir(dir)
				if err != nil {
					continue
				}

				for _, entry := range entries {
					if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".rb") {
						name := strings.TrimSuffix(entry.Name(), ".rb")
						result[name] = struct{}{}
					}
				}
				break
			}
		}
	}

	return result, nil
}
