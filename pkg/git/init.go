package git

import (
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func createBranchesDir() {
	err := os.Mkdir("branches", 0777)
	check(err)
}

func createHooksDir() {
	err := os.Mkdir("hooks", 0777)
	check(err)
}

func createInfoDir() {
	err := os.Mkdir("info", 0777)
	check(err)
	err = os.Chdir("info")
	check(err)
	a := []byte("")
	err = os.WriteFile("exclude.txt", a, 0777)
	check(err)
	err = os.Chdir("..")
	check(err)
}

func createObjectsDir() {
	err := os.Mkdir("objects", 0777)
	check(err)
	err = os.Chdir("objects")
	check(err)
	err = os.Mkdir("info", 0777)
	check(err)
	err = os.Mkdir("pack", 0777)
	check(err)
	err = os.Chdir("..")
	check(err)
}

func createRefsDir() {
	err := os.Mkdir("refs", 0777)
	check(err)
	err = os.Chdir("refs")
	check(err)
	err = os.Mkdir("heads", 0777)
	check(err)
	err = os.Mkdir("tags", 0777)
	check(err)
	err = os.Chdir("..")
	check(err)
}

func createConfig() {
	a := []byte("")
	err := os.WriteFile("config.txt", a, 0777)
	check(err)
}

func createDesc() {
	a := []byte("")
	err := os.WriteFile("description.txt", a, 0777)
	check(err)
}

func createHEAD() {
	a := []byte("")
	err := os.WriteFile("HEAD.txt", a, 0777)
	check(err)
}

func Init() {
	if _, err := os.Stat(".vcs"); !os.IsNotExist(err) {
		os.RemoveAll(".vcs")
	}

	err := os.MkdirAll(".vcs", 0777)
	check(err)

	err = os.Chdir(".vcs")
	check(err)

	createConfig()
	createDesc()
	createHEAD()
	createBranchesDir()
	createHooksDir()
	createInfoDir()
	createObjectsDir()
	createRefsDir()

}
