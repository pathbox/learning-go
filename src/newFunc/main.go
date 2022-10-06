package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func main() {
	err := makeConstructor()
	if err != nil {
		fmt.Printf("newc: [ERROR] %v\n", err)
		fmt.Printf("It seems like there is some trouble here. Try this:\n")
		fmt.Printf("\t1. Check and upgrade this tool (https://github.com/Bin-Huang/newc)\n")
		fmt.Printf("\t2. Submit an issue on Github (https://github.com/Bin-Huang/newc/issues)\n")
		os.Exit(1)
	}
}

func makeConstructor() error {
	pkg, err := GetPackageInfo(".")
	if err != nil {
		return err
	}
	// skip if generated recently
	genFilename, err := filepath.Abs("./constructor_gen.go")
	if err != nil {
		return err
	}
	if isGeneratedRecently(genFilename) {
		return nil
	}
	allImports := []ImportInfo{}
	allResults := []StructInfo{}
	for _, filename := range pkg.GoFiles {
		has, err := IncludeMakeMark(filename)
		if err != nil {
			return err
		}
		if !has {
			continue
		}
		results, imports, err := ParseCodeFile(filename)
		if err != nil {
			return err
		}
		if len(results) == 0 {
			continue
		}
		allImports = append(allImports, imports...)
		allResults = append(allResults, results...)
	}
	if len(allResults) == 0 {
		return nil
	}
	code, err := GenerateCode(pkg.Name, allImports, allResults)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(genFilename, []byte(code), 0644)
	if err != nil {
		return err
	}
	fmt.Printf("newc: [INFO] wrote %v\n", genFilename)
	return nil
}

func isGeneratedRecently(genFilename string) bool {
	stat, err := os.Stat(genFilename)
	if err != nil {
		return false
	}
	return time.Now().Sub(stat.ModTime()) < 5*time.Second
}
