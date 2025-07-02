package snap

import (
	"fmt"
	"os"
	"path/filepath"
	"qp/internal/consts"
	"qp/internal/pkgdata"
	"time"
)

type SnapdResponse struct {
	Type   string     `json:"type"`
	Status string     `json:"status"`
	Result []SnapInfo `json:"result"`
}

type SnapInfo struct {
	InstalledSize int64     `json:"installed-size"`
	Name          string    `json:"name"`
	Title         string    `json:"title"`
	Base          string    `json:"base"`
	Type          string    `json:"type"`
	Version       string    `json:"version"`
	Revision      string    `json:"revision"`
	InstallDate   string    `json:"install-date"`
	RefreshDate   string    `json:"refresh-date"`
	Status        string    `json:"status"`
	License       string    `json:"license"`
	Summary       string    `json:"summary"`
	Description   string    `json:"description"`
	Developer     string    `json:"developer"`
	Publisher     Publisher `json:"publisher"`
	Links         Links     `json:"links"`
}

type Publisher struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display-name"`
	Validation  string `json:"validation"`
}

type Links struct {
	Website []string `json:"website"`
	Source  []string `json:"source"`
	Issues  []string `json:"issues"`
}

type ConnectionsResponse struct {
	Type   string `json:"type"`
	Status string `json:"status"`
	Result Result `json:"result"`
}

type Result struct {
	Established []Connection `json:"established"`
}

type Connection struct {
	Interface string `json:"interface"`
	Slot      Slot   `json:"slot"`
	Plug      Plug   `json:"plug"`
}

type Slot struct {
	Snap string `json:"snap"`
	Slot string `json:"slot"`
}

type Plug struct {
	Snap string `json:"snap"`
	Plug string `json:"plug"`
}

func parseSnaps(snapResp SnapdResponse, deps map[string][]string) ([]*pkgdata.PkgInfo, error) {
	pkgs := make([]*pkgdata.PkgInfo, 0, len(snapResp.Result))

	for _, snapInfo := range snapResp.Result {
		depends := make([]pkgdata.Relation, 0, len(deps[snapInfo.Name]))
		for _, dep := range deps[snapInfo.Name] {
			depends = append(depends, pkgdata.Relation{
				Name:  dep,
				Depth: 1,
			})
		}

		reason := consts.ReasonExplicit
		if snapInfo.Type != typeOs && snapInfo.Type != typeSnapd {
			binPath := filepath.Join(binDir, snapInfo.Name)
			_, err := os.Stat(binPath)
			if err != nil {
				reason = consts.ReasonDependency
			}
		}

		if snapInfo.Base != "" {
			depends = append(depends, pkgdata.Relation{
				Name:  snapInfo.Base,
				Depth: 1,
			})
		}

		var updateTimestamp int64
		snapFile := fmt.Sprintf("%s_%s%s", snapInfo.Name, snapInfo.Revision, dotSnap)
		snapPath := filepath.Join(snapsDir, snapFile)
		fileInfo, err := os.Stat(snapPath)
		if err == nil {
			updateTimestamp = fileInfo.ModTime().Unix()
		}

		pkg := &pkgdata.PkgInfo{
			UpdateTimestamp:  updateTimestamp,
			InstallTimestamp: parseTimestamp(snapInfo.InstallDate),
			Size:             snapInfo.InstalledSize,
			Name:             snapInfo.Name,
			Title:            snapInfo.Title,
			Reason:           reason,
			Version:          snapInfo.Version,
			Origin:           consts.OriginSnap,
			License:          snapInfo.License,
			Description:      snapInfo.Summary,
			Author:           snapInfo.Developer,
			Url:              getUrl(snapInfo.Links),
			Packager:         snapInfo.Publisher.DisplayName,
			PkgType:          snapInfo.Type,
			Validation:       snapInfo.Publisher.Validation,
			Depends:          depends,
		}

		pkgs = append(pkgs, pkg)
	}

	return pkgs, nil
}

func parseConnections(connsResp ConnectionsResponse) map[string][]string {
	deps := make(map[string][]string)
	for _, conn := range connsResp.Result.Established {
		if conn.Interface == interfaceContent {
			deps[conn.Plug.Snap] = append(deps[conn.Plug.Snap], conn.Slot.Snap)
		}
	}

	return deps
}

func getUrl(links Links) string {
	for _, link := range links.Website {
		return link
	}

	for _, link := range links.Source {
		return link
	}

	for _, link := range links.Issues {
		return link
	}

	return ""
}

func parseTimestamp(timeStr string) int64 {
	if timeStr == "" {
		return 0
	}

	timestamp, err := time.Parse(timestampFormat, timeStr)
	if err != nil {
		timestamp, err = time.Parse(time.RFC3339Nano, timeStr)

		if err == nil {
			return timestamp.Unix()
		}

		return 0
	}

	return timestamp.Unix()
}
