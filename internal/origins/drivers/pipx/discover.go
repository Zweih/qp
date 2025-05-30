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

	if custom := os.Getenv(pipxHomeEnv); custom != "" {
		venvRootPath = filepath.Join(custom, "venvs")
		_, err := os.Stat(venvRootPath)
		if err == nil {
			return venvRootPath, nil
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		home = os.Getenv(homeEnv)
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

	for _, entry := range entries {
		if entry.IsDir() && len(entry.Name()) > 6 && entry.Name()[:6] == pythonEntry {
			return filepath.Join(libRoot, entry.Name()), nil
		}
	}

	return "", fmt.Errorf("no pythonX.Y found under %s", libRoot)
}

func findDistPath(sitePkgsPath string, name string) (string, error) {
	distInfoMatcher := name + "-*" + dotDistInfo
	matches, _ := filepath.Glob(filepath.Join(sitePkgsPath, distInfoMatcher, metadataFile))
	if len(matches) > 0 {
		return matches[0], nil
	}

	return "", fmt.Errorf("no .dist-info/METADATA found for: %s", name)
}

func inferArchitecture(sitePkgsPath string) (string, error) {
	distInfoMatcher := "*" + dotDistInfo
	matches, err := filepath.Glob(filepath.Join(sitePkgsPath, distInfoMatcher, wheelFile))
	if err != nil {
		return "", fmt.Errorf("found no .dist-info directories in %s: %v", sitePkgsPath, err)
	}

	bestMatch := anyArch

	for _, match := range matches {
		arch, err := parseWheelFile(match)
		if err != nil {
			continue
		}

		parts := strings.Split(arch, "-")
		suffix := parts[len(parts)-1]

		if suffix != anyArch {
			if !strings.Contains(suffix, universalArch) {
				return suffix, nil
			}

			bestMatch = suffix
		}
	}

	return bestMatch, nil
}
