package asset

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"strings"
)

type Asset struct {
	Hash string
	Path string
}

func NewAsset(path string, hashLen int) (*Asset, error) {
	path = strings.TrimPrefix(path, "/")
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return nil, err
	}

	hash := hex.EncodeToString(hasher.Sum(nil))

	return &Asset{
		Path: path,
		Hash: hash[0:hashLen],
	}, nil
}
