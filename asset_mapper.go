// Package asset provides asset mapper functionality to help using static assets (css, scripts, images)
// in Go templates
//
// Package inspired by Symfony AssetMapper, only much simpler, without importmap functionality and compiling assets.
package asset

import (
	"errors"
	"fmt"
	"html"
	"html/template"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	cssRe   = regexp.MustCompile(`\.css$`)
	jsRe    = regexp.MustCompile(`\.js$`)
	imageRe = regexp.MustCompile(`(\.webp|\.jpg|\.jpeg|\.jpe|\.jfif|\.jif|\.png|\.gif|\.tiff|\.tif|\.svg|\.avif)$`)
)

type AssetMapper struct {
	JSAssets    map[string]*Asset
	CSSAssets   map[string]*Asset
	ImageAssets map[string]*Asset
	OtherAssets map[string]*Asset
	PublicPath  string
	HashLen     int
}

func NewAssetMapper() *AssetMapper {
	return &AssetMapper{
		JSAssets:    map[string]*Asset{},
		CSSAssets:   map[string]*Asset{},
		ImageAssets: map[string]*Asset{},
		OtherAssets: map[string]*Asset{},
		PublicPath:  "/",
		HashLen:     10,
	}
}

func isCSS(path string) bool {
	return cssRe.MatchString(path)
}

func isJS(path string) bool {
	return jsRe.MatchString(path)
}

func isImage(path string) bool {
	return imageRe.MatchString(path)
}

// ScanDir walks directory and maps all files to AssetMapper, storing its path and hash.
func (a *AssetMapper) ScanDir(dirName string) error {
	err := filepath.Walk(dirName, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		asset, assetErr := NewAsset(path, a.HashLen)
		if assetErr != nil {
			return assetErr
		}

		switch {
		case isCSS(path):
			a.CSSAssets[path] = asset
			a.CSSAssets[a.PublicPath+path] = asset
		case isJS(path):
			a.JSAssets[path] = asset
			a.JSAssets[a.PublicPath+path] = asset
		case isImage(path):
			a.ImageAssets[path] = asset
			a.ImageAssets[a.PublicPath+path] = asset
		default:
			a.OtherAssets[path] = asset
			a.OtherAssets[a.PublicPath+path] = asset
		}

		return nil
	})

	return err
}

func (a *AssetMapper) extractAssetPathFromMap(m map[string]*Asset, search string) string {
	if asset, ok := m[search]; ok {
		path := asset.Path
		if !strings.HasPrefix(path, a.PublicPath) {
			path = a.PublicPath + path
		}
		return path + "?v=" + asset.Hash
	}
	return search
}

// CSSLink returns versioned css path. If asset not found returns raw path.
func (a *AssetMapper) CSSLink(path string) string {
	return a.extractAssetPathFromMap(a.CSSAssets, path)
}

// JSLink returns versioned javascript path. If asset not found returns raw path.
func (a *AssetMapper) JSLink(path string) string {
	return a.extractAssetPathFromMap(a.JSAssets, path)
}

// ImageLink returns versioned image path. If asset not found returns raw path.
func (a *AssetMapper) ImageLink(path string) string {
	return a.extractAssetPathFromMap(a.ImageAssets, path)
}

// OtherLink returns versioned file path. If asset not found returns raw path.
func (a *AssetMapper) OtherLink(path string) string {
	return a.extractAssetPathFromMap(a.OtherAssets, path)
}

// Get returns asset url including version. If asset not found returns path param as is.
func (a *AssetMapper) Get(path string) string {
	switch {
	case isCSS(path):
		return a.CSSLink(path)
	case isJS(path):
		return a.JSLink(path)
	case isImage(path):
		return a.ImageLink(path)
	}

	return a.OtherLink(path)
}

func attributeMapToString(m map[string]string) string {
	s := []string{}

	for k, v := range m {
		if k == "async" || k == "defer" {
			s = append(s, k)
			continue
		}
		s = append(s, fmt.Sprintf(`%s="%s"`, html.EscapeString(k), html.EscapeString(v)))
	}

	return strings.Join(s, " ")
}

func tagAttributes(attrs []string) (map[string]string, error) {
	attrMap := map[string]string{}

	if len(attrs)%2 != 0 {
		return nil, errors.New("attrs must be an even number of strings")
	}

	for i := 0; i < len(attrs); i += 2 {
		attrMap[attrs[i]] = attrs[i+1]
	}

	return attrMap, nil
}

// ScriptTag returns HTML script tag
// attrs param can be used to pass additional attributes to the tag. Unfortunately Go html.Template
// does not allow create new maps inside html templates, attrs must be an even number of strings
// reperesenting key value pairs.
//
// Example usage in template:
//
// <head>
//
//	{{ scriptTag "main.js" }}
//
//	<!-- Passing additional attributes to script tag -->
//	{{ scriptTag "other.js" "type" "module" "id" "other-script" }}
//	<!-- Should render: <script src="other.js" type="module" id="other-script"></script> -->
//
//	<!-- Example set defer or async attributes -->
//	{{ scriptTag "defered.js" "defer" "" }}
//		<!-- Produces: <script defer src="defered.js"></script> -->
//	{{ scriptTag "some-async.js" "async" "" }}
//		<!-- Produces: <script async src="some-async.js"></script> -->
//
// </head>
//
// For more complete examples follow example dir.
func (a *AssetMapper) ScriptTag(path string, attrs ...string) (template.HTML, error) {
	link := a.JSLink(path)

	attrMap, err := tagAttributes(attrs)
	if err != nil {
		return "", err
	}

	attrMap["src"] = link

	return template.HTML(fmt.Sprintf("<script %s></script>", attributeMapToString(attrMap))), nil
}

// LinkTag returns HTML link tag
// attrs param can be used to pass additional attributes to the tag. Unfortunately Go html.Template
// does not allow create new maps inside html templates, attrs must be an even number of strings
// reperesenting key value pairs.
//
// Example usage in template:
//
// <head>
//
//	{{ linkTag "style.css" }}
//
//	<!-- Passing additional attributes to link tag -->
//	{{ linkTag "homepage.css" "id" "homepage-css" "media" "screen" }}
//	<!-- Should render: <link href="homepage.css" id="homepage-css" media="screen" /> -->
//
// </head>
func (a *AssetMapper) LinkTag(path string, attrs ...string) (template.HTML, error) {
	link := a.Get(path)

	attrs = append([]string{"rel", "stylesheet"}, attrs...)
	attrMap, err := tagAttributes(attrs)
	if err != nil {
		return "", err
	}

	attrMap["href"] = link

	return template.HTML(fmt.Sprintf("<link %s />", attributeMapToString(attrMap))), nil
}
