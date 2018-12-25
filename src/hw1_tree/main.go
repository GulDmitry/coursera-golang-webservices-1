package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
)

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

func dirTree(out io.Writer, path string, printFiles bool) error {
	err := drawFileStructureRecursively(out, path, printFiles, "")
	if err != nil {
		return err
	}

	return nil
}

func drawFileStructureRecursively(out io.Writer, path string, printFiles bool, prefix string) error {
	entrypointDir, err := os.Open(path)
	if err != nil {
		return err
	}
	defer closeResource(entrypointDir)

	files, err := entrypointDir.Readdir(-1)
	if err != nil {
		return err
	}

	if !printFiles {
		files = filterFiles(files)
	}
	sortFiles(files)

	for fileIndex := 0; fileIndex < len(files); fileIndex++ {
		file := files[fileIndex]
		currentPrefix := prefix + "├───"
		nextPrefix := prefix + "│\t"

		// Last element.
		if len(files)-1 == fileIndex {
			currentPrefix = prefix + "└───"
			nextPrefix = prefix + "\t"
		}

		switch file.Mode().IsDir() {
		case true: // Directory.
			_, err := out.Write([]byte(currentPrefix + file.Name() + "\n"))
			if err != nil {
				return err
			}

			err = drawFileStructureRecursively(out, path+string(os.PathSeparator)+file.Name(), printFiles, nextPrefix)
			if err != nil {
				return err
			}
		case false: // File.
			var fileSizeStr string
			if file.Size() == 0 {
				fileSizeStr = "empty"
			} else {
				fileSizeStr = strconv.Itoa(int(file.Size())) + "b"
			}
			formattedName := fmt.Sprintf("%s%s (%s)\n", currentPrefix, file.Name(), fileSizeStr)

			_, err := out.Write([]byte(formattedName))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Sorts files alphabetically.
func sortFiles(files []os.FileInfo) {
	sort.Slice(files, func(i, j int) bool {
		if files[i].Name() != files[j].Name() {
			return files[i].Name() < files[j].Name()
		}

		return files[i].Name() < files[j].Name()
	})
}

// Leaves only directories.
func filterFiles(files []os.FileInfo) (res []os.FileInfo) {
	for _, file := range files {
		if file.Mode().IsDir() {
			res = append(res, file)
		}
	}
	return
}

func closeResource(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}
