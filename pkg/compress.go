package pkg

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
)

func CompressFileStoreInObjects(fp string) string {
	f, err := os.Open(fp)
	Check(err)
	defer f.Close()

	s := GetSHAofFile(f)
	stringSHA := fmt.Sprintf("%x", s)

	os.Chdir(VCSDirPath)
	err = os.Chdir("objects")
	Check(err)

	dir_name := stringSHA[:2]
	// fmt.Println(stringSHA)
	err = os.MkdirAll(dir_name, 0777)
	Check(err)

	err = os.Chdir(dir_name)
	Check(err)

	file_name := stringSHA[2:]
	_, err = os.Stat(file_name)
	if err == nil {
		return stringSHA
	}

	outputFile, err := os.OpenFile(file_name, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
	Check(err)
	defer outputFile.Close()

	f.Seek(0, 0)

	zlibWriter := zlib.NewWriter(outputFile)
	_, err = io.Copy(zlibWriter, f)
	Check(err)
	// zlibWriter.Close()
	return stringSHA
}

func CompressStringStoreInObjects(s string) string {
	var buffer bytes.Buffer

	w := zlib.NewWriter(&buffer)

	w.Write([]byte(s))

	defer w.Close()

	os.Chdir(VCSDirPath)
	err := os.Chdir("objects")
	Check(err)

	sha := GetSHAofText(s)
	stringsha := fmt.Sprintf("%x", sha)
	dir_name := stringsha[:2]
	err = os.MkdirAll(dir_name, 0777)
	Check(err)

	err = os.Chdir(dir_name)
	Check(err)

	file_name := stringsha[2:]
	_, err = os.Stat(file_name)
	if err == nil {
		return stringsha
	}

	outputFile, err := os.Create(stringsha[2:])
	Check(err)

	_, err = io.Copy(outputFile, &buffer)
	Check(err)

	return stringsha
}
