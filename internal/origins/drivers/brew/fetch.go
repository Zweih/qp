package brew

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"qp/internal/origins/worker"
	"qp/internal/pkgdata"
	"strings"
	"sync"

	json "github.com/goccy/go-json"
)

type installedPkg struct {
	Name        string
	Version     string
	ReceiptPath string
	VersionPath string
}

func fetchPackages(
	origin string,
	prefix string,
) ([]*pkgdata.PkgInfo, error) {
	binRoot := filepath.Join(prefix, binSubPath)
	cellarRoot := filepath.Join(prefix, cellarSubPath)
	installedPkgs, err := getInstalledPkgs(cellarRoot, binRoot)
	if err != nil {
		return nil, err
	}

	wanted := make(map[string]struct{}, len(installedPkgs))
	for _, iPkg := range installedPkgs {
		wanted[iPkg.Name] = struct{}{}
	}

	var formulaMeta map[string]*FormulaMetadata
	var metaErr error
	var metaWg sync.WaitGroup

	metaWg.Add(1)
	go func() {
		defer metaWg.Done()
		formulaMeta, metaErr = loadFormulaMetadataSubset(wanted)
	}()

	inputChan := make(chan installedPkg, len(installedPkgs))
	for _, iPkg := range installedPkgs {
		inputChan <- iPkg
	}

	close(inputChan)

	stage1Out, stage1Err := worker.RunWorkers(
		inputChan,
		func(iPkg installedPkg) (*pkgdata.PkgInfo, error) {
			return parseInstallReceipt(iPkg.ReceiptPath, iPkg.Version)
		},
		0,
		len(installedPkgs),
	)

	stage2Out, stage2Err := worker.RunWorkers(
		stage1Out,
		func(pkg *pkgdata.PkgInfo) (*pkgdata.PkgInfo, error) {
			versionPath := filepath.Join(prefix, cellarSubPath, pkg.Name, pkg.Version)
			if size, err := getInstallSize(versionPath); err == nil {
				pkg.Size = size
			}

			return pkg, nil
		},
		0,
		len(installedPkgs),
	)

	metaWg.Wait()
	if metaErr != nil {
		return nil, metaErr
	}

	stage3Out, stage3Err := worker.RunWorkers(
		stage2Out,
		func(pkg *pkgdata.PkgInfo) (*pkgdata.PkgInfo, error) {
			if meta, ok := formulaMeta[pkg.Name]; ok {
				mergeFormulaMetadata(pkg, meta)
			}
			pkg.Origin = origin
			return pkg, nil
		},
		0,
		len(installedPkgs),
	)

	allErrs := worker.MergeErrors(stage1Err, stage2Err, stage3Err)
	return worker.CollectOutput(stage3Out, allErrs)
}

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

func loadFormulaMetadataSubset(wanted map[string]struct{}) (map[string]*FormulaMetadata, error) {
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

	var formulas []*FormulaMetadata
	if err := json.Unmarshal([]byte(container.Payload), &formulas); err != nil {
		return nil, fmt.Errorf("failed to parse formula payload: %w", err)
	}

	result := make(map[string]*FormulaMetadata, len(wanted))
	for _, f := range formulas {
		if _, ok := wanted[f.Name]; ok {
			result[f.Name] = f
		}
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
