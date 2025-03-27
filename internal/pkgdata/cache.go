package pkgdata

import (
	"errors"
	"fmt"
	"os"
	pb "yaylog/internal/protobuf"

	"google.golang.org/protobuf/proto"
)

const (
	cachePath    = "/tmp/yaylog.cache"
	cacheVersion = 1
)

func getDbModTime() (int64, error) {
	dirInfo, err := os.Stat(PacmanDbPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read pacman DB mod time: %v", err)
	}

	return dirInfo.ModTime().Unix(), nil
}

func SaveProtoCache(pkgs []*PkgInfo) error {
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

func LoadProtoCache() ([]*PkgInfo, error) {
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
			Name:     rel.Name,
			Version:  rel.Version,
			Operator: pb.RelationOp(rel.Operator),
		}
	}

	return pbRels
}

func pkgsToProtos(pkgs []*PkgInfo) []*pb.PkgInfo {
	pbPkgs := make([]*pb.PkgInfo, len(pkgs))
	for i, pkg := range pkgs {
		pbPkgs[i] = &pb.PkgInfo{
			Timestamp:  pkg.Timestamp,
			Size:       pkg.Size,
			Name:       pkg.Name,
			Reason:     pkg.Reason,
			Version:    pkg.Version,
			Arch:       pkg.Arch,
			License:    pkg.License,
			Url:        pkg.Url,
			Depends:    relationsToProtos(pkg.Depends),
			RequiredBy: relationsToProtos(pkg.RequiredBy),
			Provides:   relationsToProtos(pkg.Provides),
			Conflicts:  relationsToProtos(pkg.Conflicts),
		}
	}

	return pbPkgs
}

func protosToRelations(pbRels []*pb.Relation) []Relation {
	rels := make([]Relation, len(pbRels))
	for i, pbRel := range pbRels {
		rels[i] = Relation{
			Name:     pbRel.Name,
			Version:  pbRel.Version,
			Operator: RelationOp(pbRel.Operator),
		}
	}

	return rels
}

func protosToPkgs(pbPkgs []*pb.PkgInfo) []*PkgInfo {
	pkgs := make([]*PkgInfo, len(pbPkgs))
	for i, pbPkg := range pbPkgs {
		pkgs[i] = &PkgInfo{
			Timestamp:  pbPkg.Timestamp,
			Size:       pbPkg.Size,
			Name:       pbPkg.Name,
			Reason:     pbPkg.Reason,
			Version:    pbPkg.Version,
			Arch:       pbPkg.Arch,
			License:    pbPkg.License,
			Url:        pbPkg.Url,
			Depends:    protosToRelations(pbPkg.Depends),
			RequiredBy: protosToRelations(pbPkg.RequiredBy),
			Provides:   protosToRelations(pbPkg.Provides),
			Conflicts:  protosToRelations(pbPkg.Conflicts),
		}
	}

	return pkgs
}
