package core

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gookit/color"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

type DepTreeMap = map[string]map[string][]*Dependency
type DependencyPath = [][]string

var depTreeMap DepTreeMap

func ShowRootPackages(packageName string) error {
	err := getDepTreeFromYaml()
	if err != nil {
		return err
	}

	createMapFromDepTree()
	deps, err := getApplicableDependencies(packageName, &depTree)
	if err != nil {
		return err
	}
	printRootDependencies(deps, &depTreeMap)

	return nil
}

func packageName(dep *Dependency) string {
	return fmt.Sprintf("%s@%s", dep.Name, dep.UsedVersion)
}

func buildDepdencyPaths(dep *Dependency, depTreeMap *DepTreeMap) DependencyPath {
	var fullPath DependencyPath

	var buildDependencyPathRecursive func(dep *Dependency, fullPath *DependencyPath, path *[]string, recursionDepth int, depTreeMap *DepTreeMap)
	buildDependencyPathRecursive = func(dep *Dependency, fullPath *DependencyPath, path *[]string, recursionDepth int, depTreeMap *DepTreeMap) {

		result := append(*path, packageName(dep))

		if recursionDepth > 0 {

			if (*depTreeMap)[dep.Name][dep.UsedVersion] != nil {
				for _, parentDep := range (*depTreeMap)[dep.Name][dep.UsedVersion] {
					buildDependencyPathRecursive(parentDep, fullPath, &result, recursionDepth-1, depTreeMap)
				}
				return
			}
		}
		*fullPath = append(*fullPath, result)
	}

	recursionDepth := 20
	buildDependencyPathRecursive(dep, &fullPath, &[]string{}, recursionDepth, depTreeMap)
	return fullPath
}

func printRootDependencies(deps []*Dependency, depTreeMap *DepTreeMap) {

	for _, dep := range deps {

		leafPackage := packageName(dep)

		if (*depTreeMap)[dep.Name][dep.UsedVersion] == nil {
			// This is already a root dependency
			color.Style{color.FgGreen, color.OpBold}.Println(leafPackage)
			color.Println()
			continue
		}

		color.Style{color.White, color.OpBold}.Println(leafPackage)

		paths := buildDepdencyPaths(dep, depTreeMap)

		for _, path := range paths {

			color.Printf("\t> ")
			rootParent := path[len(path)-1]

			if len(path) >= 5 {
				color.Print("... > ")
			} else if len(path) > 2 {
				color.Printf("%s > ", strings.Join(path[1:len(path)-1], " > "))
			}
			color.Style{color.FgGreen, color.OpBold}.Println(rootParent)

		}
		color.Println()
	}
}

func getDepTreeFromYaml() error {
	b, err := readFile(".yarnique.yml")
	if err != nil {
		return errors.New("could not locate any `.yarnique.yml` in the current folder, did you run `yarnique build` first?")
	}
	if err = yaml.Unmarshal(b, &depTree); err != nil {
		return err
	}
	return nil
}

func getApplicableDependencies(packageName string, depTree *DependencyTree) ([]*Dependency, error) {

	var result []*Dependency
	var err error
	for i := range depTree.Items {

		if strings.Contains(depTree.Items[i].Item.Name, packageName) {
			if slices.IndexFunc(result, func(j *Dependency) bool {
				return j.Name == depTree.Items[i].Item.Name && j.UsedVersion == depTree.Items[i].Item.UsedVersion
			}) == -1 {
				result = append(result, &depTree.Items[i].Item)
			}
		}
	}

	if result == nil {
		err = errors.New("could not find any dependencies containing `packageName`")
	}

	return result, err
}

func ensureMapExists(sourceMap *DepTreeMap, key string) {
	if (*sourceMap)[key] == nil {
		(*sourceMap)[key] = make(map[string][]*Dependency)
	}
}

func createMapFromDepTree() {

	depTreeMap = make(DepTreeMap)

	for i := range depTree.Items {

		name := depTree.Items[i].Item.Name
		version := depTree.Items[i].Item.UsedVersion

		ensureMapExists(&depTreeMap, name)

		if depTree.Items[i].Parent != nil {
			depTreeMap[name][version] = append(depTreeMap[name][version], depTree.Items[i].Parent)
		} else {
			depTreeMap[name][version] = append(depTreeMap[name][version], nil)
		}

	}

}
