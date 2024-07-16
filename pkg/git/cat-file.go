package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"vcs/pkg"
)

func Catfile(flag string, incompletehash string) {
	if len(incompletehash) < 3 {
		fmt.Println("please enter a valid hash")
		return
	}
	entries, err := os.ReadDir(filepath.Join(pkg.VCSDirPath, "objects", incompletehash[:2]))
	if err != nil {
		fmt.Println("please enter a valid hash")
		return
	}
	count := 0
	var hash string
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), incompletehash[2:]) {
			count++
			hash = incompletehash[:2] + entry.Name()
		}
	}
	if count != 1 {
		fmt.Println("please enter a valid hash")
		return
	}

	str, size, header := pkg.ReadCompressedFile(filepath.Join(pkg.VCSDirPath, "objects", hash[:2], hash[2:]))
	if flag == "-p" {
		fmt.Println(str)
	} else if flag == "-t" {
		fmt.Println(header)
	} else if flag == "-s" {
		fmt.Println(size)
	}
}
