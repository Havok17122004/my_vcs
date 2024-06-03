package git

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"vcs/pkg"
)

func Config(s []string) {
	// place other flags here

	scopePtr := flag.String("scope", "local", "used to define the scope of config file on which operation is to be performed")
	flag.Parse()

	var err error
	var file *os.File

	switch *scopePtr {
	case "global": // how to find the path of the directory??????????????????????????????????????????????????????????????????????????????????????
		userDirPath, e := os.UserConfigDir()
		pkg.Check(e)
		file, err = os.OpenFile(filepath.Join(userDirPath, "config.txt"), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	case "local":
		file, err = os.OpenFile("config.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	default:
		err = errors.New("not a valid scope")
	}
	pkg.Check(err)
	defer file.Close()
	len := len(s)
	set(strings.Split(s[len-2], "."), s[len-1], file)
}

func set(s []string, val string, file *os.File) {

	sectionFound, _, _, lineNum, err := peekConfig(s[0], s[1], file) // sectionFound, fieldFound, value, lineNum, err
	file.Seek(0, io.SeekStart)
	if lineNum != -1 {
		fmt.Println("updating...")
		content, err := io.ReadAll(file)
		pkg.Check(err)
		lines := strings.Split(string(content), "\n")
		lines[lineNum-1] = fmt.Sprintf("\t%s = %s", s[1], val)
		updatedContent := strings.Join(lines, "\n")
		err = file.Truncate(0)
		pkg.Check(err)
		_, err = file.Seek(0, io.SeekStart)
		pkg.Check(err)
		_, err = file.WriteString(updatedContent)
		pkg.Check(err)
		return
	}

	var str string
	if !sectionFound {
		fmt.Println("adding section...")
		str = fmt.Sprintf("[%s]\n\t%s = %s\n", s[0], s[1], val)
	} else if err == io.EOF {
		fmt.Println("appending...")
		str = fmt.Sprintf("\t%s = %s\n", s[1], val)
	}
	_, err = file.WriteString(str)
	pkg.Check(err)
	// file.Close()
}

func peekConfig(section string, field string, file *os.File) (bool, bool, string, int, error) { // sectionFound, fieldFound, value, lineNum, err
	reader := bufio.NewReader(file)

	sectionFound := false
	fieldFound := false
	lineNum := -1
	var line []byte
	var lineInSection []byte
	var err error
	var value string
	j := 0
	sectionWithBrackets := fmt.Sprintf("[%s]", section)
	for {
		line, _, err = reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			pkg.Check(err)
		}
		j++
		if strings.TrimSpace(string(line)) == sectionWithBrackets {
			sectionFound = true
			for {
				lineInSection, _, err = reader.ReadLine()

				if err == io.EOF {
					break
				} else if err != nil {
					pkg.Check(err)
				}
				j++
				if strings.HasPrefix(strings.TrimSpace(string(lineInSection)), "[") {
					break
				}
				if fieldFound {
					break
				}
				if strings.HasPrefix(strings.TrimSpace(string(lineInSection)), field+" = ") {
					fieldFound = true
					valInd := 4 + len(field)
					value = string(lineInSection[valInd:])
					lineNum = j
				}
			}
			break
		}
	}
	// file.Close()
	return sectionFound, fieldFound, value, lineNum, err
}

func FindConfigData(section string, field string) (*string, error) {
	file, err := os.OpenFile("config.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	pkg.Check(err)
	var found bool
	var val string

	if _, found, val, _, _ = peekConfig(section, field, file); found {
		return &val, nil
	}

	userDirPath, e := os.UserConfigDir()
	pkg.Check(e)
	file, err = os.OpenFile(filepath.Join(userDirPath, "config.txt"), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	pkg.Check(err)
	if _, found, val, _, _ := peekConfig(section, field, file); found {
		return &val, nil
	}
	if val, found = os.LookupEnv("AUTHOR_NAME"); !found {
		fmt.Printf("%s field not updated.\n", field)
		err = errors.New("field not found")
		return &val, err
	}
	return &val, nil
}
