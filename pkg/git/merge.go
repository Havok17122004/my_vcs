package git

import (
	"fmt"
	"path/filepath"
	"vcs/pkg"
)

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
	FindDiffTrees(baseTreeHash, currentBranchTreeHash)
	FindDiffTrees(baseTreeHash, otherBranchTreeHash)
}

func findBaseCommit(branchname1, branchname2 string) string {
	latestLogContents1 := *pkg.ParseLog(filepath.Join(pkg.VCSDirPath, "logs/refs/heads", branchname1+".txt"))
	latestLogContents2 := *pkg.ParseLog(filepath.Join(pkg.VCSDirPath, "logs/refs/heads", branchname2+".txt"))
	logContents1 := latestLogContents1
	logContents2 := latestLogContents2
	// fmt.Println(logContents1, logContents2)
	for {
		for _, content1 := range logContents1 {
			if logContents2[0].Currentsha == content1.Currentsha {
				return content1.Currentsha
			}
		}
		if logContents1[0].Operation == "branch" {
			logContents1 = *pkg.ParseLog(filepath.Join(pkg.VCSDirPath, "logs/refs/heads", logContents1[0].Parentsha))
		} else if logContents1[0].Operation == "commit(initial)" { //may run infinitely if log file is not correctly read or written
			break
		}
	}
	logContents1 = latestLogContents1
	for {
		for _, content2 := range logContents2 {
			if logContents1[0].Currentsha == content2.Currentsha {
				return content2.Currentsha
			}
		}
		if logContents2[0].Operation == "branch" {
			logContents2 = *pkg.ParseLog(filepath.Join(pkg.VCSDirPath, "logs/refs/heads", logContents2[0].Parentsha))
		} else if logContents2[0].Operation == "commit(initial)" { //may run infinitely if log file is not correctly read or written
			break
		}
	}
	return ""
}
