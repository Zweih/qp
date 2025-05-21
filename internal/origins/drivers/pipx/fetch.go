package pipx

import (
	"fmt"
	"os"
	"path/filepath"
	"qp/internal/pkgdata"
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

		libPath := filepath.Join(venvRoot, dir.Name(), "lib")
		versionRoot, err := findVersionedPython(libPath)
		if err != nil {
			return []*pkgdata.PkgInfo{}, fmt.Errorf("couldn't locate versioned root for %s: %v", dir.Name(), err)
		}

		sitePkgsPath := filepath.Join(versionRoot, "site-packages")
		metadataPath, err := findMetadataFile(sitePkgsPath, dir.Name())
		if err != nil {
			return []*pkgdata.PkgInfo{}, fmt.Errorf("couldn't locate metadata file for %s: %v", dir.Name(), err)
		}

		pkg, err := parseMetadataFile(metadataPath)
		if err != nil {
			return []*pkgdata.PkgInfo{}, fmt.Errorf("metadata parsing failed for %s: %v", dir.Name(), err)
		}

		pkg.Origin = origin
		pkgs = append(pkgs, pkg)
	}

	return pkgs, nil
}

func findVersionedPython(libRoot string) (string, error) {
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

func findMetadataFile(sitePkgsPath string, name string) (string, error) {
	matches, _ := filepath.Glob(filepath.Join(sitePkgsPath, name+"-*.dist-info", "METADATA"))
	if len(matches) == 1 {
		return matches[0], nil
	}

	return "", nil
}
