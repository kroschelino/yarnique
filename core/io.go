package core

import (
	"errors"
	"os"
)

func writeFile(filename string, data []byte) error {
	err := os.WriteFile(filename, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

func readFile(filename string) ([]byte, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return b, err
	}
	return b, nil
}

func findAndReadPackageJson() ([]byte, error) {
	content, err := readFile("package.json")
	if err != nil {
		return nil, errors.New("Could not locate the `package.json`")
	}
	return content, nil
}

func findAndReadYarnLock() ([]byte, error) {
	content, err := readFile("yarn.lock")
	if err != nil {
		return nil, errors.New("Could not locate the `yarn.lock`, you may have to generate one via `yarn`")
	}
	return content, nil
}
