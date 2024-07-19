package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// returns hash which is referenced and the bool which is true if the text is a branchname
func FindHashofCommit(text string) (string, bool, LogContents, error) {
	if text == "HEAD" {
		sha, _ := FetchHeadsSHAfromPath(ParseHEADgivePath())
		return sha, true, LogContents{}, nil
	}
	re := regexp.MustCompile(`^(.+)~(\d+)$`)
	matches := re.FindStringSubmatch(text)
	fmt.Println(matches)
	if len(matches) == 3 {
		branch := matches[1]
		n, err := strconv.Atoi(matches[2])
		Check(err)
		var logs *[]LogContents
		if branch == "HEAD" {
			logs = ParseLog(filepath.Join(VCSDirPath, "logs/HEAD.txt"))
		} else {
			_, err := os.Open(filepath.Join(VCSDirPath, "/refs/heads/", text+".txt"))
			if err != nil {
				fmt.Println("does not match any branch or commit hash")
				return "", false, LogContents{}, err
			}
			logs = ParseLog(filepath.Join(VCSDirPath, "logs/refs/heads", branch))
		}
		if n < 0 || len(*logs)-1-n < 0 {
			return "", false, LogContents{}, fmt.Errorf("n must not be negative")
		}
		logLine := (*logs)[len(*logs)-1-n]
		return logLine.Currentsha, false, logLine, nil
	}

	file, err := os.Open(filepath.Join(VCSDirPath, "/refs/heads/", text+".txt"))
	if err == nil {
		b := make([]byte, 40)
		file.Read(b)
		return string(b), true, LogContents{}, nil
	} else {
		logcontents := ParseLog(filepath.Join(VCSDirPath, "/logs/HEAD.txt"))
		if logcontents == nil {
			fmt.Printf("%s does not match any branchname or commit hash\n", text)
			return "", false, LogContents{}, fmt.Errorf("%s does not match any branchname or commit hash", text) // do bar print hoga kya?
		}
		cnt := 0
		// fmt.Print(logcontents)
		var line LogContents
		for _, logline := range *logcontents {
			// fmt.Println(logline.Currentsha, branch)
			if strings.HasPrefix(logline.Currentsha, text) {
				cnt++
				line = logline
			}
		}
		if cnt != 1 {
			// fmt.Print(cnt)
			fmt.Printf("%s does not match any branchname or commit hash\n", text)
			return "", false, LogContents{}, fmt.Errorf("%s does not match any branchname or commit hash", text) // do bar print hoga kya?
		} else {
			return line.Currentsha, false, line, nil
		}

	}
}
