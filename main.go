package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	path := "./"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	modules, err := readGoWork(path)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Modules: %v\n", modules)

	moduleNames := make([]string, 0)
	modulePaths := make(map[string]string)
	for _, module := range modules {
		moduleName, err := readGoModuleName(path, module)
		if err != nil {
			panic(err)
		}

		moduleNames = append(moduleNames, moduleName)
		modulePaths[moduleName] = module
	}
	fmt.Printf("Module Names: %v\n", moduleNames)
	fmt.Printf("Module Paths: %v\n", modulePaths)

	for _, module := range modules {
		fmt.Println()
		fmt.Printf("Checking module: %s\n", module)
		imports, err := readImports(path, module)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Imports: %v\n", imports)

		replacements, err := readGoModReplacements(path, module)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Replacements: %v\n", replacements)

		// Find all imports that are not in the replacements map
		missingImports := make([]string, 0)
		for _, imp := range imports {
			if _, ok := replacements[imp]; !ok {
				missingImports = append(missingImports, imp)
			}
		}
		fmt.Printf("Missing Imports: %v\n", missingImports)

		// Find all missing imports that exist in the moduleNames list
		missingModules := make([]string, 0)
		for _, missingImport := range missingImports {
			if _, ok := modulePaths[missingImport]; ok {
				missingModules = append(missingModules, missingImport)
			}
		}
		fmt.Printf("Missing Modules: %v\n", missingModules)

		// Add replacements for all missing modules
		for _, missingModule := range missingModules {
			actualModulePath := modulePaths[missingModule]

			var moduleAbsPath string
			moduleAbsPath, err = filepath.Abs(filepath.Join(path, module))
			if err != nil {
				panic(err)
			}

			var missingAbsPath string
			missingAbsPath, err = filepath.Abs(filepath.Join(path, actualModulePath))
			if err != nil {
				panic(err)
			}

			var relPath string
			relPath, err = filepath.Rel(moduleAbsPath, missingAbsPath)
			if err != nil {
				panic(err)
			}

			fmt.Printf("Adding replacement for: %s -> %s\n", missingModule, relPath)
			err = addGoModReplacement(path, module, missingModule, relPath)
			if err != nil {
				panic(err)
			}
		}
	}
}
