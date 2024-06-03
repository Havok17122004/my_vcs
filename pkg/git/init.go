package git

import (
	"os"
	"vcs/pkg"
)

func createBranchesDir() {
	err := os.MkdirAll("branches", 0777)
	pkg.Check(err)
}

func createHooksDir() {
	err := os.MkdirAll("hooks", 0777)
	pkg.Check(err)
}

func createInfoDir() {
	err := os.MkdirAll("info", 0777)
	pkg.Check(err)
	err = os.Chdir("info")
	pkg.Check(err)
	a := []byte("")
	err = os.WriteFile("exclude.txt", a, 0777)
	pkg.Check(err)
	err = os.Chdir("..")
	pkg.Check(err)
}

func createObjectsDir() {
	err := os.MkdirAll("objects", 0777)
	pkg.Check(err)
	err = os.Chdir("objects")
	pkg.Check(err)
	err = os.MkdirAll("info", 0777)
	pkg.Check(err)
	err = os.MkdirAll("pack", 0777)
	pkg.Check(err)
	err = os.Chdir("..")
	pkg.Check(err)
}

func createRefsDir() {
	err := os.MkdirAll("refs", 0777)
	pkg.Check(err)
	err = os.Chdir("refs")
	pkg.Check(err)
	err = os.MkdirAll("heads", 0777)
	pkg.Check(err)
	err = os.MkdirAll("tags", 0777)
	pkg.Check(err)
	err = os.Chdir("..")
	pkg.Check(err)
}

func createConfig() {
	a := []byte("")
	err := os.WriteFile("config.txt", a, 0777)
	pkg.Check(err)
}

func createDesc() {
	a := []byte("")
	err := os.WriteFile("description.txt", a, 0777)
	pkg.Check(err)
}

func createHEAD() {
	a := []byte("")
	err := os.WriteFile("HEAD.txt", a, 0777)
	pkg.Check(err)
}

func Init() {
	err := os.MkdirAll(".vcs", 0777)
	pkg.Check(err)

	err = os.Chdir(".vcs")
	pkg.Check(err)

	createConfig()
	createDesc()
	createHEAD()
	createBranchesDir()
	createHooksDir()
	createInfoDir()
	createObjectsDir()
	createRefsDir()

}
