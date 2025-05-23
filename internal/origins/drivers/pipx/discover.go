package pipx

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func findVenvRoot() (string, error) {
	var venvRootPath string

	if custom := os.Getenv("PIPX_HOME"); custom != "" {
		venvRootPath = filepath.Join(custom, "venvs")
		_, err := os.Stat(venvRootPath)
		if err == nil {
			return venvRootPath, nil
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		home = os.Getenv("HOME")
	}

	possibleRoots := []string{defaultVenvRoot, otherVenvRoot}
	for _, root := range possibleRoots {
		venvRootPath = filepath.Join(home, root)
		_, err := os.Stat(venvRootPath)
		if err == nil {
			return venvRootPath, nil
		}
	}

	return "", errors.New("no pipx venv root found")
}

func findVersionPath(libRoot string) (string, error) {
	entries, err := os.ReadDir(libRoot)
	if err != nil {
		return "", err
	}

	for _, e := range entries {
		if e.IsDir() && e.Name()[:6] == "python" {
			return filepath.Join(libRoot, e.Name()), nil
		}
	}
	return "", fmt.Errorf("no pythonX.Y found under %s", libRoot)
}

func findDistPath(sitePkgsPath string, name string) (string, error) {
	matches, _ := filepath.Glob(filepath.Join(sitePkgsPath, name+"-*.dist-info", "METADATA"))
	if len(matches) > 0 {
		return matches[0], nil
	}

	return "", nil
}

func inferArchitecture(sitePkgsPath string) (string, error) {
	matches, err := filepath.Glob(filepath.Join(sitePkgsPath, "*.dist-info", "WHEEL"))
	if err != nil {
		return "", fmt.Errorf("found no .dist-info directories in %s: %v", sitePkgsPath, err)
	}

	bestMatch := "any"

	for _, match := range matches {
		arch, err := parseWheelFile(match)
		if err != nil {
			continue
		}

		parts := strings.Split(arch, "-")
		suffix := parts[len(parts)-1]

		if suffix != "any" {
			if !strings.Contains(suffix, "universal") {
				return suffix, nil
			}

			bestMatch = suffix
		}
	}

	return bestMatch, nil
}
