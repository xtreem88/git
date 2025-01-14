package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/git-starter-go/internal/git"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		git.Init()
	case "cat-file":
		if len(os.Args) != 4 || os.Args[2] != "-p" {
			fmt.Fprintf(os.Stderr, "usage: mygit cat-file -p <object>\n")
			os.Exit(1)
		}
		git.CatFile(os.Args[3])
	case "hash-object":
		if len(os.Args) != 4 || os.Args[2] != "-w" {
			fmt.Fprintf(os.Stderr, "usage: mygit hash-object -w <file>\n")
			os.Exit(1)
		}
		git.HashObject(os.Args[3])
	case "ls-tree":
		if len(os.Args) != 4 || os.Args[2] != "--name-only" {
			fmt.Fprintf(os.Stderr, "usage: mygit ls-tree --name-only <tree-sha>\n")
			os.Exit(1)
		}
		git.LsTree(os.Args[3])
	case "write-tree":
		if len(os.Args) != 2 {
			fmt.Fprintf(os.Stderr, "usage: mygit write-tree\n")
			os.Exit(1)
		}
		git.WriteTree()
	case "commit-tree":
		if len(os.Args) != 7 || os.Args[3] != "-p" || os.Args[5] != "-m" {
			fmt.Fprintf(os.Stderr, "usage: mygit commit-tree <tree-sha> -p <commit-sha> -m <message>\n")
			os.Exit(1)
		}
		git.CommitTree(os.Args[2], os.Args[4], os.Args[6])
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
