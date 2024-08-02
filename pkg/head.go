package pkg

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// fetch the SHA1 of the '.git/refs/heads/branchname.txt' from the relative branch file path provided
func FetchHeadsSHAfromPath(relativebranchfilepath string) (string, error) {
	totalfilepath := relativebranchfilepath + ".txt"
	file, err := os.Open(filepath.Join(VCSDirPath, totalfilepath))
	if err != nil {
		log := ParseLog(filepath.Join(VCSDirPath, "logs/HEAD.txt"))
		for _, entry := range *log {
			if entry.Currentsha == relativebranchfilepath {
				return entry.Currentsha, nil
			}
		}
		return "0000000000000000000000000000000000000000", err
	}
	s, err := io.ReadAll(file)
	Check(err)
	defer file.Close()
	fmt.Println(string(s))
	return string(s), err
}

// update the '.git/refs/heads/branchname.txt' to the s string provided
func UpdateHeads(s string, relativebranchfilepath string) {
	os.MkdirAll(filepath.Join(VCSDirPath, "refs/heads"), 0777)
	os.MkdirAll(filepath.Join(VCSDirPath, "refs/tags"), 0777)
	fmt.Println(relativebranchfilepath)
	if strings.HasPrefix(relativebranchfilepath, "/refs/") || strings.HasPrefix(relativebranchfilepath, "refs/") {
		totalfilepath := relativebranchfilepath + ".txt"
		file, err := os.OpenFile(filepath.Join(VCSDirPath, totalfilepath), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)

		Check(err)
		_, err = file.WriteString(s)
		Check(err)
		defer file.Close()
		// fmt.Println("i am not here")
	}
	fmt.Println("sha:", s)
}

// update the HEAD.txt with the branchname
func UpdateHEAD(branchrelativepath string) {
	file, err := os.OpenFile(filepath.Join(VCSDirPath, "HEAD.txt"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777) //truncate open karne ke bad file ke contents remove kar deta hai
	Check(err)
	defer file.Close()
	fileWriter := bufio.NewWriter(file)
	s := fmt.Sprintf("%s\n", branchrelativepath) // Add a newline at the end
	_, err = fileWriter.WriteString(s)
	Check(err)
	err = fileWriter.Flush() // Flush the writer
	Check(err)
}
