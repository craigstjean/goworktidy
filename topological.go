package main

func nextUnmarked(modules map[string]Module, marked map[string]bool) string {
	for module, _ := range modules {
		if _, ok := marked[module]; !ok {
			return module
		}
	}

	return ""
}

func visit(module string, modules map[string]Module, temporaryMark map[string]bool, permanentMark map[string]bool, ordered []Module) []Module {
	if _, ok := permanentMark[module]; ok {
		return ordered
	}

	if _, ok := temporaryMark[module]; ok {
		return nil
	}

	temporaryMark[module] = true

	for _, dependency := range modules[module].dependencies {
		r := visit(dependency.name, modules, temporaryMark, permanentMark, ordered)
		if r == nil {
			return nil
		}

		ordered = r
	}

	delete(temporaryMark, module)
	permanentMark[module] = true

	return append(ordered, modules[module])
}

func appendDependencies(dependencies []Dependency, modules map[string]Module) []Dependency {
	result := make([]Dependency, 0)
	seen := make(map[string]bool)
	for _, dependency := range dependencies {
		result = append(result, dependency)
		seen[dependency.name] = true

		subDependencies := appendDependencies(modules[dependency.name].dependencies, modules)
		for _, subDependency := range subDependencies {
			if _, ok := seen[subDependency.name]; !ok {
				subDependency.direct = false
				result = append(result, subDependency)
				seen[subDependency.name] = true
			}
		}
	}

	return result
}

func sortModules(modules map[string]Module) []Module {
	ordered := make([]Module, 0)
	permanentMarks := make(map[string]bool)
	temporaryMarks := make(map[string]bool)
	for len(permanentMarks) < len(modules) {
		m := nextUnmarked(modules, permanentMarks)
		r := visit(m, modules, temporaryMarks, permanentMarks, ordered)

		if r == nil {
			return nil
		}

		ordered = r
	}

	for i, module := range ordered {
		if len(module.dependencies) > 0 {
            ordered[i].dependencies = appendDependencies(module.dependencies, modules)
		}
	}

	return ordered
}
