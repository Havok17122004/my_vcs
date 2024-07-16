package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"vcs/pkg"
)

func CreateBranch(branchname string) {
	file, err := os.OpenFile(filepath.Join(pkg.VCSDirPath, "/refs/heads", branchname+".txt"), os.O_WRONLY|os.O_CREATE, 0777)
	pkg.Check(err)
	s, err := pkg.FetchHeadsSHAfromPath(pkg.ParseHEADgivePath())
	pkg.Check(err)
	file.WriteString(s)
	username, err1 := ParseConfigData("user", "name")
	useremail, err2 := ParseConfigData("user", "email")
	if err1 != nil || err2 != nil {
		return
	}
	fmt.Println("Created branch ", branchname)
	logmessage := fmt.Sprintf("Created from %s\n", filepath.Base(pkg.ParseHEADgivePath()))
	pkg.UpdateBranchLog(branchname, s, s, username, useremail, time.Now().Unix(), strings.Split(time.Now().String(), " ")[2], logmessage, "branch")
	defer file.Close()
}

func ListBranches() {
	list, err := os.ReadDir(filepath.Join(pkg.VCSDirPath, "refs/heads"))
	if err != nil {
		fmt.Println("Corrupted internal .vcs structure. refs/heads can't be reached")
		return
	}
	headBranch := filepath.Base(pkg.ParseHEADgivePath())
	for _, entry := range list {
		if entry.Name() == headBranch {
			fmt.Printf("\x1b[42m*%s\x1b[49m\n", headBranch)
		} else {
			fmt.Println(entry.Name())
		}
	}
}
