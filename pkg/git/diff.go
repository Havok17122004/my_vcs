package git

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"vcs/pkg"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type pair struct {
	exists bool
	sha    string
}

type pairofpairs struct {
	p1 pair
	p2 pair
}

type pairstrings struct {
	name  string
	isDir bool
}

func Diff(args []string) {

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
		hash1, _, err := pkg.FindHash(args[0])
		pkg.Check(err)
		hash2, _, err := pkg.FindHash(args[1])
		pkg.Check(err)
		// fmt.Println(hash1, hash2)
		commit1 := ParseCommit(hash1)
		commit2 := ParseCommit(hash2)
		// fmt.Println(commit1.Treesha)
		// fmt.Println(commit2.Treesha)

		FindDiffTrees(commit1.Treesha, commit2.Treesha)
	case 3:
		hash1, _, err := pkg.FindHash(args[0])
		pkg.Check(err)
		hash2, _, err := pkg.FindHash(args[1])
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

func FindDiffCommittedFile(path string, commithash1 string, commithash2 string) {
	treeentry1 := GetCommittedEntry(filepath.Join(pkg.WorkingDirPath, path), commithash1)
	treeentry2 := GetCommittedEntry(filepath.Join(pkg.WorkingDirPath, path), commithash2)
	var s1, s2 string
	if treeentry1.Sha == "" {
		s1 = ""
		// fmt.Println("gere")
	} else {
		s1 = pkg.ReadCompressedFile(filepath.Join(pkg.VCSDirPath, "objects", treeentry1.Sha[:2], treeentry1.Sha[2:]))
	}

	if treeentry2.Sha == "" {
		s2 = ""
		// fmt.Println("here")
	} else {
		s2 = pkg.ReadCompressedFile(filepath.Join(pkg.VCSDirPath, "objects", treeentry2.Sha[:2], treeentry2.Sha[2:]))
	}
	// fmt.Println(s1)
	// fmt.Println(s2)
	FindDiffStrings(s1, s2)
}

func FindDiffTrees(treesha1 string, treesha2 string) {
	tree1 := ParseTree(treesha1)
	tree2 := ParseTree(treesha2)
	// var mp map [string] pairofpairs does not initialise
	// fmt.Println(tree1)
	// fmt.Println(tree2)
	mp := make(map[pairstrings]pairofpairs)
	for _, entry := range *tree1 {

		mp[pairstrings{entry.Name, entry.FileType == "tree"}] = pairofpairs{p1: pair{true, entry.Sha}}
	}
	// fmt.Println(mp)
	for _, entry := range *tree2 {
		mp[pairstrings{entry.Name, entry.FileType == "tree"}] = pairofpairs{p1: mp[pairstrings{entry.Name, entry.FileType == "tree"}].p1, p2: pair{true, entry.Sha}}
	}
	// fmt.Println(mp)
	for key, pop := range mp {
		// fmt.Println(key)
		var s1 string
		var s2 string
		if pop.p1.sha == pop.p2.sha {
			continue
		}
		if !pop.p1.exists {
			// fmt.Println("s1 empty for ", key)
			s1 = ""
		}
		if !pop.p2.exists {
			// fmt.Println("s2 empty for ", key)
			s2 = ""
		}
		if key.isDir {
			FindDiffTrees(pop.p1.sha, pop.p2.sha)
		} else {
			if pop.p1.exists {
				s1 = pkg.ReadCompressedFile(filepath.Join(pkg.VCSDirPath, "objects", pop.p1.sha[:2], pop.p1.sha[2:]))
			}
			if pop.p2.exists {
				s2 = pkg.ReadCompressedFile(filepath.Join(pkg.VCSDirPath, "objects", pop.p2.sha[:2], pop.p2.sha[2:]))
			}
			FindDiffStrings(s1, s2)
		}
	}
}

func FindDiffStagingArea(path string) {
	index := pkg.ParseIndex()
	for key, entry := range index.Entries {
		if !strings.HasPrefix(key, path) {
			continue
		}
		file, err := os.Open(key)
		var sha []byte
		if err == nil {
			sha = pkg.GetSHAofFile(file)
		}
		if string(sha) != entry.Id {
			s1 := pkg.ReadCompressedFile(filepath.Join(pkg.VCSDirPath, "objects", entry.Id[:2], entry.Id[2:]))
			data, err := os.ReadFile(key)
			if err != nil && err != io.EOF {
				pkg.Check(err)
			}
			FindDiffStrings(s1, string(data))
		}
	}
}

func FindDiffFileStagingArea(path string) {
	index = pkg.ParseIndex()
	value, exists := index.Entries[path]
	if !exists {
		fmt.Println(path + " does not exist in staging area")
		return
	}
	s1 := pkg.ReadCompressedFile(filepath.Join(pkg.VCSDirPath, "objects", value.Id[:2], value.Id[2:]))
	data, err := os.ReadFile(path)
	if err != nil && err != io.EOF {
		pkg.Check(err)
	}
	FindDiffStrings(s1, string(data))
}
