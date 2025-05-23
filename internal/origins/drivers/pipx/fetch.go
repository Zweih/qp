package pipx

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"qp/internal/consts"
	"qp/internal/origins/shared"
	"qp/internal/pkgdata"
	"strings"
)

func fetchPackages(venvRoot string, origin string) ([]*pkgdata.PkgInfo, error) {
	dirs, err := os.ReadDir(venvRoot)
	if err != nil {
		return []*pkgdata.PkgInfo{}, fmt.Errorf("failed to read pipx venv root: %w", err)
	}

	var pkgs []*pkgdata.PkgInfo

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		dirPath := filepath.Join(venvRoot, dir.Name())
		libPath := filepath.Join(dirPath, "lib")
		versionRoot, err := findVersionPath(libPath)
		if err != nil {
			return []*pkgdata.PkgInfo{}, fmt.Errorf("couldn't locate versioned root for %s: %v", dir.Name(), err)
		}

		sitePkgsPath := filepath.Join(versionRoot, "site-packages")
		metadataPath, err := findDistPath(sitePkgsPath, dir.Name())
		if err != nil {
			return []*pkgdata.PkgInfo{}, fmt.Errorf("couldn't locate metadata file for %s: %v", dir.Name(), err)
		}

		pkg, err := parseMetadataFile(metadataPath)
		if err != nil {
			return []*pkgdata.PkgInfo{}, fmt.Errorf("metadata parsing failed for %s: %v", dir.Name(), err)
		}

		dirInfo, err := dir.Info()
		if err != nil {
			return []*pkgdata.PkgInfo{}, err
		}

		arch, err := findArchitecture(sitePkgsPath)
		if err != nil {
			return []*pkgdata.PkgInfo{}, err
		}

		size, err := shared.GetInstallSize(dirPath)
		if err != nil {
			return []*pkgdata.PkgInfo{}, err
		}

		pkg.Arch = arch
		pkg.Size = size
		pkg.InstallTimestamp = dirInfo.ModTime().Unix()

		pkg.Origin = origin
		pkg.Reason = consts.ReasonExplicit
		pkgs = append(pkgs, pkg)
	}

	return pkgs, nil
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

func findArchitecture(sitePkgsPath string) (string, error) {
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

func parseWheelFile(wheelPath string) (string, error) {
	file, err := os.Open(wheelPath)
	if err != nil {
		return "", err
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		if key == "Tag" {
			return strings.TrimSpace(parts[1]), nil
		}
	}

	return "", fmt.Errorf("no tag field in wheel file for %s", wheelPath)
}
