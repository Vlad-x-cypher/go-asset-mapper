package asset

import (
	"encoding/json"
	"os"
)

type ManifestType int

const (
	ViteManifestType ManifestType = iota
	WebpackManifestType
)

// ManifestConfig provides information about manifest filepath and how to parse it correctly.
// Currently supports Vite or Webpack types.
type ManifestConfig struct {
	// manifest filepath
	Path string
	// manifest generator type
	Type ManifestType
}

type viteManifestRecord struct {
	File           string   `json:"file"`
	Src            string   `json:"src"`
	Name           string   `json:"name"`
	IsEntry        bool     `json:"isEntry"`
	CSS            []string `json:"css"`
	Imports        []string `json:"imports"`
	IsDynamicEntry bool     `json:"isDynamicEntry"`
}

func parseViteManifest(path string, a *AssetMapper) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	for decoder.More() {
		var data map[string]viteManifestRecord

		err = decoder.Decode(&data)
		if err != nil {
			return err
		}

		for k, v := range data {
			asset := &Asset{
				Path:       k,
				PublicPath: a.PublicPath + v.File,
				Hash:       "",
			}

			a.Assets[k] = asset
			if v.IsEntry {
				entry := a.CreateEntry(v.Name)
				entry.Add(asset.PublicPath)

				for _, css := range v.CSS {
					entry.Add(a.PublicPath + css)
				}
			}

		}
	}

	return nil
}

func parseWebpackManifest(path string, a *AssetMapper) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	for decoder.More() {
		var data map[string]string

		err = decoder.Decode(&data)
		if err != nil {
			return err
		}

		for k := range data {
			asset := &Asset{
				Path:       k,
				PublicPath: a.PublicPath,
				Hash:       "",
			}

			a.Assets[k] = asset
		}

	}
	return nil
}
