package json

import (
	"fmt"
	"os"
	"strings"

	"github.com/buger/jsonparser"
)

type Path struct {
	Paths []string
}

func FindInJson(filename string, libName string, version string) ([]Path, error) {
	jsonString, err := readItems(filename)
	if err != nil {
		return nil, err
	}

	var deps []Path

	jsonparser.ArrayEach(jsonString, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if value, err := jsonparser.GetString(value, "name"); err == nil && strings.HasPrefix(value, libName) {
			deps = append(deps, Path{Paths: []string{value}})
		}

	}, "data", "trees")

	return deps, nil
}

func PrintDeps(deps []Path) {
	for _, x := range deps {
		fmt.Println(x)
	}
}

func readItems(filename string) ([]byte, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return b, err
	}
	return b, nil
}
