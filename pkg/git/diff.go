package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"vcs/pkg"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
)

type pairstrings struct {
	path  string
	isDir bool
}

type PairFilesWithContents struct {
	Content1     string
	FilePathRel1 string
	Content2     string
	FilePathRel2 string
}

func Diff(args []string) { // does not consider file rename and file move!!!

	switch len(args) {
	case 0: // compares working directory and staging area
		allContents := FindDiffStagingArea(pkg.WorkingDirPath)
		for _, content := range allContents {
			fmt.Println(FindDiffStrings(content))
		}
	case 1:
		path := filepath.Join(pkg.WorkingDirPath, args[0])
		file, err := os.Open(path)
		if err != nil {
			fmt.Println(path + " does not exist")
		}
		info, err := file.Stat()
		pkg.Check(err)
		if info.IsDir() {
			allContents := FindDiffStagingArea(path)
			for _, content := range allContents {
				fmt.Println(FindDiffStrings(content))
			}
		} else {
			content := FindDiffFileStagingArea(path)
			fmt.Println(FindDiffStrings(content))
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

		allContents := FindDiffTrees(commit1.Treesha, commit2.Treesha)
		for _, content := range allContents {
			fmt.Println(FindDiffStrings(content))
		}
	case 3:
		hash1, _, _, err := pkg.FindHashofCommit(args[0])
		pkg.Check(err)
		hash2, _, _, err := pkg.FindHashofCommit(args[1])
		pkg.Check(err)
		content := FindDiffCommittedFile(args[2], hash1, hash2)
		fmt.Println(FindDiffStrings(content))
	}
}

// find the difference between two strings and returns a string containing the diffs
func FindDiffStrings(content PairFilesWithContents) string {
	edits := myers.ComputeEdits(span.URIFromPath(content.FilePathRel1), content.Content1, content.Content2)
	allLines := fmt.Sprint(gotextdiff.ToUnified(content.FilePathRel1, content.FilePathRel2, content.Content1, edits))
	lines := strings.Split(allLines, "\n")
	var mainText string
	for _, line := range lines {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++ ") {
			mainText = mainText + Green + line + Default + "\n"
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			mainText = mainText + Red + line + Default + "\n"
		} else {
			mainText = mainText + line + "\n"
		}
	}
	return mainText
}

/*
	find the difference between two versions of files. version1 being file version in commithash1. version2 being file version in commithash2.

returns a pair of file versions, each value containing both filepath and filecontents.
*/
func FindDiffCommittedFile(path string, commithash1 string, commithash2 string) PairFilesWithContents { // not used traverse tree here bcoz getCommittedEntry seems faster.
	treeentry1 := pkg.GetCommittedEntry(filepath.Join(pkg.WorkingDirPath, path), commithash1)
	treeentry2 := pkg.GetCommittedEntry(filepath.Join(pkg.WorkingDirPath, path), commithash2)
	var content PairFilesWithContents
	content.FilePathRel1, _ = filepath.Rel(pkg.WorkingDirPath, treeentry1.Path)
	content.FilePathRel2, _ = filepath.Rel(pkg.WorkingDirPath, treeentry2.Path)
	if treeentry1.FileType == "tree" && treeentry2.FileType == "tree" {
		FindDiffTrees(treeentry1.Sha, treeentry2.Sha)
		return content
	}
	if treeentry1.FileType != treeentry2.FileType {
		fmt.Println(path, "is of different filetypes in both commit hashes")
		return content
	}
	if treeentry1.Sha == "" {
		// content.content1 = "" not needed
		// fmt.Println("gere")
	} else {
		content.Content1, _, _ = pkg.ReadCompressedFile(filepath.Join(pkg.VCSDirPath, "objects", treeentry1.Sha[:2], treeentry1.Sha[2:]))
	}

	if treeentry2.Sha == "" {
		// content.content2 = "" not needed
		// fmt.Println("here")
	} else {
		content.Content2, _, _ = pkg.ReadCompressedFile(filepath.Join(pkg.VCSDirPath, "objects", treeentry2.Sha[:2], treeentry2.Sha[2:]))
	}
	// fmt.Println(content.content1)
	// fmt.Println(content.content2)
	return content
}

/*
	find the difference between two versions of all files. version1 being file versions in commithash1. version2 being file versions in commithash2.

returns a slice of pair of file versions, each value containing both filepath and filecontents.
*/
func FindDiffTrees(treesha1 string, treesha2 string) []PairFilesWithContents {
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
	var allContents []PairFilesWithContents
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
			relPath, _ := filepath.Rel(pkg.WorkingDirPath, key.path)
			allContents = append(allContents, PairFilesWithContents{Content1: s1, Content2: s2, FilePathRel1: relPath, FilePathRel2: relPath})
		}
	}
	return allContents
}

/*
	find the difference between the working directory and staging area of all files.

returns a slice of pair of file versions, each value containing both filepath and filecontents.
*/
func FindDiffStagingArea(path string) []PairFilesWithContents {
	var allContents []PairFilesWithContents
	index := pkg.ParseIndex()
	for key := range index.Entries {
		if !strings.HasPrefix(key, path) { // checked for folders
			continue
		}
		allContents = append(allContents, FindDiffFileStagingArea(key))
	}
	return allContents
}

/*
	find the difference between the working directory and staging area of the file specified.

returns a pair of file versions, each value containing both filepath and filecontents.
*/
func FindDiffFileStagingArea(path string) PairFilesWithContents {
	index = pkg.ParseIndex()
	value, exists := index.Entries[path]
	if !exists {
		fmt.Println(path + " does not exist in staging area")
		return PairFilesWithContents{}
	}
	file, err := os.Open(path)
	var sha string
	if err == nil {
		sha = pkg.GetSHAofFile(file)
	}
	if sha != value.Id {
		s1, _, _ := pkg.ReadCompressedFile(filepath.Join(pkg.VCSDirPath, "objects", value.Id[:2], value.Id[2:]))
		data, _ := os.ReadFile(path)
		relPath, _ := filepath.Rel(pkg.WorkingDirPath, path)
		return PairFilesWithContents{Content1: s1, Content2: string(data), FilePathRel1: relPath, FilePathRel2: relPath}
	}
	return PairFilesWithContents{}
}
