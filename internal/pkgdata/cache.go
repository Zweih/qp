package pkgdata

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	pb "qp/internal/protobuf"
	"runtime"

	"google.golang.org/protobuf/proto"
)

const (
	cacheVersion    = 21 // bump when updating structure of PkgInfo/Relation/pkginfo.proto OR when dependency resolution is updated
	xdgCacheHomeEnv = "XDG_CACHE_HOME"
	homeEnv         = "HOME"
	qpCacheDir      = "query-packages"
)

func GetCachePath() (string, error) {
	cachePath := filepath.Join(GetBaseCachePath(), qpCacheDir)
	if err := os.MkdirAll(cachePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	return cachePath, nil
}

func GetBaseCachePath() string {
	if runtime.GOOS == "darwin" {
		return filepath.Join(os.Getenv(homeEnv), "Library/Caches")
	}

	userCacheDir := os.Getenv(xdgCacheHomeEnv)
	if userCacheDir == "" {
		userCacheDir = filepath.Join(os.Getenv(homeEnv), ".cache")
	}

	return userCacheDir
}

func SaveProtoCache(pkgs []*PkgInfo, cachePath string, lastModified int64) error {
	if cachePath == "" {
		return errors.New("invalid cache path, skipping cache save")
	}

	cachedPkgs := &pb.CachedPkgs{
		Pkgs:         pkgsToProtos(pkgs),
		LastModified: lastModified,
		Version:      cacheVersion,
	}

	byteData, err := proto.Marshal(cachedPkgs)
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %v", cachedPkgs)
	}

	return os.WriteFile(cachePath, byteData, 0644)
}

func LoadProtoCache(cachePath string, sourceModTime int64) ([]*PkgInfo, error) {
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

	if sourceModTime > cachedPkgs.LastModified {
		return nil, errors.New("cache is stale")
	}

	pkgs := protosToPkgs(cachedPkgs.Pkgs)

	return pkgs, nil
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
			InstallTimestamp: pkg.InstallTimestamp,
			BuildTimestamp:   pkg.BuildTimestamp,
			Size:             pkg.Size,
			Name:             pkg.Name,
			Reason:           pkg.Reason,
			Version:          pkg.Version,
			Origin:           pkg.Origin,
			Arch:             pkg.Arch,
			License:          pkg.License,
			Url:              pkg.Url,
			Description:      pkg.Description,
			Validation:       pkg.Validation,
			PkgType:          pkg.PkgType,
			PkgBase:          pkg.PkgBase,
			Packager:         pkg.Packager,
			Groups:           pkg.Groups,
			Conflicts:        relationsToProtos(pkg.Conflicts),
			Replaces:         relationsToProtos(pkg.Replaces),
			Depends:          relationsToProtos(pkg.Depends),
			OptDepends:       relationsToProtos(pkg.OptDepends),
			RequiredBy:       relationsToProtos(pkg.RequiredBy),
			OptionalFor:      relationsToProtos(pkg.OptionalFor),
			Provides:         relationsToProtos(pkg.Provides),
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
			InstallTimestamp: pbPkg.InstallTimestamp,
			BuildTimestamp:   pbPkg.BuildTimestamp,
			Size:             pbPkg.Size,
			Name:             pbPkg.Name,
			Reason:           pbPkg.Reason,
			Version:          pbPkg.Version,
			Origin:           pbPkg.Origin,
			Arch:             pbPkg.Arch,
			License:          pbPkg.License,
			Url:              pbPkg.Url,
			Description:      pbPkg.Description,
			Validation:       pbPkg.Validation,
			PkgType:          pbPkg.PkgType,
			PkgBase:          pbPkg.PkgBase,
			Packager:         pbPkg.Packager,
			Groups:           pbPkg.Groups,
			Conflicts:        protosToRelations(pbPkg.Conflicts),
			Replaces:         protosToRelations(pbPkg.Replaces),
			Depends:          protosToRelations(pbPkg.Depends),
			OptDepends:       protosToRelations(pbPkg.OptDepends),
			RequiredBy:       protosToRelations(pbPkg.RequiredBy),
			OptionalFor:      protosToRelations(pbPkg.OptionalFor),
			Provides:         protosToRelations(pbPkg.Provides),
		}
	}

	return pkgs
}
