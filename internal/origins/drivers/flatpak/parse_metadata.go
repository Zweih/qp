package flatpak

import (
	"bufio"
	"fmt"
	"os"
	"qp/internal/consts"
	"qp/internal/pkgdata"
	"strings"
)

func parseMetadata(pkgRef *PkgRef) (*PkgRef, error) {
	file, err := os.Open(pkgRef.MetadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s metadta file %s: %w", pkgRef.Name, pkgRef.MetadataPath, err)
	}

	defer file.Close()

	reason := consts.ReasonExplicit
	if pkgRef.Type == typeRuntime {
		reason = consts.ReasonDependency
	}

	pkg := &pkgdata.PkgInfo{
		Name:    pkgRef.Name,
		Reason:  reason,
		Arch:    pkgRef.Arch,
		Env:     fmt.Sprintf("%s (%s)", pkgRef.Scope, pkgRef.Branch),
		PkgType: pkgRef.Type,
	}

	scanner := bufio.NewScanner(file)

	var currentSection string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.Trim(line, "[]")
			if strings.HasPrefix(currentSection, sectionExtension) {
				applyExtensionSection(pkg, currentSection)
			}

			continue
		}

		if strings.Contains(line, "=") {
			key, value := parseKeyValue(line)
			applyMetadataField(pkg, currentSection, key, value)
			continue
		}
	}

	pkgRef.Pkg = pkg
	return pkgRef, nil
}

func parseKeyValue(line string) (string, string) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", ""
	}

	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
}

func applyMetadataField(
	pkg *pkgdata.PkgInfo,
	section string,
	key string,
	value string,
) {
	switch section {
	case sectionApplication:
		applyApplicationField(pkg, key, value)
	case sectionExtensionOf:
		applyExtensionOfField(pkg, key, value)
	default:
		// TODO: extension
	}
}

func applyApplicationField(pkg *pkgdata.PkgInfo, key string, value string) {
	switch key {
	case fieldRuntime:
		rel, err := parseRuntime(value)
		if err == nil {
			pkg.Depends = append(pkg.Depends, rel)
		}
	}
}

func applyExtensionOfField(pkg *pkgdata.PkgInfo, key string, value string) {
	switch key {
	case fieldRef:
		rel, err := parseRef(value)
		if err == nil {
			pkg.OptionalFor = append(pkg.OptionalFor, rel)
		}
	}
}

func applyExtensionSection(pkg *pkgdata.PkgInfo, section string) {
	sectionParts := strings.SplitN(section, " ", 2)
	if len(sectionParts) != 2 {
		return
	}

	extPart := sectionParts[1]
	lastDot := strings.LastIndex(extPart, ".")
	baseName := extPart[:lastDot]

	if baseName == pkg.Name {
		extType := extPart[lastDot:]

		switch extType {
		case ".Locale", ".Debug", ".Source":
			pkg.OptDepends = append(pkg.OptDepends, pkgdata.Relation{Name: extPart, PkgType: pkg.Name})
			return
		}
	}

	pkg.OptDepends = append(pkg.OptDepends, pkgdata.Relation{Name: extPart})
}

func parseRef(refDir string) (pkgdata.Relation, error) {
	parts := strings.SplitN(refDir, string(os.PathSeparator), 2)
	if len(parts) != 2 {
		return pkgdata.Relation{}, fmt.Errorf("malformed ref value: %s", refDir)
	}

	return parseRuntime(parts[1])
}

func parseRuntime(runtimeDir string) (pkgdata.Relation, error) {
	parts := strings.SplitN(runtimeDir, string(os.PathSeparator), 3)
	if len(parts) != 3 {
		return pkgdata.Relation{}, fmt.Errorf("malformed runtime value: %s", runtimeDir)
	}

	name := parts[0]
	version := parts[len(parts)-1]

	return pkgdata.Relation{Name: name, Version: version, Operator: pkgdata.OpEqual}, nil
}
