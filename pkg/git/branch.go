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
	s, err := pkg.FetchHeadsSHA(pkg.ParseHEAD())
	pkg.Check(err)
	file.WriteString(s)
	username, err1 := FindConfigData("user", "name")
	useremail, err2 := FindConfigData("user", "email")
	if err1 != nil || err2 != nil {
		return
	}
	fmt.Println("Created branch ", branchname)
	var logmessage string
	pkg.UpdateBranchLog(branchname, s, s, username, useremail, time.Now().Unix(), strings.Split(time.Now().String(), " ")[2], logmessage)
	defer file.Close()
}
