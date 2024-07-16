package git

import (
	"bytes"
	"fmt"
	"net/url"
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

		FindDiffTrees(commit1.Treesha, commit2.Treesha)
	case 3:
		hash1, _, _, err := pkg.FindHashofCommit(args[0])
		pkg.Check(err)
		hash2, _, _, err := pkg.FindHashofCommit(args[1])
		pkg.Check(err)
		FindDiffCommittedFile(args[2], hash1, hash2)
	}
}

// DiffPrettyText converts a []Diff into a colored text report.
func DiffPrettyText(dmp *diffmatchpatch.DiffMatchPatch, diffs []diffmatchpatch.Diff) string {
	var buff bytes.Buffer
	for _, diff := range diffs {
		var str string
		// changed := false
		text := diff.Text

		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			// Green background for insertions
			// _, _ = buff.WriteString("\x1b[42m")
			// _, _ = buff.WriteString(text)
			// _, _ = buff.WriteString("\x1b[49m")
			str = fmt.Sprintf("%s\x1b[42m%s\x1b[49m", str, text)
			// changed = true
		case diffmatchpatch.DiffDelete:
			// Red background for deletions
			// _, _ = buff.WriteString("\x1b[41m")
			// _, _ = buff.WriteString(text)
			// _, _ = buff.WriteString("\x1b[49m")
			str = fmt.Sprintf("%s\x1b[41m%s\x1b[49m", str, text)
			// changed = true
		case diffmatchpatch.DiffEqual:
			// Default background for equalities
			// _, _ = buff.WriteString("\x1b[49m")
			// _, _ = buff.WriteString(text)
			// str = fmt.Sprintf("%s\x1b[49m%s", str, text)
		}
		// if changed {
		buff.WriteString(str)
		// }
	}

	return buff.String()
}

func FindDiffStrings(text1 string, text2 string) {
	// fmt.Println(text1)
	// fmt.Println(text2, "here")
	// dmp := diffmatchpatch.New()

	// diffs := dmp.DiffMain(text1, text2, false) // if bool is true, calculates diff line-by-line and if false it calculates diff by divide and conquer method.
	// fmt.Println(diffs)
	// fmt.Println(dmp.DiffPrettyText(diffs))

	// dmp := diffmatchpatch.New()

	// fileAdmp, fileBdmp, dmpStrings := dmp.DiffLinesToChars(text1, text2)
	// diffs := dmp.DiffMain(fileAdmp, fileBdmp, false)
	// diffs = dmp.DiffCharsToLines(diffs, dmpStrings)
	// diffs = dmp.DiffCleanupSemantic(diffs)
	// // s := dmp.DiffPrettyText(diffs)
	// fmt.Println(DiffPrettyText(dmp, diffs))
	// fmt.Print(s)
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(text1, text2, false)
	patches := dmp.PatchMake(text1, diffs)
	patchText := dmp.PatchToText(patches)
	fmt.Println(cleanPatchText(patchText))
}

func cleanPatchText(patchText string) string {
	// Decode URL-encoded sequences
	decodedText, err := url.QueryUnescape(patchText)
	if err != nil {
		fmt.Println("Error decoding patch text:", err)
		return patchText
	}

	// Remove extraneous newlines
	var buff bytes.Buffer
	lines := bytes.Split([]byte(decodedText), []byte("\n"))
	for _, line := range lines {
		if len(line) > 0 {
			buff.Write(line)
			buff.WriteByte('\n')
		}
	}
	// var buff bytes.Buffer
	buff.Write([]byte(decodedText))
	return buff.String()
}

func FindDiffCommittedFile(path string, commithash1 string, commithash2 string) { // not used traverse tree here bcoz getCommittedEntry seems faster.
	treeentry1 := pkg.GetCommittedEntry(filepath.Join(pkg.WorkingDirPath, path), commithash1)
	treeentry2 := pkg.GetCommittedEntry(filepath.Join(pkg.WorkingDirPath, path), commithash2)
	if treeentry1.FileType == "tree" && treeentry2.FileType == "tree" {
		FindDiffTrees(treeentry1.Sha, treeentry2.Sha)
		return
	}
	if treeentry1.FileType != treeentry2.FileType {
		fmt.Println(path, "is of different filetypes in both commit hashes")
		return
	}
	var s1, s2 string
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
	FindDiffStrings(s1, s2)
}

func FindDiffTrees(treesha1 string, treesha2 string) {
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
			FindDiffStrings(s1, s2)
		}
	}
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
		FindDiffStrings(s1, string(data))
	}
}
