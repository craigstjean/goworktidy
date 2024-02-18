package main

import (
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

func readImports(path string, module string) ([]string, error) {
	fset := token.NewFileSet()
	packages, err := parser.ParseDir(fset, filepath.Join(path, module), nil, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}

	allImports := make([]string, 0)
	for _, pkg := range packages {
		for _, file := range pkg.Files {
			for _, imp := range file.Imports {
				allImports = append(allImports, strings.Replace(imp.Path.Value, "\"", "", -1))
			}
		}
	}

	return allImports, nil
}
