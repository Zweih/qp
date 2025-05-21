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

	// var pkgs []*pkgdata.PkgInfo

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		metaPath := filepath.Join(venvRoot, dir.Name(), "lib")
		versionRoot, err := findVersionedPython(metaPath)
		if err != nil {
			return []*pkgdata.PkgInfo{}, fmt.Errorf("can't find versioned root for %s: %v", dir.Name(), err)
		}

		sitePkgsPath := filepath.Join(versionRoot, "site-packages")

		fmt.Println(sitePkgsPath)
	}

	return []*pkgdata.PkgInfo{}, nil
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

	return "", error
}
