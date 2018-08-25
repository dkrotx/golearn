package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func filterDirs(files []os.FileInfo, printFiles bool) []os.FileInfo {
	if printFiles {
		return files
	}

	var filtered []os.FileInfo
	for _, file := range files {
		if file.IsDir() {
			filtered = append(filtered, file)
		}
	}

	return filtered
}

func humanizeSize(size int64) string {
	if size == 0 {
		return "empty"
	}
	return fmt.Sprintf("%db", size)
}

func printEntryWithSize(out io.Writer, file os.FileInfo) {
	fmt.Fprint(out, file.Name())
	if !file.IsDir() {
		fmt.Fprintf(out, " (%s)", humanizeSize(file.Size()))
	}
	fmt.Fprintln(out)
}

func dirTreeImpl(out io.Writer, path string, printFiles bool, prefix string) (err error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	filtered := filterDirs(files, printFiles)


	for i, file := range filtered {
		isLast := i == len(filtered) - 1

		fmt.Fprint(out, prefix)
		if isLast {
			fmt.Fprint(out, "└───")
		}  else {
			fmt.Fprint(out, "├───")
		}
		printEntryWithSize(out, file)

		if file.IsDir() {
			dirPath := filepath.Join(path, file.Name())

			var childPrefix string
			if isLast {
				childPrefix = prefix + "\t"
			} else {
				childPrefix = prefix + "│\t"
			}

			if err = dirTreeImpl(out, dirPath, printFiles, childPrefix); err != nil {
				return
			}
		}
	}

	return
}

func dirTree(out io.Writer, path string, printFiles bool) (err error) {
	return dirTreeImpl(out, path, printFiles, "")
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
