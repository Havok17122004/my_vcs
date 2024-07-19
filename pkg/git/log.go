package git

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"vcs/pkg"
)

func Log(args []string) {
	toBeDisplayed := make(map[int64]pkg.LogContents)
	exclusions := make(map[int64]pkg.LogContents)
	if len(args) == 0 {
		contents := pkg.ParseLog(filepath.Join(pkg.VCSDirPath, "logs/HEAD.txt"))
		for _, content := range *contents {
			toBeDisplayed[content.Timestamp] = content
		}
	} else {
		for _, val := range args {
			excluded := false
			if strings.HasPrefix(val, "^") {
				excluded = true
				val = strings.TrimLeft(val, "^")
			}
			_, isBranch, logLine, err := pkg.FindHashofCommit(val)
			pkg.Check(err)
			if isBranch {
				var contents *[]pkg.LogContents
				if val == "HEAD" {
					contents = pkg.ParseLog(filepath.Join(pkg.VCSDirPath, "logs/HEAD.txt"))
				} else {
					contents = pkg.ParseLog(filepath.Join(pkg.VCSDirPath, "logs/refs/heads", val+".txt"))
				}
				if excluded {
					for _, content := range *contents {
						exclusions[content.Timestamp] = content
					}
				} else {
					for _, content := range *contents {
						toBeDisplayed[content.Timestamp] = content
					}
				}

			} else {
				if excluded {
					exclusions[logLine.Timestamp] = logLine
				} else {
					toBeDisplayed[logLine.Timestamp] = logLine
				}
			}
		}
	}
	fmt.Println(toBeDisplayed)
	fmt.Println(exclusions)
	for key, value := range toBeDisplayed {
		val, found := exclusions[key]
		// fmt.Println(val, value)
		if found && val == value { // can i compare like this?
			delete(toBeDisplayed, key)
		}
	}
	keys := make([]int64, 0, len(toBeDisplayed))
	for k := range toBeDisplayed {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})
	for i := 0; i < len(keys); i++ {
		oneObj := toBeDisplayed[keys[i]]
		fmt.Println(oneObj)
		fmt.Printf("%scommit %s%s\n", Yellow, oneObj.Currentsha, Default) // why was there (branchname) written on some lines? add that too!
		fmt.Printf("Author:\t%s <%s>\n", oneObj.Authorname, oneObj.Authoremail)
		t := time.Unix(oneObj.Timestamp, 0)
		formattedTime := t.Format("Mon Jan 2 15:04:05 2006 -0700")
		fmt.Printf("Date:\t%s\n\n", formattedTime)
		fmt.Printf("\t%s\n\n", oneObj.Message)
	}
}
