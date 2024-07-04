package git

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"vcs/pkg"
)

type CommitObject struct {
	Treesha        string
	Parentsha      string
	Authorname     string
	Authoremail    string
	CommitterName  string
	CommitterEmail string
	Time           int64
	Timezone       string
	Message        string
}

type TreeEntry struct {
	FileMode string
	FileType string
	Sha      string
	Name     string
}

var index *pkg.Index

func Commit() {
	var c CommitObject
	var err1, err2 error
	c.Authorname, err1 = FindConfigData("user", "name")
	c.Authoremail, err2 = FindConfigData("user", "email")
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
	// fmt.Println(files)

	var s string
	for _, file := range files {
		if file.Name() == ".vcs" || file.Name() == ".git" {
			continue
		}
		// fmt.Println(file.Name(), file.IsDir())
		// if _, exist := index.Entries[filepath.Join(path, file.Name())]; !exist {
		// 	continue
		// }
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
				s = fmt.Sprintf("%s%s tree %s %s\n", s, fileMode, shatree, file.Name())
			}
		} else {
			entry := index.Entries[filepath.Join(path, file.Name())]
			s = fmt.Sprintf("%s%s blob %s %s\n", s, fileMode, entry.Id, file.Name())
		}
	}
	if len(s) == 0 {
		return s, false
	}
	// treesha := pkg.GetSHAofText(s)
	// stringsha := fmt.Sprintf("%x", treesha)
	// commitsha, err := pkg.FetchHeadsSHA(pkg.ParseHEAD())
	// // pkg.Check(err)
	// var prev TreeEntry
	// if err == nil {
	// 	prev := GetCommittedEntry(path, commitsha)
	// 	fmt.Println(prev, path)
	// 	if prev.Sha == stringsha {
	// 		return stringsha, false
	// 	}
	// }
	// fmt.Println(stringsha, prev.Sha)
	// fmt.Println("Created tree ", stringsha, " for ", path)
	stringsha := pkg.CompressStringStoreInObjects(s)
	return stringsha, true
}

func makeCommitObject(c *CommitObject) {
	var err error
	c.Parentsha, err = pkg.FetchHeadsSHA(pkg.ParseHEAD())
	s := fmt.Sprintf("tree %s\n", c.Treesha)
	if err == nil {
		s = fmt.Sprintf("%sparent %s\n", s, c.Parentsha)
	}
	prevCommit := ParseCommit(c.Parentsha)
	if prevCommit.Treesha == c.Treesha {
		branchname := filepath.Base(pkg.ParseHEAD())
		fmt.Printf("On branch %s\nnothing to commit, working tree clean\n", branchname)
		return
	}
	c.Time = time.Now().Unix()
	c.Timezone = strings.Split(time.Now().String(), " ")[2]
	s = fmt.Sprintf("%sauthor %s <%s> %d %s\ncommitter %s <%s> %d %s\n\n%s\n", s, c.Authorname, c.Authoremail, c.Time, c.Timezone, c.CommitterName, c.CommitterEmail, c.Time, c.Timezone, c.Message)
	sha := pkg.CompressStringStoreInObjects(s)
	fmt.Println("Created commit object ", sha)
	pkg.UpdateHeads(sha, pkg.ParseHEAD())
	var logmessage string = c.Message // to be updated!
	pkg.UpdateHEADlog(c.Parentsha, sha, c.Authorname, c.Authoremail, c.Time, c.Timezone, logmessage)

	relativebranchfilepath := pkg.ParseHEAD()

	_, err = os.Open(filepath.Join(pkg.VCSDirPath, relativebranchfilepath+".txt"))
	// pkg.Check(err)
	if err == nil {
		// fmt.Println("here")
		pkg.UpdateBranchLog(filepath.Base(relativebranchfilepath), c.Parentsha, sha, c.Authorname, c.Authoremail, c.Time, c.Timezone, logmessage)
	} else {

	}

}

func ParseCommit(hash string) *CommitObject {
	// file, err := os.Open(filepath.Join(pkg.VCSDirPath, "objects", hash[:2], hash[2:]))
	contents := pkg.ReadCompressedFile(filepath.Join(pkg.VCSDirPath, "objects", hash[:2], hash[2:]))
	var c CommitObject

	lines := strings.Split(contents, "\n")
	if strings.HasPrefix(lines[0], "tree") {
		c.Treesha = lines[0][5:]
	} else {
		pkg.Check(fmt.Errorf("Commit file does not contain reference to tree"))
	}
	i := 0
	if strings.HasPrefix(lines[1], "parent") {
		c.Parentsha = lines[1][7:]
	} else {
		c.Parentsha = "0000000000000000000000000000000000000000"
		i = 1
	}
	// after "author"

	l := strings.SplitN(lines[2-i], " ", 5)
	c.Authorname = l[1]
	c.Authoremail, _ = strings.CutSuffix(l[2], ">")
	c.Authoremail, _ = strings.CutPrefix(c.Authoremail, "<")
	c.Time, _ = strconv.ParseInt(l[3], 10, 64)
	c.Timezone = l[4]

	l = strings.SplitN(lines[3-i], " ", 5)
	// fmt.Println(lines)
	c.CommitterName = l[1]
	c.CommitterEmail, _ = strings.CutSuffix(l[2], ">")
	c.CommitterEmail, _ = strings.CutPrefix(c.Authoremail, "<")

	c.Message = lines[5-i]
	// fmt.Println(lines[5])
	// fmt.Println(lines[6])
	// fmt.Println(lines[7])
	return &c
}

func ParseTree(hash string) *[]TreeEntry {
	// fmt.Println(hash)
	if hash == "" {
		return &[]TreeEntry{{"", "", "", ""}}
	}
	contents := pkg.ReadCompressedFile(filepath.Join(pkg.VCSDirPath, "objects", hash[:2], hash[2:]))
	lines := strings.Split(contents, "\n")
	var treeObjects []TreeEntry
	// fmt.Println(lines)
	for _, line := range lines {
		// fmt.Println(line)
		var obj TreeEntry
		l := strings.SplitN(line, " ", 4)
		if len(l) == 1 {
			continue
		}
		obj.FileMode = l[0]
		obj.FileType = l[1]
		obj.Sha = l[2]
		obj.Name = l[3]
		treeObjects = append(treeObjects, obj)
	}
	return &treeObjects
}
