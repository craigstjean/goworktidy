package main

import (
	"golang.org/x/mod/modfile"
	"os"
	"path/filepath"
)

func readGoWork(path string) ([]string, error) {
	fullpath := filepath.Join(path, "go.work")
	_, err := os.Stat(fullpath)
	if err != nil {
		return nil, err
	}

	var contents []byte
	contents, err = os.ReadFile(fullpath)
	if err != nil {
		return nil, err
	}

	var work *modfile.WorkFile
	work, err = modfile.ParseWork(fullpath, contents, nil)
	if err != nil {
		return nil, err
	}

	var modules []string
	for _, use := range work.Use {
		modules = append(modules, use.Path)
	}

	return modules, nil
}
