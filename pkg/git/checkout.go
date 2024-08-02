package git

import (
	"fmt"
	"path/filepath"
	"vcs/pkg"
)

/*
This function is used to switch from one branch or hash to another branch or hash.
This also changes the working directory wrt the branch on which it is shifted.
*/
func Checkout(args []string) {
	sha, branchexists, _, err := pkg.FindHashofCommit(args[0])
	if err != nil {
		fmt.Println("no such commit exists as", args[0])
		return
	}
	if branchexists {
		pkg.UpdateHEAD("refs: " + filepath.Join("/refs/heads", args[0]))
		fmt.Println("On branch ", args[0])
	} else {
		pkg.UpdateHEAD(sha)
		fmt.Println("On commit ", sha, ". Entering detached HEAD state.")
	}
	if len(args) == 1 {
		pkg.RecoverWorkingDirToCommitWithDeletions(pkg.WorkingDirPath, sha)
		pkg.RecoverIndexToCommit(pkg.WorkingDirPath, sha)
	}
	for _, path := range args[1:] {
		path = filepath.Join(pkg.WorkingDirPath, path)
		pkg.RecoverWorkingDirToCommitWithDeletions(path, sha)
		pkg.RecoverIndexToCommit(path, sha)
	}
}
