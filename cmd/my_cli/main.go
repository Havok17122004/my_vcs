package main

import (
	"fmt"
	"os"
	"vcs/pkg/git"
)

func main() {

	args := os.Args[1:]

	switch args[0] {
	case "init":
		git.Init()
	case "config":
		os.Chdir(".vcs")
		// fmt.Println(args)
		git.Config(args)
	case "findconfig":
		os.Chdir(".vcs")
		val, _ := git.FindConfigData(args[1], args[2])
		fmt.Println(*val)
	}
}
