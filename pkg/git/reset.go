package git

import (
	"flag"
	"fmt"
	"path/filepath"
	"vcs/pkg"
)

func Reset(args []string) {
	flagSet := flag.NewFlagSet("reset", flag.ExitOnError)
	soft := flagSet.Bool("soft", false, "use the soft flag")
	hard := flagSet.Bool("hard", false, "use the hard flag")
	mixed := flagSet.Bool("mixed", true, "use the mixed flag")

	var hash string
	var err error

	// Parse the flags and arguments after "reset"
	flagSet.Parse(args[:])
	remainingArgs := flagSet.Args()
	var files []string
	if len(remainingArgs) == 0 {
		hash, _ = pkg.FetchHeadsSHAfromPath(pkg.ParseHEADgivePath())
		files = append(files, "")
	} else {
		hash, _, _, err = pkg.FindHashofCommit(remainingArgs[len(remainingArgs)-1])
		if err != nil {
			fmt.Println("here")
			hash, _ = pkg.FetchHeadsSHAfromPath(pkg.ParseHEADgivePath())
			files = remainingArgs
		} else {
			files = remainingArgs[:len(remainingArgs)-1]
		}
	}

	var flagValue string
	if *soft {
		flagValue = "soft"
	} else if *hard {
		flagValue = "hard"
		for _, path := range files {
			path = filepath.Join(pkg.WorkingDirPath, path)
			pkg.RecoverWorkingDirToCommitWithDeletions(path, hash)
			pkg.RecoverIndexToCommit(path, hash)
		}
	} else if *mixed {
		// Default flag
		for _, path := range files {
			path = filepath.Join(pkg.WorkingDirPath, path)
			// fmt.Println(path)
			pkg.RecoverIndexToCommit(path, hash)
		}
		flagValue = "mixed"
	}
	pkg.UpdateHeads(hash, pkg.ParseHEADgivePath())

	fmt.Println("Flag:", flagValue)
	fmt.Println("Files:", files)

}
