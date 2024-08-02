package git

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"vcs/pkg"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
)

type threestrings struct {
	basecontent string
	content1    string
	content2    string
}

// Merges otherBranchName to the current branch
func Merge(otherBranchName string) {
	otherBranchHash, isBranch, _, err := pkg.FindHashofCommit(otherBranchName)
	if !isBranch {
		fmt.Println("enter a valid branch name")
		return
	}
	if err != nil {
		return
	}
	currentBranchName := filepath.Base(pkg.ParseHEADgivePath())
	currentBranchHash, _ := pkg.FetchHeadsSHAfromPath(pkg.ParseHEADgivePath())
	currentBranchTreeHash := pkg.ParseCommit(currentBranchHash).Treesha
	// fmt.Println(otherBranchHash, otherBranchName, currentBranchHash, currentBranchName, currentBranchTreeHash)
	otherBranchTreeHash := pkg.ParseCommit(otherBranchHash).Treesha
	baseHash := findBaseCommit(currentBranchName, otherBranchName)
	baseTreeHash := pkg.ParseCommit(baseHash).Treesha
	fmt.Println(baseHash, baseTreeHash)
	pairslice1 := FindDiffTrees(baseTreeHash, currentBranchTreeHash) //if entry exists in the slice, it means that the file has been changed.
	pairslice2 := FindDiffTrees(baseTreeHash, otherBranchTreeHash)   // the entries which are the same are not included in the slice
	mp := make(map[string]threestrings)
	for _, pair1 := range pairslice1 {
		fmt.Println(pair1.Content1, "here in pair1", pair1.Content2, pair1.FilePathRel1)
		mp[pair1.FilePathRel1] = threestrings{pair1.Content1, pair1.Content2, ""}
	}
	for _, pair2 := range pairslice2 {
		fmt.Println(pair2.Content1, "here", pair2.Content2, pair2.FilePathRel1)

		// Retrieve the current value from the map
		existing := mp[pair2.FilePathRel1]

		// Update the content2 field
		existing.content2 = pair2.Content2

		// Store the updated struct back into the map
		mp[pair2.FilePathRel1] = existing
	}
	fmt.Println(mp)
	var mergedString string
	for key, val := range mp {
		fmt.Println("HERE?")
		if val.content1 == "" {
			fmt.Println("here 2 for", key)
			mergedString = val.content2
		} else if val.content2 == "" {
			mergedString = val.content1
			fmt.Println("here 1 for", key)
		} else {
			mergedString = MergeStringsUserPrompt(key, val)
			fmt.Println("I WAS HEREEE!")
		}
		file, _ := os.OpenFile(filepath.Join(pkg.WorkingDirPath, key), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
		_, err := file.WriteString(mergedString)
		pkg.Check(err)
	}
}

/*
merge strings.
unchanged lines will be added as they are, if the lines are changed in one of those and not in other, the extra lines are included in the merge.
if both lines are changed, ask the user which one to include in the string.
*/
func MergeStringsUserPrompt(path string, val threestrings) string {
	var mergedString string
	edits1 := myers.ComputeEdits(span.URIFromPath(path), val.content1, val.content2)
	x := gotextdiff.ToUnified(path, path, val.content1, edits1)
	for _, hunk := range x.Hunks {
		lines := hunk.Lines
		fmt.Println("These are lines", lines)
		idx := 0
		for idx < len(lines) {
			if lines[idx].Kind == gotextdiff.Equal {
				mergedString += fmt.Sprint(lines[idx].Content)
				idx++
			} else {
				var optionString1, optionString2 string
				for idx < len(lines) && lines[idx].Kind != gotextdiff.Equal {
					if lines[idx].Kind == gotextdiff.Delete {
						optionString1 += lines[idx].Content
					} else {
						optionString2 += lines[idx].Content
					}
					idx++
				}
				fmt.Println("options start", optionString1, "mid", optionString2, "options end")
				mergedString += fmt.Sprint(userPromptChooseString(optionString1, optionString2))
			}
		}
	}
	return mergedString
}

// ask the user which one to include in the string.
func userPromptChooseString(str1 string, str2 string) string {
	fmt.Print("\n\n\n<<<<<<< one\n", str1, "\n", "=======\n", str2, ">>>>>>> two\n\n\nWrite either one, two, or custom\n")
	var op string
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	op = strings.TrimSpace(input.Text())
	if op == "one" {
		return str1
	} else if op == "two" {
		return str2
	} else if op == "custom" {
		fmt.Println("Enter the custom text that you want to be replaced by both of these :")
		input = bufio.NewScanner(os.Stdin)
		input.Scan()
		return input.Text()
	} else {
		pkg.Check(fmt.Errorf("enter a valid operation. Entered %s", op))
		return ""
	}
}

// returns the base commit of two branches
func findBaseCommit(branchname1, branchname2 string) string {
	rootToCommitSlice1 := pkg.TraverseGetRootToTip(branchname1, "")
	rootToCommitSlice2 := pkg.TraverseGetRootToTip(branchname2, "")
	j := 0
	i := 0
	n1 := len(rootToCommitSlice1)
	n2 := len(rootToCommitSlice2)
	latestCommonHash := "0000000000000000000000000000000000000000"
	for {
		if i >= n1 || j >= n2 {
			return latestCommonHash
		}
		for rootToCommitSlice1[i].Operation[:6] != "commit" {
			i++
			if i >= n1 {
				return latestCommonHash
			}
		}
		for rootToCommitSlice2[j].Operation[:6] != "commit" {
			j++
			if j >= n2 {
				return latestCommonHash
			}
		}
		if rootToCommitSlice1[i].Currentsha == rootToCommitSlice2[j].Currentsha {
			latestCommonHash = rootToCommitSlice1[i].Currentsha
			i++
			j++
		} else {
			return latestCommonHash
		}
	}
}
