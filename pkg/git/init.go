package git

import (
	"os"
	"path/filepath"
	"vcs/pkg"
)

func createInfoDir() {
	os.MkdirAll(filepath.Join(pkg.VCSDirPath, "info"), 0777)
	a := []byte("")
	err := os.WriteFile(filepath.Join(pkg.VCSDirPath, "info/exclude.txt"), a, 0777)
	pkg.Check(err)
}

func createConfig() {
	a := []byte("")
	err := os.WriteFile(filepath.Join(pkg.VCSDirPath, "config.txt"), a, 0777)
	pkg.Check(err)
}

func createDesc() {
	a := []byte("")
	err := os.WriteFile(filepath.Join(pkg.VCSDirPath, "description.txt"), a, 0777)
	pkg.Check(err)
}

func createHEAD() {
	a := []byte("")
	err := os.WriteFile(filepath.Join(pkg.VCSDirPath, "HEAD.txt"), a, 0777)
	pkg.Check(err)
}

func Init(args []string) {
	var wd string
	wd, _ = os.Getwd()
	if len(args) != 1 {
		wd = filepath.Join(wd, args[1])
		os.MkdirAll(wd, 0777)
	}

	pkg.SetWorkingDirPath(wd)
	os.Chdir(wd)
	os.MkdirAll(filepath.Join(pkg.WorkingDirPath, ".vcs"), 0777)

	pkg.SetVCSDirPath(filepath.Join(wd, ".vcs"))
	createConfig()
	createDesc()
	createHEAD()
	os.MkdirAll(filepath.Join(pkg.VCSDirPath, "branches"), 0777)
	os.MkdirAll(filepath.Join(pkg.VCSDirPath, "hooks"), 0777)
	createInfoDir()
	os.MkdirAll(filepath.Join(pkg.VCSDirPath, "objects"), 0777)
	os.MkdirAll(filepath.Join(pkg.VCSDirPath, "objects/info"), 0777)
	os.MkdirAll(filepath.Join(pkg.VCSDirPath, "objects/pack"), 0777)
	os.MkdirAll(filepath.Join(pkg.VCSDirPath, "refs"), 0777)
	os.MkdirAll(filepath.Join(pkg.VCSDirPath, "refs/heads"), 0777)
	os.MkdirAll(filepath.Join(pkg.VCSDirPath, "refs/tags"), 0777)
}
