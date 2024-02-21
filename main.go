package main

import (
	"os"
	"path/filepath"
    "strconv"
    "time"

	"github.com/rs/zerolog"
)

type Dependency struct {
	name   string
	direct bool
}

type Module struct {
	name         string
	path         string
	dependencies []Dependency
    imports []string
}

func main() {
	loglevel, err := strconv.Atoi(os.Getenv("LOGLEVEL"))
	if err != nil {
		loglevel = int(zerolog.InfoLevel)
	}

	logger := zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}).Level(zerolog.Level(loglevel)).With().Timestamp().Caller().Logger()
	logger.Trace().Msg("Starting up...")

	path := "./"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
    logger.Trace().Msgf("Using path: %s", path)

	modulePaths, err := readGoWork(path)
	if err != nil {
		panic(err)
	}

    modules := make(map[string]Module)
    modulesByPath := make(map[string]Module)

	for _, modulePath := range modulePaths {
		moduleName, err := readGoModuleName(path, modulePath)
		if err != nil {
			panic(err)
		}

        module := Module{
            name: moduleName,
            path: modulePath,
        }

        modules[moduleName] = module
        modulesByPath[modulePath] = module

        logger.Debug().Msgf("Module: %s -> %s", moduleName, modulePath)
	}

	for _, module := range modules {
        logger.Debug().Msgf("Checking module: %s", module.name)

		imports, err := readImports(path, module.path)
		if err != nil {
			panic(err)
		}

        module.imports = imports
		logger.Trace().Msgf("Imports: %v", imports)

        // Find all imports that are in the modules list
        module.dependencies = make([]Dependency, 0)
        for _, imp := range imports {
            if _, ok := modules[imp]; ok {
                module.dependencies = append(module.dependencies, Dependency{
                    name: imp,
                    direct: true,
                })
            }
        }

        if len(module.dependencies) > 0 {
            logger.Debug().Msgf("Dependencies: %v", module.dependencies)
        }

        modules[module.name] = module
    }

    logger.Debug().Msg("Sorting modules...")
    sorted := sortModules(modules)
	for _, module := range sorted {
        if len(module.dependencies) == 0 {
            logger.Trace().Msgf("No dependencies for module: %s", module.name)
            continue
        }

        logger.Debug().Msgf("Module: %s", module.name)
        for _, dep := range module.dependencies {
            logger.Debug().Msgf("  Dependency: %s (direct: %v)", dep.name, dep.direct)
        }

		replacements, err := readGoModReplacements(path, module.path)
		if err != nil {
			panic(err)
		}
		logger.Trace().Msgf("Replacements: %v", replacements)

		// Find all imports that are not in the replacements map
		missingImports := make([]Dependency, 0)
		for _, imp := range module.dependencies {
			if _, ok := replacements[imp.name]; !ok {
				missingImports = append(missingImports, imp)
			}
		}
		logger.Trace().Msgf("Missing Imports: %v", missingImports)

		// Find all missing imports that exist in the moduleNames list
		missingModules := make([]Dependency, 0)
		for _, missingImport := range missingImports {
			if _, ok := modules[missingImport.name]; ok {
				missingModules = append(missingModules, missingImport)
			}
		}
		logger.Trace().Msgf("Missing Modules: %v", missingModules)

		// Add replacements for all missing modules
		for _, missingModule := range missingModules {
			actualModulePath := modules[missingModule.name].path

			var moduleAbsPath string
			moduleAbsPath, err = filepath.Abs(filepath.Join(path, module.path))
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

			logger.Debug().Msgf("Adding replacement for: %s -> %s", missingModule.name, relPath)
			err = addGoModReplacement(path, module.path, missingModule.name, relPath, missingModule.direct)
			if err != nil {
				panic(err)
			}
		}
	}
}

