package core

import (
	"encoding/json"
)

func GetDependencies() ([]Dependency, error) {
	content, err := findAndReadPackageJson()
	if err != nil {
		return nil, err
	}

	var packageJson Package
	err = json.Unmarshal(content, &packageJson)
	if err != nil {
		return nil, err
	}

	var deps []Dependency
	for dep, version := range packageJson.Dependencies {
		deps = append(deps, Dependency{Name: dep, RequiredVersion: version})
	}

	return deps, nil
}
