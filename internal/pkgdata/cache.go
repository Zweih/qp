package pkgdata

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	pb "qp/internal/protobuf"

	"google.golang.org/protobuf/proto"
)

const (
	cacheVersion    = 13 // bump when updating structure of PkgInfo/Relation/pkginfo.proto OR when dependency resolution is updated
	xdgCacheHomeEnv = "XDG_CACHE_HOME"
	homeEnv         = "HOME"
	qpCacheDir      = "query-packages"
	packageManager  = "pacman"
)

func GetCachePath() (string, error) {
	userCacheDir := os.Getenv(xdgCacheHomeEnv)
	if userCacheDir == "" {
		userCacheDir = filepath.Join(os.Getenv(homeEnv), ".cache")
	}

	cachePath := filepath.Join(userCacheDir, qpCacheDir)
	if err := os.MkdirAll(cachePath, 0755); err != nil {
		return "", err
	}

	cacheFileName := "qp-" + packageManager + ".cache"

	return filepath.Join(cachePath, cacheFileName), nil
}

func getDbModTime() (int64, error) {
	dirInfo, err := os.Stat(PacmanDbPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read pacman DB mod time: %v", err)
	}

	return dirInfo.ModTime().Unix(), nil
}

func SaveProtoCache(pkgs []*PkgInfo, cachePath string) error {
	if cachePath == "" {
		return errors.New("invalid cache path, skipping cache save")
	}

	lastModified, err := getDbModTime()
	if err != nil {
		return err
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

func LoadProtoCache(cachePath string) ([]*PkgInfo, error) {
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

	dbModTime, err := getDbModTime()
	if err != nil {
		return nil, err
	}

	if dbModTime > cachedPkgs.LastModified {
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
			PkgType:          pb.PkgType(pkg.PkgType),
			Arch:             pkg.Arch,
			License:          pkg.License,
			Url:              pkg.Url,
			Description:      pkg.Description,
			PkgBase:          pkg.PkgBase,
			Validation:       pkg.Validation,
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
			PkgType:          PkgType(pbPkg.PkgType),
			Arch:             pbPkg.Arch,
			License:          pbPkg.License,
			Url:              pbPkg.Url,
			Description:      pbPkg.Description,
			PkgBase:          pbPkg.PkgBase,
			Validation:       pbPkg.Validation,
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
