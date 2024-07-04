package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// returns hash which is referenced and the bool which is true if the text is a branchname
func FindHash(text string) (string, bool, error) {
	file, err := os.Open(filepath.Join(VCSDirPath, "/refs/heads/", text+".txt"))
	var sha string
	if err == nil {
		b := make([]byte, 40)
		file.Read(b)
		return fmt.Sprintf("%x", b), true, nil
	} else {
		logcontents := ParseLog(filepath.Join(VCSDirPath, "/logs/HEAD.txt"))
		if logcontents == nil {
			fmt.Printf("%s does not match any branchname or commit hash\n", text)
			return "", false, fmt.Errorf("%s does not match any branchname or commit hash", text) // do bar print hoga kya?
		}
		cnt := 0
		// fmt.Print(logcontents)
		for _, logline := range *logcontents {
			// fmt.Println(logline.Currentsha, branch)
			if strings.HasPrefix(logline.Currentsha, text) {
				cnt++
				sha = logline.Currentsha
			}
		}
		if cnt != 1 {
			// fmt.Print(cnt)
			fmt.Printf("%s does not match any branchname or commit hash\n", text)
			return "", false, fmt.Errorf("%s does not match any branchname or commit hash", text) // do bar print hoga kya?
		} else {
			return sha, false, nil
		}

	}
}
