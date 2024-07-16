package pkg

import (
	"os"
	"path/filepath"
	"strings"
)

func TraverseDir(path string) []string {
	var files []string
	direntries, err := os.ReadDir(path)
	Check(err)
	for _, entry := range direntries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		if entry.IsDir() {
			files = append(files, TraverseDir(filepath.Join(path, entry.Name()))...)
			continue
		}
		files = append(files, filepath.Join(path, entry.Name()))
	}
	return files
}

func TraverseTree(treeHash string) []TreeEntry { // includes the "tree" marked objects in entries too.
	if treeHash == "" {
		return []TreeEntry{}
	}
	var entries []TreeEntry
	tree := ParseTree(treeHash)
	for _, entry := range *tree {
		entries = append(entries, entry)
		if entry.FileType == "tree" {
			entries = append(entries, TraverseTree(entry.Sha)...)
		}
	}
	return entries
}
