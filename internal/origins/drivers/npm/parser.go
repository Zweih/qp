package npm

import (
	"fmt"
	"os"
	"path/filepath"
	"qp/internal/consts"
	"qp/internal/origins/worker"
	"qp/internal/pkgdata"
	"strings"

	json "github.com/goccy/go-json"
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
	pkgDirInfo, err := os.Stat(pkgDir)
	if err != nil {
		return nil, fmt.Errorf("the directory %s does not exist: %w", pkgDir, err)
	}

	packageJsonPath := filepath.Join(pkgDir, packageJsonFile)
	_, err = os.Stat(packageJsonPath)
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

	arch := archAny
	if len(pkgJson.Cpu) > 0 {
		arch = strings.Join(pkgJson.Cpu, ", ")
	}

	pkg := &pkgdata.PkgInfo{
		UpdateTimestamp: pkgDirInfo.ModTime().Unix(),
		Name:            pkgJson.Name,
		Version:         pkgJson.Version,
		Reason:          consts.ReasonExplicit,
		Arch:            arch,
		Description:     pkgJson.Description,
		License:         extractLicense(pkgJson.License),

		Url: pkgJson.Homepage,
	}

	return pkg, nil
}

// not a fan of type assertion, but a custom unmarshaller wouldn't gain much performance here
func extractLicense(license interface{}) string {
	switch value := license.(type) {
	case string:
		return value

	case map[string]interface{}:
		if licenseType, ok := value[fieldType].(string); ok {
			return licenseType
		}

	case nil:
		return ""
	}

	return ""
}

func extractAuthor(author interface{}) string {
	switch value := author.(type) {
	case string:
		return value

	case map[string]interface{}:
		var parts []string
		if name, ok := value[fieldName].(string); ok {
			parts = append(parts, name)
		}

		if email, ok := value[fieldEmail].(string); ok {
			parts = append(parts, "<"+email+">")
		}

		return strings.Join(parts, " ")

	case nil:
		return ""
	}

	return ""
}
