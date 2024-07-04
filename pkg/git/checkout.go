package git

import (
	"fmt"
	"path/filepath"
	"vcs/pkg"
)

func Checkout(text string) {
	sha, branchexists, err := pkg.FindHash(text)
	pkg.Check(err)
	if branchexists {
		pkg.UpdateHEAD("refs: " + filepath.Join("/refs/heads", text))
		fmt.Println("On branch ", text)
	} else {
		pkg.UpdateHEAD(sha)
		fmt.Println("On commit ", sha, ". Entering detached HEAD state.")
	}
	// _, err := os.Open(filepath.Join(pkg.VCSDirPath, "/refs/heads/", branch+".txt"))
	// var sha string
	// if err == nil {
	// 	pkg.UpdateHEAD("refs: " + filepath.Join("/refs/heads", branch))
	// 	fmt.Println("On branch ", branch)
	// } else {
	// 	logcontents := pkg.ParseLog(filepath.Join(pkg.VCSDirPath, "/logs/HEAD.txt"))
	// 	if logcontents == nil {
	// 		fmt.Printf("%s does not match any branchname or commit hash\n", branch)
	// 		return
	// 	}
	// 	cnt := 0
	// 	// fmt.Print(logcontents)
	// 	for _, logline := range *logcontents {
	// 		// fmt.Println(logline.Currentsha, branch)
	// 		if strings.HasPrefix(logline.Currentsha, branch) {
	// 			cnt++
	// 			sha = logline.Currentsha
	// 		}
	// 	}
	// 	if cnt != 1 {
	// 		// fmt.Print(cnt)
	// 		fmt.Printf("%s does not match any branchname or commit hash\n", branch)
	// 	} else {
	// 		pkg.UpdateHEAD(sha)
	// 		fmt.Println("On commit ", sha, ". Entering detached HEAD state.")
	// 	}

	// }
}
