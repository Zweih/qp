package pkgdata

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	pb "qp/internal/protobuf"
	"runtime"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"
)

const (
	cacheVersion    = 22 // bump when updating structure of PkgInfo/Relation/pkginfo.proto OR when dependency resolution is updated
	historyVersion  = 2
	xdgCacheHomeEnv = "XDG_CACHE_HOME"
	homeEnv         = "HOME"
	sudoUserEnv     = "SUDO_USER"
	userEnv         = "USER"
	qpCacheDir      = "query-packages"
	dotCache        = ".cache"
	dotModTime      = ".modtime"
	dotHistory      = ".history"
	dotLock         = ".lock"
	darwinCacheDir  = "Library/Caches"
)

func GetCachePath() (string, error) {
	userCacheDir, err := GetUserCachePath()
	if err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	cachePath := filepath.Join(userCacheDir, qpCacheDir)
	if err := os.MkdirAll(cachePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	return cachePath, nil
}

func GetUserCachePath() (string, error) {
	err := switchToRealUser()
	if err != nil {
		return "", err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		home = os.Getenv(homeEnv)
	}

	if runtime.GOOS == "darwin" {
		return filepath.Join(home, darwinCacheDir), nil
	}

	userCacheDir := os.Getenv(xdgCacheHomeEnv)
	if userCacheDir == "" {
		userCacheDir = filepath.Join(home, dotCache)
	}

	return userCacheDir, nil
}

func switchToRealUser() error {
	realUser := os.Getenv(sudoUserEnv)
	if realUser == "" {
		return nil
	}

	usr, err := user.Lookup(realUser)
	if err != nil {
		return err
	}

	os.Setenv(homeEnv, usr.HomeDir)
	os.Setenv(userEnv, realUser)
	os.Unsetenv(xdgCacheHomeEnv)

	return nil
}

func SaveCacheModTime(cacheRoot string, modTime int64) error {
	modPath := cacheRoot + dotModTime
	return os.WriteFile(modPath, []byte(strconv.FormatInt(modTime, 10)), 0644)
}

func LoadCacheModTime(cacheRoot string) (int64, error) {
	modPath := cacheRoot + dotModTime
	data, err := os.ReadFile(modPath)
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
}

func LoadProtoCache(cacheRoot string) ([]*PkgInfo, error) {
	cachePath := cacheRoot + dotCache
	if cachePath == "" {
		return nil, errors.New("invalid cache path, skipping cache load")
	}

	byteData, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, err
	}

	var cachedPkgs pb.CachedPkgs
	err = proto.Unmarshal(byteData, &cachedPkgs)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal cache: %v", err)
	}

	if cachedPkgs.Version != cacheVersion {
		return nil, errors.New("cache version mismatch, regenerating fresh cache")
	}

	pkgs := protosToPkgs(cachedPkgs.Pkgs)
	return pkgs, nil
}

func SaveProtoCache(cacheRoot string, pkgs []*PkgInfo) error {
	cachePath := cacheRoot + dotCache

	cachedPkgs := &pb.CachedPkgs{
		Pkgs:    pkgsToProtos(pkgs),
		Version: cacheVersion,
	}

	byteData, err := proto.Marshal(cachedPkgs)
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %v", cachedPkgs)
	}

	return os.WriteFile(cachePath, byteData, 0644)
}

func relationsToProtos(rels []Relation) []*pb.Relation {
	pbRels := make([]*pb.Relation, len(rels))
	for i, rel := range rels {
		pbRels[i] = &pb.Relation{
			Operator:     pb.RelationOp(rel.Operator),
			Depth:        rel.Depth,
			Name:         rel.Name,
			Version:      rel.Version,
			ProviderName: rel.ProviderName,
			Why:          rel.Why,
		}
	}

	return pbRels
}

