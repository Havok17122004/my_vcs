package pkg

import (
	"compress/zlib"
	"io"
	"os"
)

func CompressStoreInObjects(f *os.File) {
	s := GetSHA(f)

	err := os.Chdir("objects")
	Check(err)

	dir_name := string(s[:2])
	err = os.MkdirAll(dir_name, 0777)
	Check(err)

	err = os.Chdir(dir_name)
	Check(err)

	outputFile, err := os.Create(string(s[2:]))
	Check(err)

	zlibWriter := zlib.NewWriter(outputFile)
	_, err = io.Copy(zlibWriter, f)
	Check(err)
}
