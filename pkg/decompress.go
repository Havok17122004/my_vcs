package pkg

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func ReadCompressedFile(fp string) (string, string, int64) { // header of this file is `blob <size>\0` from here get the filesize

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
	// return b.String()
	data := b.String()
	fmt.Println(data)
	nullIndex := strings.IndexByte(data, '\x00')
	if nullIndex == -1 {
		panic("Invalid object: no null byte found in header")
	}
	// // Extract the header and content
	header := data[:nullIndex]
	content := data[nullIndex+1:]

	// // The header is in the format "blob <size>"
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || (headerParts[0] != "blob" && headerParts[0] != "commit" && headerParts[0] != "tree") {
		panic("Invalid blob header format")
	}

	// // Parse the size
	size, err := strconv.ParseInt(headerParts[1], 10, 64)
	if err != nil {
		panic(err)
	}

	// // Return the content and the size
	return content, headerParts[0], size
}
