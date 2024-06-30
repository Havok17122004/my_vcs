package pkg

import (
	"bytes"
	"compress/zlib"
	"io"
	"os"
)

func ReadCompressedFile(fp string) string {

	file, err := os.Open(fp)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, file)
	if err != nil {
		panic(err)
	}

	r, err := zlib.NewReader(buffer)
	if err != nil {
		panic(err)
	}
	var b bytes.Buffer
	_, err = io.Copy(&b, r)
	if err != nil {
		panic(err)
	}
	return b.String()
}
