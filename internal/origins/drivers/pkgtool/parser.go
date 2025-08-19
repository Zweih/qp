package pkgtool

import (
	"bytes"
	"fmt"
	"os"
	"qp/internal/consts"
	"qp/internal/pkgdata"
	"strconv"
	"strings"
	"unicode"
)

func parsePackageFile(packagePath string, origin string) (*pkgdata.PkgInfo, error) {
	data, err := os.ReadFile(packagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open package metadata: %w", err)
	}

	pkg := &pkgdata.PkgInfo{
		Origin: origin,
		Reason: consts.ReasonExplicit,
	}

	inDesc := false
	var description []string

	for line := range bytes.SplitSeq(data, []byte("\n")) {
		parts := strings.SplitN(string(line), ":", 2)
		if len(parts) != 2 {
			continue
		}

		key, value := parts[0], parts[1]

		switch key {
		case fieldPackageName:
			extractMetadata(value, pkg)
		case fieldUncompressedSize:
			pkg.Size = extractSize(value)
		case fieldPackageDescription:
			inDesc = true
		case fieldFileList:
			// exit metadata file
			break
		default:
			if inDesc {
				description = append(description, strings.TrimSpace(value))
			}

		}
	}

	fileInfo, _ := os.Stat(packagePath)
	pkg.UpdateTimestamp = fileInfo.ModTime().Unix()
	pkg.Description = strings.Join(description, " ")

	return pkg, nil
}

func extractMetadata(value string, pkg *pkgdata.PkgInfo) {
	trimmed := strings.TrimSpace(value)
	parts := strings.Split(trimmed, "-")
	if len(parts) < 4 {
		return
	}

	pkg.Arch = parts[len(parts)-2]
	pkg.Version = parts[len(parts)-3]
	pkg.Name = strings.Join(parts[:len(parts)-3], "-")
}

func extractSize(value string) int64 {
	sizeStr := strings.TrimSpace(value)
	firstLetterIdx := strings.IndexFunc(sizeStr, func(char rune) bool {
		return unicode.IsLetter(char)
	})

	rawSize, err := strconv.ParseFloat(sizeStr[:firstLetterIdx], 64)
	if err != nil {
		return 0
	}

	unit := sizeStr[firstLetterIdx:]
	var size int64

	switch unit {
	case "G":
		size = int64(rawSize * consts.GB)
	case "M":
		size = int64(rawSize * consts.MB)
	case "K":
		size = int64(rawSize * consts.KB)
	case "B":
		fallthrough
	default:
		size = int64(rawSize)
	}

	return size
}
