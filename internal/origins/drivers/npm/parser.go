package npm

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"qp/internal/consts"
	"qp/internal/origins/worker"
	"qp/internal/pkgdata"
	"strings"
)

type PackageJson struct {
	Name        string      `json:"name"`
	Version     string      `json:"version"`
	Description string      `json:"description"`
	Cpu         []string    `json:"cpu"`
	Homepage    string      `json:"homepage"`
	License     interface{} `json:"license"`
}

func parsePackageJson(pkgDir string) (*pkgdata.PkgInfo, error) {
	packageJsonPath := filepath.Join(pkgDir, packageJsonFile)
	jsonInfo, err := os.Stat(packageJsonPath)
	if err != nil {
		return nil, fmt.Errorf("%w. No package.json at %s. It looks like NPM did not remove it after uninstalling. You may want to manually remove it.", worker.ErrSkip, pkgDir)
	}

	data, err := os.ReadFile(packageJsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read package.json in %s: %w", pkgDir, err)
	}

	var pkgJson PackageJson
	if err := json.Unmarshal(data, &pkgJson); err != nil {
		return nil, fmt.Errorf("failed to parse package.json in %s: %w", pkgDir, err)
	}

	arch := "any"
	if len(pkgJson.Cpu) > 0 {
		arch = strings.Join(pkgJson.Cpu, ", ")
	}

	pkg := &pkgdata.PkgInfo{
		UpdateTimestamp: jsonInfo.ModTime().Unix(),
		Name:            pkgJson.Name,
		Version:         pkgJson.Version,
		Reason:          consts.ReasonExplicit,
		Arch:            arch,
		Description:     pkgJson.Description,
		License:         extractLicense(pkgJson.License),
		Url:             pkgJson.Homepage,
	}

	return pkg, nil
}

// not a fan of type assertion, but a custom unmarshaller wouldn't gain much performance here
func extractLicense(license interface{}) string {
	switch v := license.(type) {
	case string:
		return v
	case map[string]interface{}:
		if licType, ok := v["type"].(string); ok {
			return licType
		}
	case nil:
		return ""
	}

	return ""
}
