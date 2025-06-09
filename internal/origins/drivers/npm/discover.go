package npm

import (
	"os"
	"os/exec"
	"path/filepath"
	"qp/internal/consts"
	"runtime"
	"strings"
)

func getGlobalModulesDirs() ([]string, error) {
	var modulesDirs []string

	if prefix := os.Getenv(envPrefix); prefix != "" {
		modulesDirs = append(modulesDirs, getNodeModulesDir(prefix))
	}

	nodeBinPaths, err := getAllBins("node")
	if err != nil {
		return modulesDirs, err
	}

	for _, nodeBinPath := range nodeBinPaths {
		prefix := getPrefix(nodeBinPath)

		if runtime.GOOS != consts.Windows {
			if destDir := os.Getenv(envDestDir); destDir != "" {
				prefix = filepath.Join(destDir, prefix)
			}
		}

		modulesDirs = append(modulesDirs, getNodeModulesDir(prefix))
	}

	return modulesDirs, nil
}

func getNodeBin() (string, error) {
	nodeBinName := nodeBinUnix
	if runtime.GOOS == consts.Windows {
		nodeBinName = nodeBinWin
	}

	binPath, err := exec.LookPath(nodeBinName)
	if err != nil {
		return "", err
	}

	resolvedPath, err := filepath.EvalSymlinks(binPath)
	if err != nil {
		return binPath, nil
	}

	return resolvedPath, nil
}

func getPrefix(nodeBinPath string) string {
	dir := filepath.Dir(nodeBinPath)

	if runtime.GOOS != consts.Windows {
		dir = filepath.Dir(dir) // on non-windows it's dirname/dirname/bin
	}

	return dir
}

func getNodeModulesDir(prefix string) string {
	if runtime.GOOS == consts.Windows {
		return filepath.Join(prefix, winModulesDir)
	}

	return filepath.Join(prefix, unixModulesDir)
}

func getAllBins(binName string) ([]string, error) {
	var bins []string
	seen := make(map[string]bool)

	if runtime.GOOS == consts.Windows && !strings.HasSuffix(binName, ".exe") {
		binName += ".exe"
	}

	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return bins, nil
	}

	pathDirs := filepath.SplitList(pathEnv)

	for _, dir := range pathDirs {
		if dir == "" {
			continue
		}

		binPath := filepath.Join(dir, binName)

		if info, err := os.Stat(binPath); err == nil {
			if info.Mode()&0111 != 0 {
				resolvedPath, err := filepath.EvalSymlinks(binPath)
				if err != nil {
					resolvedPath = binPath
				}

				if !seen[resolvedPath] {
					seen[resolvedPath] = true
					bins = append(bins, resolvedPath)
				}
			}
		}
	}

	return bins, nil
}
