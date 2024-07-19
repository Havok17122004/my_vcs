package git //change there

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
	var c pkg.CommitObject
	var err1, err2 error
	c.Authorname, err1 = ParseConfigData("user", "name")
	c.Authoremail, err2 = ParseConfigData("user", "email")
	c.CommitterName = c.Authorname
	c.CommitterEmail = c.Authoremail

	if err1 != nil || err2 != nil {
		return
	}
	fmt.Print("Enter the commit message : ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	c.Message = strings.TrimSpace(input.Text())
	if len(c.Message) == 0 {
		fmt.Println("commit operation aborted. valid commit message not entered")
		return
	}
	err := os.WriteFile("COMMIT_EDITMSG.txt", []byte(c.Message), 0777)
	pkg.Check(err)
	index = pkg.ParseIndex()
	var flag bool
	c.Treesha, flag = makeTrees(pkg.WorkingDirPath)
	// prevCommit := ParseCommit()
	if flag {
		makeCommitObject(&c)
	} else {
		branchname := filepath.Base(pkg.ParseHEADgivePath())
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
	fmt.Println(index)
	for _, file := range files {
		if file.Name() == ".vcs" || file.Name() == ".git" {
			continue
		}
		exist := false
		for key := range index.Entries {
			if strings.HasPrefix(key, filepath.Join(path, file.Name())) {
				exist = true
				break
			}
		}
		if !exist {
			continue
		}
		fileMode := GetFileMode(file)

		if file.IsDir() {
			newPath := filepath.Join(path, file.Name())
			fmt.Println("recursive call for ", newPath)
			shatree, flag := makeTrees(newPath)
			if flag {
				s = fmt.Sprintf("%s%s tree %s %s\n", s, fileMode, shatree, newPath)
			}
		} else {
			entry := index.Entries[filepath.Join(path, file.Name())]
			s = fmt.Sprintf("%s%s blob %s %s\n", s, fileMode, entry.Id, filepath.Join(path, file.Name()))
		}
	}
	if len(s) == 0 {
		return s, false
	}
	stringsha := pkg.CompressStringStoreInObjects(s, "tree")
	return stringsha, true
}

func makeCommitObject(c *pkg.CommitObject) {
	var err error
	c.Parentsha, err = pkg.FetchHeadsSHAfromPath(pkg.ParseHEADgivePath())
	s := fmt.Sprintf("tree %s\n", c.Treesha)
	var operation string
	if err == nil {
		s = fmt.Sprintf("%sparent %s\n", s, c.Parentsha)
		if c.Parentsha != "0000000000000000000000000000000000000000" {
			prevCommit := pkg.ParseCommit(c.Parentsha)
			if prevCommit.Treesha == c.Treesha {
				branchname := filepath.Base(pkg.ParseHEADgivePath())
				fmt.Printf("On branch %s\nnothing to commit, working tree clean\n", branchname)
				return
			}
			operation = "commit"
		} else {
			operation = "commit(initial)"
		}
	}
	c.Time = time.Now().Unix()
	c.Timezone = strings.Split(time.Now().String(), " ")[2]
	s = fmt.Sprintf("%sauthor %s <%s> %d %s\ncommitter %s <%s> %d %s\n\n%s\n", s, c.Authorname, c.Authoremail, c.Time, c.Timezone, c.CommitterName, c.CommitterEmail, c.Time, c.Timezone, c.Message)
	sha := pkg.CompressStringStoreInObjects(s, "commit")
	fmt.Println("Created commit object ", sha)
	pkg.UpdateHeads(sha, pkg.ParseHEADgivePath())
	var logmessage string = c.Message // to be updated!
	pkg.UpdateHEADlog(c.Parentsha, sha, c.Authorname, c.Authoremail, c.Time, c.Timezone, logmessage, operation)

	relativebranchfilepath := pkg.ParseHEADgivePath()

	_, err = os.Open(filepath.Join(pkg.VCSDirPath, relativebranchfilepath+".txt"))
	if err == nil {
		pkg.UpdateBranchLog(filepath.Base(relativebranchfilepath), c.Parentsha, sha, c.Authorname, c.Authoremail, c.Time, c.Timezone, logmessage, operation)
	} else {

	}

}
