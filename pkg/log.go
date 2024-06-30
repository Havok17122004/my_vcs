package pkg

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type logContents struct {
	// branchname  string
	Parentsha   string
	Currentsha  string
	Authorname  string
	Authoremail string
	Timestamp   string
	Gmt         string
	Message     string
}

func UpdateHEADlog(parentsha string, currentsha string, authorname string, authoremail string, timestamp int64, gmt string, message string) {
	str := fmt.Sprintf("%s %s %s <%s> %d %s %s\n", parentsha, currentsha, authorname, authoremail, timestamp, gmt, message)
	err := os.MkdirAll(filepath.Join(VCSDirPath, "logs/"), 0777)
	Check(err)
	file, err := os.OpenFile(filepath.Join(VCSDirPath, "logs/HEAD.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	Check(err)
	defer file.Close()
	fileWriter := bufio.NewWriter(file)
	_, err = fileWriter.WriteString(str)
	Check(err)
	err = fileWriter.Flush()
	Check(err)
}

func UpdateBranchLog(branch string, parentsha string, currentsha string, authorname string, authoremail string, timestamp int64, gmt string, message string) {
	str := fmt.Sprintf("%s %s %s <%s> %d %s %s\n", parentsha, currentsha, authorname, authoremail, timestamp, gmt, message)
	branch = branch + ".txt"
	err := os.MkdirAll(filepath.Join(VCSDirPath, "logs/refs/heads"), 0777)
	Check(err)
	file, err := os.OpenFile(filepath.Join(VCSDirPath, "logs/refs/heads", branch), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	Check(err)
	defer file.Close()
	fileWriter := bufio.NewWriter(file)
	_, err = fileWriter.WriteString(str)
	Check(err)
	err = fileWriter.Flush()
	Check(err)
}

func ParseLog(logpath string) *[]logContents {
	dir := filepath.Dir(logpath)
	err := os.MkdirAll(dir, 0777)
	Check(err)

	file, err := os.OpenFile(logpath, os.O_RDONLY|os.O_CREATE, 0777)
	Check(err)
	defer file.Close()

	var parsedContentsLog []logContents
	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		contents := strings.Split(line, " ")
		// if len(contents) < 7 {
		// 	continue // Skip lines that do not have enough fields
		// }
		var log logContents
		log.Parentsha = contents[0]
		log.Currentsha = contents[1]
		log.Authorname = contents[2]
		log.Authoremail = contents[3]
		log.Timestamp = contents[4]
		log.Gmt = contents[5]
		log.Message = contents[6]
		// log.branchname = filepath.Base(logpath) // Uncomment this if branch name is required

		parsedContentsLog = append(parsedContentsLog, log)
	}

	if err := fileScanner.Err(); err != nil {
		Check(err)
	}

	return &parsedContentsLog
}
