package pkg

import (
	"crypto/sha1"
	"io"
	"os"
)

func GetSHAofFile(f *os.File) []byte {
	h := sha1.New()

	_, err := io.Copy(h, f)
	Check(err)

	return h.Sum(nil)
}

func GetSHAofText(s string) [20]byte {
	return sha1.Sum([]byte(s))
}
