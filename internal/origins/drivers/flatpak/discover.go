package flatpak

import (
	"os"
	"path/filepath"
	"qp/internal/origins/worker"
	"qp/internal/pkgdata"
	"strings"
	"sync"
)

type PkgRef struct {
	Pkg          *pkgdata.PkgInfo
	Name         string
	Type         string
	Arch         string
	Branch       string
	Remote       string
	Scope        string
	CommitDir    string
	NameDir      string
	MetadataPath string
	MetainfoPath string
	DesktopPath  string
}

type PackageLocation struct {
	InstallDir string
	PkgType    string
	PkgDir     string
}

type PackageCommit struct {
	InstallDir string
	PkgType    string
	Name       string
	Arch       string
	Branch     string
	CommitDir  string
}

func discoverPackages(
	installDirs []string,
	estimatedCount int,
	errChan chan<- error,
	errGroup *sync.WaitGroup,
) <-chan *PkgRef {
	jobChan := make(chan PackageLocation, len(installDirs)*2)

	for _, installDir := range installDirs {
		jobChan <- PackageLocation{
			InstallDir: installDir,
			PkgType:    typeApp,
			PkgDir:     filepath.Join(installDir, typeApp),
		}

		jobChan <- PackageLocation{
			InstallDir: installDir,
			PkgType:    typeRuntime,
			PkgDir:     filepath.Join(installDir, typeRuntime),
		}
	}
	close(jobChan)

	stage1Out := worker.RunWorkers(
		jobChan,
		errChan,
		errGroup,
		scanPackageDirectory,
		0,
		len(installDirs)*2,
	)

	commitChan := flattenCommits(stage1Out, estimatedCount)

	stage2Out := worker.RunWorkers(
		commitChan,
		errChan,
		errGroup,
		validatePackageCommit,
		0,
		estimatedCount,
	)

	stage3Out := worker.RunWorkers(
		stage2Out,
		errChan,
		errGroup,
		findPackageFiles,
		0,
		estimatedCount,
	)

	stage4Out := worker.RunWorkers(
		stage3Out,
		errChan,
		errGroup,
		determinePackageRemote,
		0,
		estimatedCount,
	)

	return stage4Out
}

func scanPackageDirectory(loc PackageLocation) ([]PackageCommit, error) {
	var commits []PackageCommit

	nameEntries, err := os.ReadDir(loc.PkgDir)
	if err != nil {
		return nil, nil // dir doesn't exist, not an error
	}

	for _, nameEntry := range nameEntries {
		if !nameEntry.IsDir() {
			continue
		}

		nameDir := filepath.Join(loc.PkgDir, nameEntry.Name())

		archEntries, err := os.ReadDir(nameDir)
		if err != nil {
			continue
		}

		for _, archEntry := range archEntries {
			if !archEntry.IsDir() {
				continue
			}

			archDir := filepath.Join(nameDir, archEntry.Name())

			branchEntries, err := os.ReadDir(archDir)
			if err != nil {
				continue
			}

			for _, branchEntry := range branchEntries {
				if !branchEntry.IsDir() {
					continue
				}

				branchDir := filepath.Join(archDir, branchEntry.Name())

				commitEntries, err := os.ReadDir(branchDir)
				if err != nil {
					continue
				}

				for _, commitEntry := range commitEntries {
					if !commitEntry.IsDir() {
						continue
					}

					commits = append(commits, PackageCommit{
						InstallDir: loc.InstallDir,
						PkgType:    loc.PkgType,
						Name:       nameEntry.Name(),
						Arch:       archEntry.Name(),
						Branch:     branchEntry.Name(),
						CommitDir:  filepath.Join(branchDir, commitEntry.Name()),
					})
				}
			}
		}
	}

	return commits, nil
}

func flattenCommits(commitSlices <-chan []PackageCommit, bufferSize int) <-chan PackageCommit {
	out := make(chan PackageCommit, bufferSize)

	go func() {
		defer close(out)
		for commits := range commitSlices {
			for _, commit := range commits {
				out <- commit
			}
		}
	}()

	return out
}

func validatePackageCommit(commit PackageCommit) (*PkgRef, error) {
	metadataPath := filepath.Join(commit.CommitDir, metadataFile)
	if _, err := os.Stat(metadataPath); err != nil {
		return nil, worker.ErrSkip // no metadata, skip commit
	}

	if commit.PkgType == typeApp {
		activePath := filepath.Join(filepath.Dir(commit.CommitDir), activeFile)
		target, err := filepath.EvalSymlinks(activePath)

		if err != nil || target != commit.CommitDir {
			return nil, worker.ErrSkip // not the active version
		}
	}

	scope := scopeSystem
	if strings.Contains(commit.InstallDir, userInstallDir) {
		scope = scopeUser
	}

	return &PkgRef{
		Name:         commit.Name,
		Type:         commit.PkgType,
		Arch:         commit.Arch,
		Branch:       commit.Branch,
		Scope:        scope,
		CommitDir:    commit.CommitDir,
		NameDir:      filepath.Join(commit.InstallDir, commit.PkgType, commit.Name),
		MetadataPath: metadataPath,
	}, nil
}

func findPackageFiles(ref *PkgRef) (*PkgRef, error) {
	metainfoSearchPaths := []string{
		filepath.Join(ref.CommitDir, metainfoDir, ref.Name+dotMetainfoXml),
		filepath.Join(ref.CommitDir, metainfoDir, ref.Name+dotAppdataXml),
		filepath.Join(ref.CommitDir, appdataDir, ref.Name+dotAppdataXml),
	}

	for _, path := range metainfoSearchPaths {
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			ref.MetainfoPath = path
			break
		}
	}

	desktopPath := filepath.Join(ref.CommitDir, applicationsDir, ref.Name+dotDesktop)
	if info, err := os.Stat(desktopPath); err == nil && !info.IsDir() {
		ref.DesktopPath = desktopPath
	}

	return ref, nil
}

func determinePackageRemote(ref *PkgRef) (*PkgRef, error) {
	pathParts := strings.Split(ref.CommitDir, string(filepath.Separator))
	var installDir string

	for i, part := range pathParts {
		if part == ref.Type && i >= 1 {
			installDir = filepath.Join(pathParts[:i]...)
			break
		}
	}

	if installDir == "" {
		ref.Remote = remoteUnknown
		return ref, nil
	}

	remotesDir := filepath.Join(installDir, remotesDir)
	remoteEntries, err := os.ReadDir(remotesDir)
	if err != nil {
		ref.Remote = remoteUnknown
		return ref, nil
	}

	refPath := filepath.Join(ref.Type, ref.Name, ref.Arch, ref.Branch)

	for _, remoteEntry := range remoteEntries {
		if !remoteEntry.IsDir() {
			continue
		}

		remoteRefPath := filepath.Join(remotesDir, remoteEntry.Name(), refPath)
		if _, err := os.Stat(remoteRefPath); err == nil {
			ref.Remote = remoteEntry.Name()
			return ref, nil
		}
	}

	ref.Remote = remoteUnknown
	return ref, nil
}

func estimatePackageCount(installDirs []string) int {
	estimate := 0

	for _, installDir := range installDirs {
		for _, pkgType := range []string{typeApp, typeRuntime} {
			baseDir := filepath.Join(installDir, pkgType)
			if entries, err := os.ReadDir(baseDir); err == nil {
				// TODO: perhaps we should use 1.5 and round the result
				estimate += len(entries) * 2 // 1-2 commits per package on average
			}
		}
	}

	return max(estimate, 32)
}
