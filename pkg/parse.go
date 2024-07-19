package pkg

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ParseCommit(hash string) *CommitObject {
	var c CommitObject
	if hash == "0000000000000000000000000000000000000000" {
		return &c
	}
	// file, err := os.Open(filepath.Join(VCSDirPath, "objects", hash[:2], hash[2:]))
	contents, filetype, _ := ReadCompressedFile(filepath.Join(VCSDirPath, "objects", hash[:2], hash[2:]))
	if filetype != "commit" {
		Check(fmt.Errorf("invalid commit object: header does not match commit header"))
	}

	// nullIndex := strings.IndexByte(contents, '\x00')
	// if nullIndex == -1 {
	// 	Check(fmt.Errorf("invalid commit object: no null byte found in header"))
	// }
	// // commitSize := contents[7:nullIndex] // not needed yet.
	// entries := contents[nullIndex+1:]
	lines := strings.Split(contents, "\n")

	if strings.HasPrefix(lines[0], "tree") {
		c.Treesha = lines[0][5:]
	} else {
		Check(fmt.Errorf("commit file does not contain reference to tree"))
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
	// fmt.Println(c)
	return &c
}

func ParseTree(hash string) *[]TreeEntry {
	// fmt.Println(hash)
	if hash == "" {
		return &[]TreeEntry{{"", "", "", ""}}
	}
	entries, filetype, _ := ReadCompressedFile(filepath.Join(VCSDirPath, "objects", hash[:2], hash[2:]))
	if filetype != "tree" {
		Check(fmt.Errorf("invalid tree object: header does not match tree header"))
	}

	var treeObjects []TreeEntry
	lines := strings.Split(entries, "\n")
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
		obj.Path = l[3]
		treeObjects = append(treeObjects, obj)
	}
	return &treeObjects
}

func ParseLog(logpath string) *[]LogContents {
	dir := filepath.Dir(logpath)
	err := os.MkdirAll(dir, 0777)
	Check(err)

	file, err := os.OpenFile(logpath, os.O_RDONLY|os.O_CREATE, 0777)
	Check(err)
	defer file.Close()

	var parsedContentsLog []LogContents
	fileScanner := bufio.NewScanner(file)
	// fmt.Print("baka")
	for fileScanner.Scan() {
		line := fileScanner.Text()
		// fmt.Println(line, "ppp")
		contents := strings.SplitN(line, " ", 8)
		if len(contents) < 8 {
			continue // Skip lines that do not have enough fields
		}
		var log LogContents
		log.Parentsha = contents[0]
		log.Currentsha = contents[1]
		log.Authorname = contents[2]
		log.Authoremail = strings.TrimRight(strings.TrimLeft(contents[3], "<"), ">")
		log.Timestamp, err = strconv.ParseInt(contents[4], 10, 64)
		Check(err)
		log.Gmt = contents[5]
		log.Operation = contents[6][:len(contents[6])-1]
		log.Message = contents[7]
		// log.branchname = filepath.Base(logpath) // Uncomment this if branch name is required

		parsedContentsLog = append(parsedContentsLog, log)
	}

	if err := fileScanner.Err(); err != nil {
		Check(err)
	}

	return &parsedContentsLog
}

func ParseHEADgivePath() string {
	file, err := os.OpenFile(filepath.Join(VCSDirPath, "HEAD.txt"), os.O_CREATE|os.O_RDONLY, 0777)
	Check(err)
	defer file.Close()
	fileScanner := bufio.NewScanner(file)
	fileScanner.Scan() // Scan the file
	line := fileScanner.Text()
	if len(line) == 0 {
		UpdateHEAD("refs: refs/heads/master")
		return "refs/heads/master"
	}
	if strings.HasPrefix(line, "refs: ") {
		return line[6:]
	} else {
		return line
	}
}
