package git

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

type TreeEntry struct {
	Mode string
	Name string
	Hash string
}

func LsTree(treeHash string) {
	repo := NewRepo(".")
	content, err := repo.ReadObject(treeHash)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading tree object: %s\n", err)
		os.Exit(1)
	}

	entries, err := parseTreeObject(content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing tree object: %s\n", err)
		os.Exit(1)
	}

	for _, entry := range entries {
		fmt.Println(entry.Name)
	}
}

func parseTreeObject(content []byte) ([]TreeEntry, error) {
	var entries []TreeEntry
	parts := bytes.SplitN(content, []byte{0}, 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid tree object format")
	}

	header := parts[0]
	if !bytes.HasPrefix(header, []byte("tree ")) {
		return nil, fmt.Errorf("invalid tree object header")
	}

	data := parts[1]
	for len(data) > 0 {
		nullIndex := bytes.IndexByte(data, 0)
		if nullIndex == -1 {
			return nil, fmt.Errorf("invalid tree entry format")
		}

		entryHeader := string(data[:nullIndex])
		data = data[nullIndex+1:]

		if len(data) < 20 {
			return nil, fmt.Errorf("invalid tree entry: insufficient data for SHA")
		}

		sha := hex.EncodeToString(data[:20])
		data = data[20:]

		parts := bytes.SplitN([]byte(entryHeader), []byte{' '}, 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid tree entry header")
		}

		mode := string(parts[0])
		name := string(parts[1])

		entries = append(entries, TreeEntry{
			Mode: mode,
			Name: name,
			Hash: sha,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})

	return entries, nil
}

func WriteTree() {
	repo := NewRepo(".")
	treeHash, err := writeTreeRecursive(repo, ".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing tree: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(treeHash)
}

func writeTreeRecursive(repo *Repo, dir string) (string, error) {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}

	var treeEntries []TreeEntry
	for _, entry := range entries {
		if entry.Name() == ".git" {
			continue
		}

		entryPath := filepath.Join(dir, entry.Name())
		var mode string
		var hash string

		if entry.IsDir() {
			mode = "40000"
			hash, err = writeTreeRecursive(repo, entryPath)
		} else {
			mode = "100644"
			if entry.Mode()&0111 != 0 {
				mode = "100755"
			}
			content, err := ioutil.ReadFile(entryPath)
			if err != nil {
				return "", err
			}
			hash, err = repo.WriteObject("blob", content)
		}

		if err != nil {
			return "", err
		}

		treeEntries = append(treeEntries, TreeEntry{
			Mode: mode,
			Name: entry.Name(),
			Hash: hash,
		})
	}

	sort.Slice(treeEntries, func(i, j int) bool {
		return treeEntries[i].Name < treeEntries[j].Name
	})

	var treeContent bytes.Buffer
	for _, entry := range treeEntries {
		hashBytes, _ := hex.DecodeString(entry.Hash)
		treeContent.WriteString(fmt.Sprintf("%s %s\x00", entry.Mode, entry.Name))
		treeContent.Write(hashBytes)
	}

	return repo.WriteObject("tree", treeContent.Bytes())
}
