package pkg

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

// get the SHA1 of file specified
func GetSHAofFile(f *os.File) string {
	h := sha1.New()

	_, err := io.Copy(h, f)
	Check(err)
	s := fmt.Sprintf("%x", h.Sum(nil))
	return s
}

// get the SHA1 of the string specified
func GetSHAofText(s string) [20]byte {
	return sha1.Sum([]byte(s))
}
