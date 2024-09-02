package main

import (
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		initRepo()
	case "cat-file":
		if len(os.Args) != 4 || os.Args[2] != "-p" {
			fmt.Fprintf(os.Stderr, "usage: mygit cat-file -p <object>\n")
			os.Exit(1)
		}
		catFile(os.Args[3])
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}

func initRepo() {
	for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			os.Exit(1)
		}
	}

	headFileContents := []byte("ref: refs/heads/main\n")
	if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Initialized git directory")
}

func catFile(objectHash string) {
	objectPath := filepath.Join(".git", "objects", objectHash[:2], objectHash[2:])
	file, err := os.Open(objectPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening object file: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()

	zlibReader, err := zlib.NewReader(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating zlib reader: %s\n", err)
		os.Exit(1)
	}
	defer zlibReader.Close()

	content, err := io.ReadAll(zlibReader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading object content: %s\n", err)
		os.Exit(1)
	}

	parts := strings.SplitN(string(content), "\x00", 2)
	if len(parts) != 2 {
		fmt.Fprintf(os.Stderr, "Invalid object format\n")
		os.Exit(1)
	}

	header := strings.SplitN(parts[0], " ", 2)
	if len(header) != 2 || header[0] != "blob" {
		fmt.Fprintf(os.Stderr, "Invalid object type\n")
		os.Exit(1)
	}

	size, err := strconv.Atoi(header[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid object size\n")
		os.Exit(1)
	}

	if len(parts[1]) != size {
		fmt.Fprintf(os.Stderr, "Object size mismatch\n")
		os.Exit(1)
	}

	fmt.Print(parts[1])
}
