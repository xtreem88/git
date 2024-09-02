package git

import (
	"fmt"
	"os"
	"time"
)

func CommitTree(treeSHA, parentSHA, message string) {
	repo := NewRepo(".")

	// Hardcoded author and committer information
	author := "John Doe <john@example.com>"
	committer := "John Doe <john@example.com>"
	timestamp := time.Now().Unix()
	timezone := "-0700" // Hardcoded timezone offset

	commitContent := fmt.Sprintf("tree %s\n", treeSHA)
	commitContent += fmt.Sprintf("parent %s\n", parentSHA)
	commitContent += fmt.Sprintf("author %s %d %s\n", author, timestamp, timezone)
	commitContent += fmt.Sprintf("committer %s %d %s\n", committer, timestamp, timezone)
	commitContent += "\n" + message + "\n"

	commitSHA, err := repo.WriteObject("commit", []byte(commitContent))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing commit object: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(commitSHA)
}
