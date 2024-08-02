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

// TODO: Refactor code by adding structs that can efficiently manage the peekConfig and other functions.
// read/edit the config file
func Config(s []string) {
	// place other flags here
	scopePtr := flag.String("scope", "local", "used to define the scope of config file on which operation is to be performed")
	flag.Parse()

	var err error
	var file *os.File

	switch *scopePtr {
	// case "global": // how to find the path of the directory??????????????????????????????????????????????????????????????????????????????????????
	// userDirPath, e := os.UserConfigDir()
	// pkg.Check(e)
	// file, err = os.OpenFile(filepath.Join(userDirPath, "config.txt"), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	case "local":
		file, err = os.OpenFile(filepath.Join(pkg.VCSDirPath, "config.txt"), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	default:
		fmt.Println("not a valid scope")
		return
	}
	pkg.Check(err)
	defer file.Close()
	len := len(s)
	setConfigValues(strings.Split(s[len-2], "."), s[len-1], file) // TODO: Improve the flag parsing code.
}

// set the values of the config of the field defined in s as val.
func setConfigValues(s []string, val string, file *os.File) {
	sectionFound, _, _, lineNum, err := peekConfig(s[0], s[1], file) // sectionFound, fieldFound, value, lineNum, err
	file.Seek(0, io.SeekStart)
	if lineNum != -1 {
		fmt.Printf("updating %s field of %s section\n", s[1], s[0])
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
		fmt.Printf("adding section %s\n", s[0])
		fmt.Printf("appending field %s\n", s[1])
		str = fmt.Sprintf("[%s]\n\t%s = %s\n", s[0], s[1], val)
	} else if err == io.EOF {
		fmt.Printf("appending field %s\n", s[1])
		str = fmt.Sprintf("\t%s = %s\n", s[1], val)
	}
	_, err = file.WriteString(str)
	pkg.Check(err)
	// file.Close()
}

// read the config values as already present in the .config file
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
		// fmt.Println(string(line))
		if err == io.EOF {
			break
		} else if err != nil {
			pkg.Check(err)
		}
		j++
		// fmt.Println(strings.TrimSpace(string(line)), sectionWithBrackets)
		if strings.TrimSpace(string(line)) == sectionWithBrackets {
			// fmt.Print("dvsd")
			sectionFound = true
			for {
				lineInSection, _, err = reader.ReadLine()
				// fmt.Println(strings.TrimSpace(string(lineInSection)), field+" = ")
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
					// fmt.Println("eg")
				}
			}
			break
		}
	}
	// file.Close()
	return sectionFound, fieldFound, value, lineNum, err
}

// find the data related to the section and field in the config file.
func ParseConfigData(section string, field string) (string, error) {
	os.Chdir(pkg.VCSDirPath)

	file, err := os.OpenFile("config.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	pkg.Check(err)
	var found bool
	var val string

	if _, found, val, _, _ = peekConfig(section, field, file); found {
		return val, nil
	}

	userDirPath, e := os.UserConfigDir()
	pkg.Check(e)
	file, err = os.OpenFile(filepath.Join(userDirPath, "config.txt"), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	pkg.Check(err)
	if _, found, val, _, _ := peekConfig(section, field, file); found {
		return val, nil
	}
	if val, found = os.LookupEnv("AUTHOR_NAME"); !found {
		fmt.Printf("%s field not updated.\n", field)
		err = errors.New("field not found")
		return val, err
	}
	return val, nil
}
