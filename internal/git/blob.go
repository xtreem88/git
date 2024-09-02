package git

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func CatFile(objectHash string) {
	repo := NewRepo(".")
	content, err := repo.ReadObject(objectHash)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading object: %s\n", err)
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

func HashObject(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
		os.Exit(1)
	}

	repo := NewRepo(".")
	hash, err := repo.WriteObject("blob", content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing object: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(hash)
}
