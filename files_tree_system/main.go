package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var printFiles bool

func main() {
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}

	printFiles = len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(os.Args[1])
	if err != nil {
		panic(err)
	}
}

func dirTree(path string) interface{} {
	return filepath.Walk(path, visit)
}

func visit(path string, f os.FileInfo, err error) error {
	repeat := strings.Count(path, "/")
	if repeat > 1 {
		repeat--
	}

	switch true {
	case repeat == 0:
		fmt.Printf("├───%s\n", filepath.Base(path))
	case f.IsDir():
		fmt.Printf("%s%s%s\n", strings.Repeat("|   ", repeat), "├───", filepath.Base(path))
	case printFiles && f.Size() == 0:
		fmt.Printf("%s%s%s (empty)\n", strings.Repeat("|   ", repeat), "├───", filepath.Base(path))
	case printFiles && f.Size() != 0:
		fmt.Printf("%s%s%s (%s)\n", strings.Repeat("|   ", repeat), "├───", filepath.Base(path), fmt.Sprint(f.Size(), "b"))
	}
	return err
}