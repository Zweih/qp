package npm

import (
	"os"
	"os/exec"
	"path/filepath"
	"qp/internal/consts"
	"runtime"
)

func getGlobalModulesDir() (string, error) {
	if prefix := os.Getenv(envPrefix); prefix != "" {
		return getNodeModulesDir(prefix), nil
	}

	nodeBinPath, err := getNodeBin()
	if err != nil {
		return "", err
	}

	prefix := getPrefix(nodeBinPath)

	if runtime.GOOS != consts.Windows {
		if destDir := os.Getenv(envDestDir); destDir != "" {
			prefix = filepath.Join(destDir, prefix)
		}
	}

	return getNodeModulesDir(prefix), nil
}

func getNodeBin() (string, error) {
	nodeBinName := nodeBinUnix
	if runtime.GOOS == consts.Windows {
		nodeBinName = nodeBinWin
	}

	execPath, err := exec.LookPath(nodeBinName)
	if err != nil {
		return "", err
	}

	resolvedPath, err := filepath.EvalSymlinks(execPath)
	if err != nil {
		return execPath, nil
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
