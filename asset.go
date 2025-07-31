package asset

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"strings"
)

type Asset struct {
	PublicPath string
	Hash       string
	Path       string
}

func NewAsset(path, publicPath string, hashLen int) (*Asset, error) {
	path = strings.TrimPrefix(path, "/")
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hash := ""
	if hashLen > 0 {
		hasher := sha256.New()
		if _, err := io.Copy(hasher, file); err != nil {
			return nil, err
		}

		hash = hex.EncodeToString(hasher.Sum(nil))
	}

	hash = hash[0:hashLen]

	pubPath := publicPath + path

	if hashLen > 0 {
		pubPath += "?v=" + hash
	}

	return &Asset{
		Path:       path,
		Hash:       hash,
		PublicPath: pubPath,
	}, nil
}
