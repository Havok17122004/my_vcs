package main

import (
	"fmt"
	"os"
	"strings"
	"vcs/cmd"
	"vcs/pkg"
	"vcs/pkg/git"
)

func main() {

	args := os.Args[1:]

	if args[0] != "init" {
		s, _ := os.Getwd()
		pkg.SetWorkingDirPath(s)
		err := os.Chdir(".vcs")
		if err != nil {
			fmt.Print(cmd.NotInitialisedMessage)
			return
		}
		s, _ = os.Getwd()
		pkg.SetVCSDirPath(s)
	}
	// 				TODO:: ERROR: if .vcs folder was modified manually, which does not match the correct way.
	switch args[0] {
	case "init":
		git.Init(args) // TODO: What happens if we reinitialise the repository??
		var InitialisedMessage string = fmt.Sprintf("Initialized empty VCS repository in %s\n", pkg.WorkingDirPath)
		fmt.Print(InitialisedMessage)
	case "config":
		if len(args) == 2 {
			val, _ := git.ParseConfigData(strings.Split(args[1], ".")[0], strings.Split(args[1], ".")[1])
			fmt.Println(val)
		} else {
			git.Config(args)
		}
	case "commit":
		git.Commit()
	case "branch":
		if len(args) >= 2 {
			git.CreateBranch(args[1:])
		} else {
			git.ListBranches()
		}
	case "checkout":
		git.Checkout(args[1:])
	case "add":
		git.Add(args[1:])
	case "diff":
		git.Diff(args[1:])
	case "cat-file":
		git.Catfile(args[1], args[2])
	case "log":
		git.Log(args[1:])
	case "status":
		git.Status()
	case "reset":
		git.Reset(args[1:])
	case "merge":
		git.Merge(args[1])
	}
}