func pkgsToProtos(pkgs []*PkgInfo) []*pb.PkgInfo {
	pbPkgs := make([]*pb.PkgInfo, len(pkgs))
	for i, pkg := range pkgs {
		pbPkgs[i] = &pb.PkgInfo{
			UpdateTimestamp: pkg.UpdateTimestamp,
			BuildTimestamp:  pkg.BuildTimestamp,
			Size:            pkg.Size,
			Name:            pkg.Name,
			Reason:          pkg.Reason,
			Version:         pkg.Version,
			Origin:          pkg.Origin,
			Arch:            pkg.Arch,
			License:         pkg.License,
			Url:             pkg.Url,
			Description:     pkg.Description,
			Validation:      pkg.Validation,
			PkgType:         pkg.PkgType,
			PkgBase:         pkg.PkgBase,
			Packager:        pkg.Packager,
			Groups:          pkg.Groups,
			Conflicts:       relationsToProtos(pkg.Conflicts),
			Replaces:        relationsToProtos(pkg.Replaces),
			Depends:         relationsToProtos(pkg.Depends),
			OptDepends:      relationsToProtos(pkg.OptDepends),
			RequiredBy:      relationsToProtos(pkg.RequiredBy),
			OptionalFor:     relationsToProtos(pkg.OptionalFor),
			Provides:        relationsToProtos(pkg.Provides),
		}
	}

	return pbPkgs
}

func protosToRelations(pbRels []*pb.Relation) []Relation {
	rels := make([]Relation, len(pbRels))
	for i, pbRel := range pbRels {
		rels[i] = Relation{
			Operator:     RelationOp(pbRel.Operator),
			Depth:        pbRel.Depth,
			Name:         pbRel.Name,
			Version:      pbRel.Version,
			ProviderName: pbRel.ProviderName,
			Why:          pbRel.Why,
		}
	}

	return rels
}

func protosToPkgs(pbPkgs []*pb.PkgInfo) []*PkgInfo {
	pkgs := make([]*PkgInfo, len(pbPkgs))
	for i, pbPkg := range pbPkgs {
		pkgs[i] = &PkgInfo{
			UpdateTimestamp: pbPkg.UpdateTimestamp,
			BuildTimestamp:  pbPkg.BuildTimestamp,
			Size:            pbPkg.Size,
			Name:            pbPkg.Name,
			Reason:          pbPkg.Reason,
			Version:         pbPkg.Version,
			Origin:          pbPkg.Origin,
			Arch:            pbPkg.Arch,
			License:         pbPkg.License,
			Url:             pbPkg.Url,
			Description:     pbPkg.Description,
			Validation:      pbPkg.Validation,
			PkgType:         pbPkg.PkgType,
			PkgBase:         pbPkg.PkgBase,
			Packager:        pbPkg.Packager,
			Groups:          pbPkg.Groups,
			Conflicts:       protosToRelations(pbPkg.Conflicts),
			Replaces:        protosToRelations(pbPkg.Replaces),
			Depends:         protosToRelations(pbPkg.Depends),
			OptDepends:      protosToRelations(pbPkg.OptDepends),
			RequiredBy:      protosToRelations(pbPkg.RequiredBy),
			OptionalFor:     protosToRelations(pbPkg.OptionalFor),
			Provides:        protosToRelations(pbPkg.Provides),
		}
	}

	return pkgs
}

func SaveInstallHistory(cacheRoot string, history map[string]int64, latestLogTime int64) error {
	historyPath := cacheRoot + dotHistory
	installHistory := &pb.InstallHistory{
		InstallTimestamps:  history,
		Version:            historyVersion,
		LatestLogTimestamp: latestLogTime,
	}

	byteData, err := proto.Marshal(installHistory)
	if err != nil {
		return fmt.Errorf("failed to marshal history: %v", err)
	}

	return os.WriteFile(historyPath, byteData, 0644)
}

func LoadInstallHistory(cacheRoot string) (map[string]int64, int64, error) {
	historyPath := cacheRoot + dotHistory
	if _, err := os.Stat(historyPath); os.IsNotExist(err) {
		return make(map[string]int64), 0, nil
	}

	byteData, err := os.ReadFile(historyPath)
	if err != nil {
		return nil, 0, err
	}

	var installHistory pb.InstallHistory
	err = proto.Unmarshal(byteData, &installHistory)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to unmarshal install history: %v", err)
	}

	if installHistory.Version != historyVersion {
		return make(map[string]int64), 0, nil
	}

	return installHistory.InstallTimestamps, installHistory.LatestLogTimestamp, nil
}

func IsLockFileExists(cacheRoot string) bool {
	lockPath := cacheRoot + dotLock
	_, err := os.Stat(lockPath)
	return err == nil
}

func CreateLockFile(cacheRoot string) error {
	lockPath := cacheRoot + dotLock
	pid := os.Getpid()
	return os.WriteFile(lockPath, []byte(strconv.Itoa(pid)), 0644)
}

func RemoveLockFile(cacheRoot string) error {
	lockPath := cacheRoot + dotLock
	return os.Remove(lockPath)
}
