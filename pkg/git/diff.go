package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"vcs/pkg"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type pairstrings struct {
	name  string
	isDir bool
}

func Diff(args []string) { // print output ko thik kaise karu

	switch len(args) {
	case 0: // compares working directory and staging area
		FindDiffStagingArea(pkg.WorkingDirPath)
	case 1:
		path := filepath.Join(pkg.WorkingDirPath, args[0])
		file, err := os.Open(path)
		if err != nil {
			fmt.Println(path + " does not exist")
		}
		info, err := file.Stat()
		pkg.Check(err)
		if info.IsDir() {
			FindDiffStagingArea(path)
		} else {
			FindDiffFileStagingArea(path)
		}
	case 2:
		hash1, _, _, err := pkg.FindHashofCommit(args[0])
		pkg.Check(err)
		hash2, _, _, err := pkg.FindHashofCommit(args[1])
		pkg.Check(err)
		// fmt.Println(hash1, hash2)
		commit1 := pkg.ParseCommit(hash1)
		commit2 := pkg.ParseCommit(hash2)
		// fmt.Println(commit1.Treesha)
		// fmt.Println(commit2.Treesha)

		set1, set2 := FindDiffTrees(commit1.Treesha, commit2.Treesha)
		for i, _ := range set1 {
			diff := FindDiffStrings(set1[i], set2[i])
			fmt.Println(diff)
		}
	case 3:
		hash1, _, _, err := pkg.FindHashofCommit(args[0])
		pkg.Check(err)
		hash2, _, _, err := pkg.FindHashofCommit(args[1])
		pkg.Check(err)
		s1, s2 := FindDiffCommittedFile(args[2], hash1, hash2)
		diff := FindDiffStrings(s1, s2)
		fmt.Println(diff)
	}
}
func FindDiffStrings(a string, b string) string {
	dmp := diffmatchpatch.New()
	fileA, fileB, dmpStrings := dmp.DiffLinesToChars(a, b)
	diffs := dmp.DiffMain(fileA, fileB, false)
	diffs = dmp.DiffCharsToLines(diffs, dmpStrings)
	diffs = dmp.DiffCleanupSemantic(diffs)

	return dmp.DiffPrettyText(diffs)
}

func FindDiffCommittedFile(path string, commithash1 string, commithash2 string) (string, string) { // not used traverse tree here bcoz getCommittedEntry seems faster.
	treeentry1 := pkg.GetCommittedEntry(filepath.Join(pkg.WorkingDirPath, path), commithash1)
	treeentry2 := pkg.GetCommittedEntry(filepath.Join(pkg.WorkingDirPath, path), commithash2)
	var s1, s2 string
	if treeentry1.FileType == "tree" && treeentry2.FileType == "tree" {
		FindDiffTrees(treeentry1.Sha, treeentry2.Sha)
		return s1, s2
	}
	if treeentry1.FileType != treeentry2.FileType {
		fmt.Println(path, "is of different filetypes in both commit hashes")
		return s1, s2
	}
	if treeentry1.Sha == "" {
		s1 = ""
		// fmt.Println("gere")
	} else {
		s1, _, _ = pkg.ReadCompressedFile(filepath.Join(pkg.VCSDirPath, "objects", treeentry1.Sha[:2], treeentry1.Sha[2:]))
	}

	if treeentry2.Sha == "" {
		s2 = ""
		// fmt.Println("here")
	} else {
		s2, _, _ = pkg.ReadCompressedFile(filepath.Join(pkg.VCSDirPath, "objects", treeentry2.Sha[:2], treeentry2.Sha[2:]))
	}
	// fmt.Println(s1)
	// fmt.Println(s2)
	return s1, s2
}

func FindDiffTrees(treesha1 string, treesha2 string) ([]string, []string) {
	tree1 := pkg.TraverseTree(treesha1)
	tree2 := pkg.TraverseTree(treesha2)
	mp := make(map[pairstrings]pkg.PairOfSHA)
	for _, entry := range tree1 {
		mp[pairstrings{entry.Path, entry.FileType == "tree"}] = pkg.PairOfSHA{P1: pkg.Pair{Exists: true, Sha: entry.Sha}}
	}
	// fmt.Println(mp)
	for _, entry := range tree2 {
		mp[pairstrings{entry.Path, entry.FileType == "tree"}] = pkg.PairOfSHA{P1: mp[pairstrings{entry.Path, entry.FileType == "tree"}].P1, P2: pkg.Pair{Exists: true, Sha: entry.Sha}}
	}
	var set1, set2 []string
	for key, pop := range mp {
		// fmt.Println(key)
		var s1 string
		var s2 string
		if pop.P1.Sha == pop.P2.Sha {
			continue
		}
		if !pop.P1.Exists {
			// fmt.Println("s1 empty for ", key)
			s1 = ""
		}
		if !pop.P2.Exists {
			// fmt.Println("s2 empty for ", key)
			s2 = ""
		}
		if !key.isDir {
			if pop.P1.Exists {
				s1, _, _ = pkg.ReadCompressedFile(filepath.Join(pkg.VCSDirPath, "objects", pop.P1.Sha[:2], pop.P1.Sha[2:]))
			}
			if pop.P2.Exists {
				s2, _, _ = pkg.ReadCompressedFile(filepath.Join(pkg.VCSDirPath, "objects", pop.P2.Sha[:2], pop.P2.Sha[2:]))
			}
			set1 = append(set1, s1)
			set2 = append(set2, s2)
		}
	}
	return set1, set2
}

func FindDiffStagingArea(path string) {
	index := pkg.ParseIndex()
	for key := range index.Entries {
		if !strings.HasPrefix(key, path) { // checked for folders
			continue
		}
		FindDiffFileStagingArea(key)
	}
}

func FindDiffFileStagingArea(path string) {
	index = pkg.ParseIndex()
	value, exists := index.Entries[path]
	if !exists {
		fmt.Println(path + " does not exist in staging area")
		return
	}
	file, err := os.Open(path)
	var sha string
	if err == nil {
		sha = pkg.GetSHAofFile(file)
	}
	if sha != value.Id {
		s1, _, _ := pkg.ReadCompressedFile(filepath.Join(pkg.VCSDirPath, "objects", value.Id[:2], value.Id[2:]))
		data, _ := os.ReadFile(path)
		diff := FindDiffStrings(s1, string(data))
		fmt.Println(diff)
	}
}
