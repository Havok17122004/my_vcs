package pkg

import (
	"compress/zlib"
	"io"
	"os"
)

func ReadCompressedFile(f *os.File) []byte {
	reader, err := zlib.NewReader(f)
	Check(err)
	b, err := io.ReadAll(reader)
	Check(err)
	reader.Close()
	return b
}
