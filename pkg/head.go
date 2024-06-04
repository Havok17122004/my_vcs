package pkg

import (
	"io"
	"os"
)

func FetchHEADSHA() (string, error) {
	os.Chdir(VCSDirPath)
	file, err := os.Open("refs/heads/master.txt")
	if err != nil {
		return "", err
	}
	s, err := io.ReadAll(file)
	Check(err)
	return string(s), err
}

func UpdateHEAD(s string) {
	os.Chdir(VCSDirPath)
	os.MkdirAll("refs/heads", 0777)
	os.MkdirAll("refs/tags", 0777)
	os.Chdir("refs/heads")
	file, err := os.OpenFile("master.txt", os.O_CREATE|os.O_WRONLY, 0777)
	Check(err)
	file.WriteString(s)
}
