package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Traverse the dir and give a slice of filepaths that are reachable through the dir
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

// Traverse the tree and give a slice of treeEntries that are reachable through the tree
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

// Traverse the branch commits. Get a slice of logcontents from root to the tip of the branch
func TraverseGetRootToTip(branchname string, tillHash string) []LogContents {
	totalLogContents := *ParseLog(filepath.Join(VCSDirPath, "logs/refs/heads", branchname+".txt"))
	// fmt.Println(totalLogContents)
	var logContents []LogContents
	for _, content := range totalLogContents {
		if content.Currentsha == tillHash {
			logContents = append(logContents, content)
			break
		} else {
			logContents = append(logContents, content)
		}
	}
	// fmt.Println(logContents)
	if len(logContents) == 0 {
		return []LogContents{}
	}
	if logContents[0].Operation == "commit(initial)" {
		return logContents
	}
	if logContents[0].Operation != "branch" {
		Check(fmt.Errorf("log content of %s contains %v as the first line. no origin found", branchname, logContents[0]))
	}
	return append(TraverseGetRootToTip(strings.Split(logContents[0].Message, " ")[2], logContents[0].Currentsha), logContents...)
}
