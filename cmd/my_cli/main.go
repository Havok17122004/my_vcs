package main

import (
	"os"
	"vcs/pkg/git"
)

func main() {

	args := os.Args[1:]

	switch args[0] {
	case "init":
		git.Init()
	}
}
