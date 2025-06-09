package flatpak

import (
	"fmt"
	"os"
	"path/filepath"
)

type PkgRef struct {
	Name         string
	Type         string
	Arch         string
	Branch       string
	InstallDir   string
	MetadataPath string
}

func discoverPackages(installDirs []string) ([]*PkgRef, error) {
	var allPkgRefs []*PkgRef

	for _, installDir := range installDirs {
		pkgRefs, err := discoverInDir(installDir)
		if err != nil {
			// TODO: log errors
			continue
		}

		allPkgRefs = append(allPkgRefs, pkgRefs...)
	}

	return allPkgRefs, nil
}

// TODO: handle these errors
func discoverInDir(installDir string) ([]*PkgRef, error) {
	var pkgRefs []*PkgRef

	appPkgRefs, err := getPkgRefsOfType(installDir, appsSubdir)
	if err == nil {
		pkgRefs = append(pkgRefs, appPkgRefs...)
	}

	runtimePkgRefs, err := getPkgRefsOfType(installDir, runtimeSubdir)
	if err == nil {
		pkgRefs = append(pkgRefs, runtimePkgRefs...)
	}

	return pkgRefs, nil
}

func getPkgRefsOfType(installDir string, pkgType string) ([]*PkgRef, error) {
	baseDir := filepath.Join(installDir, pkgType)
	if _, err := os.Stat(baseDir); err != nil {
		return nil, fmt.Errorf("failed to find %s directory: %w", baseDir, err)
	}

	var pkgRefs []*PkgRef

	nameEntries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s directory: %w", baseDir, err)
	}

	for _, nameEntry := range nameEntries {
		if !nameEntry.IsDir() {
			continue
		}

		pkgName := nameEntry.Name()
		nameDir := filepath.Join(baseDir, pkgName)

		archEntries, err := os.ReadDir(nameDir)
		if err != nil {
			continue
		}

		for _, archEntry := range archEntries {
			if !archEntry.IsDir() {
				continue
			}

			arch := archEntry.Name()
			archDir := filepath.Join(nameDir, arch)

			branchEntries, err := os.ReadDir(archDir)
			if err != nil {
				continue
			}

			for _, branchEntry := range branchEntries {
				if !branchEntry.IsDir() {
					continue
				}

				branch := branchEntry.Name()
				branchDir := filepath.Join(archDir, branch)

				activePath := filepath.Join(branchDir, activeFile)
				if _, err := os.Lstat(activePath); err != nil {
					continue // skip inactive branches
				}

				commitDir, err := filepath.EvalSymlinks(activePath)
				if err != nil {
					continue
				}

				metadataPath := filepath.Join(commitDir, metadataFile)
				if _, err := os.Stat(metadataPath); err != nil {
					continue // skip if no metadata file
				}

				pkgRefs = append(pkgRefs, &PkgRef{
					Name:         pkgName,
					Arch:         arch,
					Type:         pkgType,
					Branch:       branch,
					InstallDir:   installDir,
					MetadataPath: metadataPath,
				})
			}

		}
	}

	return pkgRefs, nil
}
