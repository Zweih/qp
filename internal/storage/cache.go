package storage

import (
	"errors"
	"fmt"
	"os"
	"qp/internal/pkgdata"
	pb "qp/internal/protobuf"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"
)

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

func LoadProtoCache(cacheRoot string) ([]*pkgdata.PkgInfo, error) {
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

func SaveProtoCache(cacheRoot string, pkgs []*pkgdata.PkgInfo) error {
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
