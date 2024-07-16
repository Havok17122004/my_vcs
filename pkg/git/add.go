package git

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"vcs/pkg"
)

func makeBlobs(path string, index *pkg.Index) {
	dir, err := os.Open(filepath.Join(pkg.WorkingDirPath, path))
	pkg.Check(err)
	files, _ := dir.ReadDir(0)
	info, err := dir.Stat()
	pkg.Check(err)
	if !info.IsDir() {
		s := pkg.CompressFileStoreInObjects(filepath.Join(pkg.WorkingDirPath, path), "blob")
		pkg.Check(err)
		meta, err := pkg.GetFileMetadata(info)
		// fmt.Println(meta)
		pkg.Check(err)

		index.ModifyIndex(filepath.Join(pkg.WorkingDirPath, path), meta, s)
		// fmt.Print("modified?")
		return
	}
	for _, file := range files {
		file_path := filepath.Join(path, file.Name())
		if file.Name() == ".vcs" || file.Name() == ".git" {
			fmt.Printf("Skipped %s\n", file.Name())
			continue
		} else if file.IsDir() {
			makeBlobs(file_path, index)
		} else {
			s := pkg.CompressFileStoreInObjects(filepath.Join(pkg.WorkingDirPath, file_path), "blob")
			info, _ := file.Info()
			meta, err := pkg.GetFileMetadata(info)
			// fmt.Println(meta)
			pkg.Check(err)
			index.ModifyIndex(filepath.Join(pkg.WorkingDirPath, file_path), meta, s)
		}
	}
}

func Add(arg []string) {
	if len(arg) == 0 {
		fmt.Println("no file path passed as argument to add")
		return
	}
	index := pkg.ParseIndex()
	fmt.Println(index)
	sort.Strings(arg)

	for _, entry := range arg {
		makeBlobs(entry, index)
	}
	index.SaveIndex()
}
