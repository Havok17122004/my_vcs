package git

import (
	"fmt"
	"path/filepath"
	"strings"
	"vcs/pkg"
)

func GetCommittedEntry(path string, commithash string) TreeEntry {
	commitObject := ParseCommit(commithash)
	treeObject := ParseTree(commitObject.Treesha)
	fmt.Println(path)
	path, err := filepath.Rel(pkg.WorkingDirPath, path)
	pkg.Check(err)
	structure := strings.Split(path, string(filepath.Separator))
	// fmt.Println(structure, len(structure))
	var target TreeEntry
	for i := 0; i < len(structure); i++ {
		found := false
		// fmt.Println(structure[i])
		fmt.Println(structure[i], treeObject)
		for _, ele := range *treeObject {
			if ele.Name == structure[i] {
				fmt.Println(ele.Name, "in")
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
			return TreeEntry{"", "", "", ""}
		}
	}
	// fmt.Println(target)
	return target
}
