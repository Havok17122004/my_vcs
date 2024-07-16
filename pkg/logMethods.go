package pkg

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func UpdateHEADlog(parentsha string, currentsha string, authorname string, authoremail string, timestamp int64, gmt string, message string, operation string) {
	str := fmt.Sprintf("%s %s %s <%s> %d %s %s: %s\n", parentsha, currentsha, authorname, authoremail, timestamp, gmt, operation, message)
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

func UpdateBranchLog(branch string, parentsha string, currentsha string, authorname string, authoremail string, timestamp int64, gmt string, message string, operation string) {
	str := fmt.Sprintf("%s %s %s <%s> %d %s %s: %s\n", parentsha, currentsha, authorname, authoremail, timestamp, gmt, operation, message)
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
