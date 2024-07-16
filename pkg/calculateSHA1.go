package pkg

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

func GetSHAofFile(f *os.File) string {
	h := sha1.New()

	_, err := io.Copy(h, f)
	Check(err)
	s := fmt.Sprintf("%x", h.Sum(nil))
	return s
}

func GetSHAofText(s string) [20]byte {
	return sha1.Sum([]byte(s))
}
