package git

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"vcs/pkg"
)

var index *pkg.Index

func Commit() {
	username, err1 := FindConfigData("user", "name")
	useremail, err2 := FindConfigData("user", "email")
	if err1 != nil || err2 != nil {
		return
	}
	fmt.Print("Enter the commit message : ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	m := strings.TrimSpace(input.Text())
	if len(m) == 0 {
		fmt.Println("commit operation aborted. valid commit message not entered")
		return
	}
	err := os.WriteFile("COMMIT_EDITMSG.txt", []byte(m), 0777)
	pkg.Check(err)
	index = pkg.ParseIndex()
	shatree, flag := makeTrees(pkg.WorkingDirPath)
	if flag {
		makeCommitObject(shatree, m, username, useremail)
	} else {
		branchname := filepath.Base(pkg.ParseHEAD())
		fmt.Printf("On branch %s\nnothing to commit, working tree clean\n", branchname)
	}

}

func GetFileMode(file os.DirEntry) string {
	info, _ := file.Info()
	if info.Mode()&0100 != 0 {
		return "100755"
	} else {
		return "100644"
	}
}

func makeTrees(path string) (string, bool) {
	dir, err := os.Open(path)
	pkg.Check(err)
	files, _ := dir.ReadDir(0)
	var s string
	for _, file := range files {
		if file.Name() == ".vcs" || file.Name() == ".git" {
			continue
		}
		if _, exist := index.Entries[filepath.Join(path, file.Name())]; !exist {
			continue
		}
		fileMode := GetFileMode(file)
		if file.IsDir() {
			newPath := filepath.Join(path, file.Name())
			shatree, flag := makeTrees(newPath)
			if flag {
				s = fmt.Sprintf("%s%s tree %s %s\n", s, fileMode, shatree, file.Name())
			}
		} else {
			newPath := filepath.Join(path, file.Name())
			f, err := os.Open(newPath)
			pkg.Check(err)
			s = fmt.Sprintf("%s%s blob %x %s\n", s, fileMode, pkg.GetSHAofFile(f), file.Name())
		}
	}
	if len(s) == 0 {
		return s, false
	}
	treesha := pkg.CompressStringStoreInObjects(s)
	fmt.Println("Created tree ", treesha, " for ", path)
	return treesha, true
}

func makeCommitObject(shatree string, message string, username *string, useremail *string) {
	shaparent, err := pkg.FetchHeadsSHA(pkg.ParseHEAD())
	s := fmt.Sprintf("tree %s\n", shatree)
	if err == nil {
		s = fmt.Sprintf("%sparent %s\n", s, shaparent)
	}
	timenow := time.Now().Unix()
	timezone := strings.Split(time.Now().String(), " ")[2]
	s = fmt.Sprintf("%sauthor %s <%s> %d %s\ncommitter %s <%s> %d %s\n\n%s\n", s, *username, *useremail, timenow, timezone, *username, *useremail, timenow, timezone, message)
	sha := pkg.CompressStringStoreInObjects(s)
	fmt.Println("Created commit object ", sha)
	pkg.UpdateHeads(sha, pkg.ParseHEAD())
	var logmessage string = message // to be updated!
	pkg.UpdateHEADlog(shaparent, sha, *username, *useremail, timenow, timezone, logmessage)

	relativebranchfilepath := pkg.ParseHEAD()

	_, err = os.Open(filepath.Join(pkg.VCSDirPath, relativebranchfilepath+".txt"))
	// pkg.Check(err)
	if err == nil {
		// fmt.Println("here")
		pkg.UpdateBranchLog(filepath.Base(relativebranchfilepath), shaparent, sha, *username, *useremail, timenow, timezone, logmessage)
	} else {

	}

}
