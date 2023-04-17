package core

import (
	"errors"
	"fmt"
	"sort"

	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

var depTree DependencyTree

func findItemInTree(toFind *DependencyTreeItem) bool {
	return slices.IndexFunc(depTree.Items, func(depTreeItem DependencyTreeItem) bool {
		if depTreeItem.Item.Name == toFind.Item.Name && depTreeItem.Item.UsedVersion == toFind.Item.UsedVersion {
			if depTreeItem.Parent != nil && toFind.Parent != nil {
				return depTreeItem.Parent.Name == toFind.Parent.Name && depTreeItem.Parent.UsedVersion == toFind.Parent.UsedVersion
			}
			return false
		}
		return false
	}) != -1
}

func buildDependencyTree(rootDeps []Dependency, yarnLock string) {

	fmt.Print("Creating dependency tree...")
	var buildDependencyTreeRecursive func(deps []Dependency, parentDep *Dependency, yarnLock *string)

	buildDependencyTreeRecursive = func(deps []Dependency, parentDep *Dependency, yarnLock *string) {
		for _, dep := range deps {
			usedVersion, ok := versionMap[dep.Name][dep.RequiredVersion]
			if !ok {
				continue
			}

			depTreeItem := DependencyTreeItem{Item: dep}
			depTreeItem.Item.UsedVersion = usedVersion

			if parentDep != nil {
				depTreeItem.Parent = parentDep
			}

			if !findItemInTree(&depTreeItem) {

				depTree.Items = append(depTree.Items, depTreeItem)

				childDeps := findDependencies(&depTreeItem.Item, yarnLock)
				if childDeps != nil {
					buildDependencyTreeRecursive(childDeps, &depTreeItem.Item, yarnLock)
				}
			}
		}
	}

	buildDependencyTreeRecursive(rootDeps, nil, &yarnLock)

	fmt.Println("done")
}

func sortDependencyTree(depTree *DependencyTree) {
	sort.Slice(depTree.Items, func(i, j int) bool {
		return depTree.Items[i].Item.Name < depTree.Items[j].Item.Name
	})
}

func BuildDependencyTree() error {
	yarnLock, err := findAndReadYarnLock()
	if err != nil {
		return err
	}
	deps, err := GetDependencies()
	if err != nil {
		return err
	}

	createVersionsMap(string(yarnLock))
	buildDependencyTree(deps, string(yarnLock))
	sortDependencyTree(&depTree)

	err = writeToYaml(&depTree)
	return err
}

func writeToYaml(tree *DependencyTree) error {
	fmt.Print("Saving output to `.yarnique.yml`...")
	d, err := yaml.Marshal(&tree)
	if err != nil {
		return errors.New("Could not create dependency tree")
	}

	err = writeFile(".yarnique.yml", d)
	if err != nil {
		return errors.New("Could not create dependency tree")
	}
	fmt.Println("done")
	return nil
}
