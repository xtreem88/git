package git

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
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
