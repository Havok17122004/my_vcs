package main

import (
	"fmt"
	"os"
	"vcs/pkg"
	"vcs/pkg/git"
)

func main() {

	args := os.Args[1:]

	if args[0] != "init" {
		s, _ := os.Getwd()
		pkg.SetWorkingDirPath(s)
		os.Chdir(".vcs")
		s, _ = os.Getwd()
		pkg.SetVCSDirPath(s)
	}

	switch args[0] {
	case "init":
		git.Init()
	case "config":
		git.Config(args)
	case "findconfig":
		val, _ := git.FindConfigData(args[1], args[2])
		fmt.Println(*val)
	case "commit":
		git.Commit()
	}
}
