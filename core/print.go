package core

import "fmt"

func PrintDeps(deps []Dependency) {
	for _, x := range deps {
		fmt.Println(x)
	}
}
