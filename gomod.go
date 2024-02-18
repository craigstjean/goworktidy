package main

import (
	"golang.org/x/mod/modfile"
	"os"
	"path/filepath"
)

func readGoMod(path string, module string) (*modfile.File, error) {
	fullpath := filepath.Join(path, module, "go.mod")
	_, err := os.Stat(fullpath)
	if err != nil {
		return nil, err
	}

	var contents []byte
	contents, err = os.ReadFile(fullpath)
	if err != nil {
		return nil, err
	}

	var mod *modfile.File
	mod, err = modfile.Parse(fullpath, contents, nil)
	if err != nil {
		return nil, err
	}

	return mod, nil
}

func readGoModuleName(path string, module string) (string, error) {
	mod, err := readGoMod(path, module)
	if err != nil {
		return "", err
	}

	return mod.Module.Mod.Path, nil
}

func readGoModReplacements(path string, module string) (map[string]string, error) {
	mod, err := readGoMod(path, module)
	if err != nil {
		return nil, err
	}

	replacements := make(map[string]string)
	for _, replace := range mod.Replace {
		replacements[replace.Old.Path] = replace.New.Path
	}

	return replacements, nil
}

func addGoModReplacement(path string, module string, old string, new string) error {
	mod, err := readGoMod(path, module)
	if err != nil {
		return err
	}

	err = mod.AddReplace(old, "", new, "")
	if err != nil {
		return err
	}

	var data []byte
	data, err = mod.Format()
	if err != nil {
		return err
	}

	fullpath := filepath.Join(path, module, "go.mod")
	err = os.WriteFile(fullpath, data, 0644)
	return err
}
