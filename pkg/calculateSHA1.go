package pkg

import (
	"crypto/sha1"
	"io"
	"os"
)

func GetSHA(f *os.File) []byte {
	h := sha1.New()

	_, err := io.Copy(h, f)
	Check(err)

	return h.Sum(nil)
}
