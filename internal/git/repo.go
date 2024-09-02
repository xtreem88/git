package git

import (
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Repo struct {
	Path string
}

func NewRepo(path string) *Repo {
	return &Repo{Path: path}
}

func (r *Repo) ReadObject(hash string) ([]byte, error) {
	objectPath := filepath.Join(r.Path, ".git", "objects", hash[:2], hash[2:])
	file, err := os.Open(objectPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	zlibReader, err := zlib.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer zlibReader.Close()

	return io.ReadAll(zlibReader)
}

func (r *Repo) WriteObject(objectType string, content []byte) (string, error) {
	header := fmt.Sprintf("%s %d\x00", objectType, len(content))
	store := append([]byte(header), content...)

	hash := sha1.Sum(store)
	hashString := fmt.Sprintf("%x", hash)

	objectPath := filepath.Join(r.Path, ".git", "objects", hashString[:2], hashString[2:])
	if err := os.MkdirAll(filepath.Dir(objectPath), 0755); err != nil {
		return "", err
	}

	file, err := os.Create(objectPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	zlibWriter := zlib.NewWriter(file)
	defer zlibWriter.Close()

	if _, err := zlibWriter.Write(store); err != nil {
		return "", err
	}

	return hashString, nil
}
