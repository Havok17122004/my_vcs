package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"vcs/pkg"
)

func Checkout(branch string) {
	_, err := os.Open(filepath.Join(pkg.VCSDirPath, "/refs/heads/", branch+".txt"))
	var sha string
	if err == nil {
		pkg.UpdateHEAD("refs: " + filepath.Join("/refs/heads", branch))
		fmt.Println("On branch ", branch)
	} else {
		logcontents := pkg.ParseLog(filepath.Join(pkg.VCSDirPath, "/logs/HEAD.txt"))
		if logcontents == nil {
			fmt.Printf("%s does not match any branchname or commit hash\n", branch)
			return
		}
		cnt := 0
		// fmt.Print(logcontents)
		for _, logline := range *logcontents {
			// fmt.Println(logline.Currentsha, branch)
			if strings.HasPrefix(logline.Currentsha, branch) {
				cnt++
				sha = logline.Currentsha
			}
		}
		if cnt != 1 {
			// fmt.Print(cnt)
			fmt.Printf("%s does not match any branchname or commit hash\n", branch)
		} else {
			pkg.UpdateHEAD(sha)
			fmt.Println("On commit ", sha, ". Entering detached HEAD state.")
		}

	}
}
