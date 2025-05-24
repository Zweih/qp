package shared

import (
	"fmt"
	"os"
	"path/filepath"
)

type Node struct {
	Depth int
	Path  string
	Entry os.DirEntry
}

func BfsStale(path string, cacheMtime int64, maxDepth int) (bool, error) {
	parentInfo, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("failed to stat %s: %w", path, err)
	}

	if parentInfo.ModTime().Unix() > cacheMtime {
		return true, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return false, fmt.Errorf("failed to read %s: %w", path, err)
	}

	bfsQueue := make([]Node, 0, len(entries))
	for _, entry := range entries {
		bfsQueue = append(bfsQueue, Node{
			Depth: 1,
			Path:  filepath.Join(path, entry.Name()),
			Entry: entry,
		})
	}

	for 0 < len(bfsQueue) {
		node := bfsQueue[0]
		bfsQueue = bfsQueue[1:]

		if node.Depth > maxDepth {
			return false, nil
		}

		if !node.Entry.IsDir() {
			continue
		}

		info, err := node.Entry.Info()
		if err != nil {
			continue
		}

		modTime := info.ModTime().Unix()
		if modTime > cacheMtime {
			return true, nil
		}

		subEntries, err := os.ReadDir(node.Path)
		if err != nil {
			return false, fmt.Errorf("failed to read %s: %w", node.Path, err)
		}

		for _, subEntry := range subEntries {
			bfsQueue = append(bfsQueue, Node{
				Depth: node.Depth + 1,
				Path:  filepath.Join(node.Path, subEntry.Name()),
				Entry: subEntry,
			})
		}
	}

	return false, nil
}

func IsDirStale(origin string, path string, cacheMtime int64) (bool, error) {
	dirInfo, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("failed to read %s DB modified time: %v", origin, err)
	}

	return dirInfo.ModTime().Unix() > cacheMtime, nil
}
