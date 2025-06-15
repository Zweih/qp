package flatpak

import (
	"fmt"
	"os"
	"path/filepath"
)

func applyTimestamps(pkgRef *PkgRef) error {
	activePath := filepath.Join(
		pkgRef.NameDir,
		pkgRef.Arch,
		pkgRef.Branch,
		"active",
	)

	info, err := os.Lstat(activePath)
	if err != nil {
		return fmt.Errorf("failed to stat activePath at %s: %w", activePath, err)
	}

	pkgRef.Pkg.UpdateTimestamp = info.ModTime().Unix()
	// TODO: add creation time, once one of our packages updates so we can actually test it
	return nil
}
