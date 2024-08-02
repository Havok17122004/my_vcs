package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"vcs/pkg"
)

const (
	Default = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
)

// gives the status of the current working directory and the staging area
func Status() {
	branch_or_hash := pkg.ParseHEADgivePath()
	if strings.Contains(branch_or_hash, string(os.PathSeparator)) {
		fmt.Println("On branch", filepath.Base(branch_or_hash))
	} else {
		fmt.Println("On commit", branch_or_hash)
	}
	toBeCommitted := CompareIndexHEAD()
	isEmpty := true
	for path, entry := range toBeCommitted {
		if entry.Status == "unchanged" {
			continue
		}
		if isEmpty {
			isEmpty = false
			fmt.Printf("Changes to be committed:\n%s", Green)
		}
		relPath, err := filepath.Rel(pkg.WorkingDirPath, path)
		pkg.Check(err)
		fmt.Printf("\t%s:\t%s\n", entry.Status, relPath)
	}
	if !isEmpty {
		fmt.Printf("%s\n", Default)
	}

	notToBeCommitted := CompareIndexWorkingDir()
	var notTracked []string
	isEmpty = true
	for path, entry := range notToBeCommitted {
		if entry.Status == "unchanged" {
			continue
		}
		if entry.Status == "untracked" {
			notTracked = append(notTracked, path)
			continue
		}
		if isEmpty {
			isEmpty = false
			fmt.Printf("Changes not staged for commit:\n%s", Red)
		}
		relPath, err := filepath.Rel(pkg.WorkingDirPath, path)
		pkg.Check(err)
		fmt.Printf("\t%s:\t%s\n", entry.Status, relPath)
	}
	if !isEmpty {
		fmt.Printf("%s\n", Default)
	}

	if len(notTracked) != 0 {
		fmt.Printf("Untracked files:\n%s", Red)
		for _, path := range notTracked {
			relPath, err := filepath.Rel(pkg.WorkingDirPath, path)
			pkg.Check(err)
			fmt.Printf("\t%s\n", relPath)
		}
		fmt.Printf("%s\n", Default)
	}
}

// changes to be committed
func CompareIndexHEAD() map[string]pkg.PairSHAandStatus {
	mp := make(map[string]pkg.PairSHAandStatus)
	index := pkg.ParseIndex()
	fmt.Println(index)
	for path, props := range index.Entries {
		mp[path] = pkg.PairSHAandStatus{
			P:      pkg.PairOfSHA{P1: pkg.Pair{Exists: true, Sha: props.Id}},
			Status: "",
		}
	}
	commithash, err := pkg.FetchHeadsSHAfromPath(pkg.ParseHEADgivePath())
	if err == nil { // if err!=nil, there wont be any HEAD. So, no additions to mp
		// var treeEntries []pkg.TreeEntry
		treeEntries := pkg.TraverseTree(pkg.ParseCommit(commithash).Treesha)

		for _, entry := range treeEntries {
			if entry.FileType == "tree" {
				continue
			}
			mp[entry.Path] = pkg.PairSHAandStatus{
				P: pkg.PairOfSHA{
					P1: mp[entry.Path].P.P1,
					P2: pkg.Pair{Exists: true, Sha: entry.Sha},
				},
				Status: "",
			}
		}
	}
	for path, pair := range mp { // does changing pair.Status make respective changes? if not, try to pass it as ref.
		if pair.P.P1.Exists { // exist in index
			if pair.P.P2.Exists {
				if pair.P.P1.Sha == pair.P.P2.Sha { // same file, same hash
					pair.Status = "unchanged"
				} else { // same file exist but sha different
					pair.Status = "modified"
					fmt.Println(pair.P.P1.Sha, pair.P.P2.Sha)
					// modified
				}
			} else { // exist in index, not in HEAD
				pair.Status = "new file"
				// deleted
			}
		} else { // does not exist in index, but exists in HEAD
			pair.Status = "deleted"
			// untracked
		}
		mp[path] = pair //because pair is a copy. changes made wont be reflected
	}
	return mp
}

// untracked files - don't exist in index but exist in working dir
// changes not staged for commit - exist in index but their hash is different from that in working dir
func CompareIndexWorkingDir() map[string]pkg.PairSHAandStatus {
	mp := make(map[string]pkg.PairSHAandStatus)
	index := pkg.ParseIndex()
	for path, props := range index.Entries {
		mp[path] = pkg.PairSHAandStatus{
			P:      pkg.PairOfSHA{P1: pkg.Pair{Exists: true, Sha: props.Id}},
			Status: "",
		}
	}
	files := pkg.TraverseDir(pkg.WorkingDirPath)
	for _, path := range files {
		file, err := os.Open(path)
		pkg.Check(err)
		mp[path] = pkg.PairSHAandStatus{
			P: pkg.PairOfSHA{
				P1: mp[path].P.P1,
				P2: pkg.Pair{Exists: true, Sha: pkg.GetSHAofFile(file)},
			},
			Status: "",
		}
	}
	for path, pair := range mp {
		if pair.P.P1.Exists { // exist in index
			if pair.P.P2.Exists {
				if pair.P.P1.Sha == pair.P.P2.Sha { // same file, same hash
					pair.Status = "unchanged"
				} else { // same file exist but sha different
					pair.Status = "modified"
					// modified
				}
			} else { // exist in index, not in working dir
				pair.Status = "deleted"
				// deleted
			}
		} else { // does not exist in index, but exists in working dir
			pair.Status = "untracked"
			// untracked
		}
		mp[path] = pair
	}
	return mp
}
