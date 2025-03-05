package display

import (
	"encoding/json"
	"fmt"
	"yaylog/internal/pkgdata"
)

func (o *OutputManager) renderJson(pkgs []pkgdata.PackageInfo, columnNames []string) {
	filteredPackages := make([]pkgdata.PackageInfoJson, len(pkgs))
	for i, pkg := range pkgs {
		filteredPackages[i] = GetColumnJsonValues(pkg, columnNames)
	}

	jsonOutput, err := json.MarshalIndent(filteredPackages, "", "  ")
	if err != nil {
		o.writeLine(fmt.Sprintf("Error genereating JSON output: %v", err))
	}

	o.writeLine(string(jsonOutput))
}
