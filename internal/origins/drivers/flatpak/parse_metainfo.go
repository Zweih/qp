package flatpak

import (
	"encoding/xml"
	"os"
	"qp/internal/pkgdata"
	"strings"
)

type MetainfoXml struct {
	XmlName xml.Name `xml:"component"`
	Id      string   `xml:"id"`
	Names   []struct {
		XmlLang string `xml:"xml:lang,attr"`
		Name    string `xml:",chardata"`
	} `xml:"name"`
	Summaries []struct {
		XmlLang string `xml:"xml:lang,attr"`
		Summary string `xml:",chardata"`
	} `xml:"summary"`
	Description struct {
		Paragraphs []string `xml:"p"`
	} `xml:"description"`
	ProjectLicense string `xml:"project_license"`
	DeveloperName  string `xml:"developer_name"`
	Urls           []struct {
		Type string `xml:"type,attr"`
		Url  string `xml:",chardata"`
	} `xml:"url"`
}

func parseMetainfo(pkgRef *PkgRef) {
	if pkgRef.MetainfoPath == "" {
		return
	}

	data, err := os.ReadFile(pkgRef.MetainfoPath)
	if err != nil {
		// return nil, fmt.Errorf("failed to read metainfo xml at %s: %w", pkgRef.MetainfoPath)
	}

	var component MetainfoXml

	if err = xml.Unmarshal(data, &component); err != nil {
	}

	applyMetainfo(pkgRef.Pkg, component)
}

func applyMetainfo(pkg *pkgdata.PkgInfo, metainfo MetainfoXml) {
	for _, summary := range metainfo.Summaries {
		if summary.XmlLang == "" {
			pkg.Description = summary.Summary
			break
		}
	}

	for _, name := range metainfo.Names {
		if name.XmlLang == "" {
			pkg.Title = name.Name
			break
		}
	}

	if metainfo.ProjectLicense != "" {
		pkg.License = strings.TrimSpace(metainfo.ProjectLicense)
	}

	if metainfo.DeveloperName != "" {
		pkg.Packager = metainfo.DeveloperName
	}

	for _, url := range metainfo.Urls {
		if url.Type == "homepage" {
			pkg.Url = url.Url
			break
		}
	}
}
