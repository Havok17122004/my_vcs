package pkg

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
)

// compress the file according to zlib compression and save it in .vcs/objects with the file SHA1 as the name of the file.
func CompressFileStoreInObjects(fp string, objectType string) string {
	f, err := os.Open(fp)
	Check(err)
	defer f.Close()

	stringSHA := GetSHAofFile(f)

	err = os.Chdir(VCSDirPath)
	Check(err)
	err = os.Chdir("objects")
	Check(err)

	dir_name := stringSHA[:2]
	err = os.MkdirAll(dir_name, 0777)
	Check(err)

	err = os.Chdir(dir_name)
	Check(err)

	file_name := stringSHA[2:]
	_, err = os.Stat(file_name)
	if err == nil {
		return stringSHA
	}

	outputFile, err := os.Create(file_name)
	Check(err)
	defer outputFile.Close()

	f.Seek(0, 0)
	info, _ := f.Stat()
	zlibWriter := zlib.NewWriter(outputFile)
	zlibWriter.Write([]byte(fmt.Sprintf("%s %d\x00", objectType, info.Size())))
	_, err = io.Copy(zlibWriter, f)
	Check(err)
	err = zlibWriter.Close()
	Check(err)
	fmt.Printf("Created blob %s for %s\n", stringSHA, fp)
	return stringSHA
}

// compress the string according to zlib compression and save it in .vcs/objects with the file SHA1 as the name of the file.
func CompressStringStoreInObjects(s string, objectType string) string {
	var buffer bytes.Buffer

	w := zlib.NewWriter(&buffer)
	w.Write(([]byte(fmt.Sprintf("%s %d\x00", objectType, len(s))))) //len(s) is taken as size, because each char occupies 1 byte
	_, err := w.Write([]byte(s))
	Check(err)
	err = w.Close()
	Check(err)

	err = os.Chdir(VCSDirPath)
	Check(err)
	err = os.Chdir("objects")
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

	outputFile, err := os.Create(file_name)
	Check(err)
	defer outputFile.Close()

	_, err = io.Copy(outputFile, &buffer)
	Check(err)

	return stringsha
}
