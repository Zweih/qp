package storage

import (
	"qp/internal/pkgdata"
	pb "qp/internal/protobuf"
)

func relationsToProtos(rels []pkgdata.Relation) []*pb.Relation {
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

func pkgsToProtos(pkgs []*pkgdata.PkgInfo) []*pb.PkgInfo {
	pbPkgs := make([]*pb.PkgInfo, len(pkgs))
	for i, pkg := range pkgs {
		pbPkgs[i] = &pb.PkgInfo{
			InstallTimestamp: pkg.InstallTimestamp,
			UpdateTimestamp:  pkg.UpdateTimestamp,
			BuildTimestamp:   pkg.BuildTimestamp,
			Size:             pkg.Size,
			Name:             pkg.Name,
			Reason:           pkg.Reason,
			Version:          pkg.Version,
			Origin:           pkg.Origin,
			Arch:             pkg.Arch,
			Env:              pkg.Env,
			License:          pkg.License,
			Url:              pkg.Url,
			Description:      pkg.Description,
			Validation:       pkg.Validation,
			PkgType:          pkg.PkgType,
			PkgBase:          pkg.PkgBase,
			Packager:         pkg.Packager,
			Groups:           pkg.Groups,
			AlsoIn:           pkg.AlsoIn,
			OtherEnvs:        pkg.OtherEnvs,
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

func protosToRelations(pbRels []*pb.Relation) []pkgdata.Relation {
	rels := make([]pkgdata.Relation, len(pbRels))
	for i, pbRel := range pbRels {
		rels[i] = pkgdata.Relation{
			Operator:     pkgdata.RelationOp(pbRel.Operator),
			Depth:        pbRel.Depth,
			Name:         pbRel.Name,
			Version:      pbRel.Version,
			ProviderName: pbRel.ProviderName,
			Why:          pbRel.Why,
		}
	}

	return rels
}

func protosToPkgs(pbPkgs []*pb.PkgInfo) []*pkgdata.PkgInfo {
	pkgs := make([]*pkgdata.PkgInfo, len(pbPkgs))
	for i, pbPkg := range pbPkgs {
		pkgs[i] = &pkgdata.PkgInfo{
			InstallTimestamp: pbPkg.InstallTimestamp,
			UpdateTimestamp:  pbPkg.UpdateTimestamp,
			BuildTimestamp:   pbPkg.BuildTimestamp,
			Size:             pbPkg.Size,
			Name:             pbPkg.Name,
			Reason:           pbPkg.Reason,
			Version:          pbPkg.Version,
			Origin:           pbPkg.Origin,
			Arch:             pbPkg.Arch,
			Env:              pbPkg.Env,
			License:          pbPkg.License,
			Url:              pbPkg.Url,
			Description:      pbPkg.Description,
			Validation:       pbPkg.Validation,
			PkgType:          pbPkg.PkgType,
			PkgBase:          pbPkg.PkgBase,
			Packager:         pbPkg.Packager,
			Groups:           pbPkg.Groups,
			AlsoIn:           pbPkg.AlsoIn,
			OtherEnvs:        pbPkg.OtherEnvs,
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
