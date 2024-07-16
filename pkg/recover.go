package pkg // change here

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type PairOfBools struct {
	ExistInTree bool
	TreeSha     string
	ExistInDir  bool
}

func GetCommittedEntry(fullPath string, commithash string) TreeEntry {
	commitObject := ParseCommit(commithash)
	if fullPath == WorkingDirPath {
		return TreeEntry{"100755", "tree", commitObject.Treesha, fullPath}
	}
	treeObject := ParseTree(commitObject.Treesha)
	fmt.Println(fullPath)
	relPath, err := filepath.Rel(WorkingDirPath, fullPath)
	Check(err)
	structure := strings.Split(relPath, string(filepath.Separator))
	fmt.Println(structure, len(structure))
	var target TreeEntry
	for i := 0; i < len(structure); i++ {
		found := false
		// fmt.Println(structure[i])
		fmt.Println(structure[i], treeObject)
		for _, ele := range *treeObject {
			ele.Path, err = filepath.Rel(WorkingDirPath, ele.Path)
			Check(err)
			fmt.Println(ele.Path)
			if filepath.Base(ele.Path) == structure[i] {
				// fmt.Println(ele.Path, "in")
				if i == len(structure)-1 {
					found = true
					target = ele
				}
				if ele.FileType == "tree" {
					target = ele
					treeObject = ParseTree(ele.Sha)
					found = true
					// fmt.Println(treeObject)
					break
				}

			}
		}
		if !found {
			fmt.Println(structure[i], "not found")
			return TreeEntry{FileMode: "", FileType: "", Sha: "", Path: ""}
		}
	}
	// fmt.Println(target)
	return target
}

func RecoverWorkingDirToCommitWithDeletions(fullPath string, hash string) {
	treeContents := TraverseTree(ParseCommit(hash).Treesha)
	dirContents := TraverseDir(WorkingDirPath)
	mp := make(map[string]PairOfBools)
	var dirs []string
	for _, content := range treeContents {
		if content.FileType == "tree" {
			dirs = append(dirs, content.Path)
		} else {
			mp[content.Path] = PairOfBools{ExistInTree: true, TreeSha: content.Sha}
		}
	}
	for _, content := range dirContents {
		mp[content] = PairOfBools{ExistInTree: mp[content].ExistInTree, TreeSha: mp[content].TreeSha, ExistInDir: true}
	}

	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0777)
		Check(err)
	}
	for path, existStruct := range mp {
		if existStruct.ExistInTree {
			file, err := os.Open(path)
			if err == nil && GetSHAofFile(file) == existStruct.TreeSha {
				fmt.Println("unmodified", path)
				continue
			}
			fmt.Println("modified", path)
			newFile, err := os.Create(path) // if file does not exists, it is created. else it is cleared.
			Check(err)
			s, _, _ := ReadCompressedFile(filepath.Join(VCSDirPath, "objects", existStruct.TreeSha[:2], existStruct.TreeSha[2:]))
			// write s string into newfile
			_, err = newFile.WriteString(s)
			Check(err)
		} else { //delete from working dir
			// os.RemoveAll(path)
			fmt.Println("deleted", path)
		}
	}
}

func RecoverIndexToCommit(fullPath string, hash string) {
	var newIndex Index
	newIndex.Entries = make(map[string]Entry)
	prevIndex := ParseIndex()
	treeContents := TraverseTree(ParseCommit(hash).Treesha)
	for _, content := range treeContents {
		if content.FileType == "tree" {
			continue
		}
		if !strings.HasPrefix(content.Path, fullPath) {
			fmt.Println(content.Path)
			fmt.Println(prevIndex.Entries[content.Path].Metadata)
			fmt.Println(prevIndex.Entries[content.Path].Id)
			newIndex.ModifyIndex(content.Path, prevIndex.Entries[content.Path].Metadata, prevIndex.Entries[content.Path].Id)
			fmt.Println("heree")
			continue
		}
		_, _, size := ReadCompressedFile(filepath.Join(VCSDirPath, "objects", content.Sha[:2], content.Sha[2:]))
		mode, err := strconv.Atoi(content.FileMode)
		Check(err)
		fmt.Println(content.Sha, "heree")
		newIndex.ModifyIndex(content.Path, Metadata{Mode: uint32(mode), Size: size}, content.Sha)
	}
	newIndex.SaveIndex()
}
