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

func Commit() {
	var m string
	fmt.Print("Enter the commit message : ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	m = input.Text()
	m = strings.TrimSpace(m)
	if len(m) == 0 {
		fmt.Println("commit operation aborted. valid commit message not entered")
		return
	}
	err := os.WriteFile("COMMIT_EDITMSG.txt", []byte(m), 0777)
	pkg.Check(err)

	makeBlobs(pkg.WorkingDirPath)
	shatree := makeTrees(pkg.WorkingDirPath)
	makeCommitObject(string(shatree[:]), m)
}

func makeBlobs(path string) {
	dir, err := os.Open(path)
	pkg.Check(err)
	files, _ := dir.ReadDir(0)
	for _, file := range files {
		file_path := filepath.Join(path, file.Name())
		if file.Name() == ".vcs" {
			continue
		} else if file.IsDir() {
			makeBlobs(file_path)
		} else {
			pkg.CompressFileStoreInObjects(file_path)
		}
	}
}

func makeTrees(path string) string {
	dir, err := os.Open(path)
	pkg.Check(err)
	files, _ := dir.ReadDir(0)
	var s string
	for _, file := range files {
		fileInfo, _ := os.Stat(filepath.Join(path, file.Name()))
		fileMode := fileInfo.Mode()
		if file.IsDir() {
			newPath := filepath.Join(path, file.Name())

			shatree := makeTrees(newPath)
			s = fmt.Sprintf("%s%d tree %s %s\n", s, fileMode, string(shatree[:]), file.Name())
		} else {
			newPath := filepath.Join(path, file.Name())
			f, err := os.Open(newPath)
			pkg.Check(err)
			s = fmt.Sprintf("%s%d blob %s %s\n", s, fileMode, pkg.GetSHAofFile(f), file.Name())
		}
	}
	return pkg.CompressStringStoreInObjects(s)
}

func makeCommitObject(shatree string, message string) {
	username, err := FindConfigData("user", "name")
	pkg.Check(err)
	useremail, err := FindConfigData("user", "email")
	pkg.Check(err)
	shaparent, err := pkg.FetchHEADSHA()
	s := fmt.Sprintf("tree %s\n", shatree)
	if err == nil {
		s = fmt.Sprintf("%sparent %s\n", s, shaparent)
	}
	timenow := time.Now().Unix()
	timezone := strings.Split(time.Now().String(), " ")[2]
	s = fmt.Sprintf("%sauthor %s <%s> %d %s\ncommitter %s <%s> %d %s\n\n%s", s, *username, *useremail, timenow, timezone, *username, *useremail, timenow, timezone, message)
	sha := pkg.CompressStringStoreInObjects(s)
	pkg.UpdateHEAD(string(sha[:]))
}
