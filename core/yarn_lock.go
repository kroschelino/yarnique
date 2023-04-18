package core

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

var versionMap map[string]map[string]string

func unifyString(input string) string {
	return strings.Trim(strings.TrimSpace(input), `"`)
}

func escape(version string) string {
	re := regexp.MustCompile("[.~<>^/]")
	return re.ReplaceAllString(version, "\\$0")
}

func createVersionsMap(yarnLock string) {
	fmt.Print("Creating version map for packages...")
	versionMap = make(map[string]map[string]string)
	packageRegex := regexp.MustCompile(`(?m)^(\S.*):\s*  version "(\S+)"`)
	if packageInfo := packageRegex.FindAllStringSubmatch(yarnLock, -1); packageInfo != nil {
		for _, match := range packageInfo {
			if len(match) == 3 {
				usedVersion := match[2]
				requiredVersions := strings.Split(match[1], ",")
				for _, packageRawData := range requiredVersions {
					packageRawData = unifyString(packageRawData)
					packageRegex := regexp.MustCompile(`^(\S+)@(.+)$`)
					if packageData := packageRegex.FindStringSubmatch(packageRawData); len(packageData) == 3 {
						packageName := packageData[1]
						requiredVersion := packageData[2]
						if versionMap[packageName] == nil {
							versionMap[packageName] = make(map[string]string)
						}
						versionMap[packageName][requiredVersion] = usedVersion
					}
				}

			}
		}
	}
	fmt.Println("done")
}

func findDirectDependencies(data *string) []Dependency {
	var result []Dependency
	depRegex := regexp.MustCompile(`(?m)  dependencies:\r?\n+?((?:    \S+ .+?\r?\n)+)`)
	if deps := depRegex.FindStringSubmatch(*data); len(deps) > 1 {
		scanner := bufio.NewScanner(strings.NewReader(deps[1]))
		for scanner.Scan() {
			line := unifyString(scanner.Text())
			if splitIndex := strings.Index(line, " "); splitIndex != -1 {
				name := line[:splitIndex]
				version := line[splitIndex:]
				result = append(result, Dependency{Name: unifyString(name), RequiredVersion: unifyString(version)})
			}

		}
	}
	return result
}

func findOptionalDependencies(data *string) []Dependency {
	var result []Dependency
	depRegex := regexp.MustCompile(`(?m)  optionalDependencies:\r?\n+?((?:    \S+ .+?\r?\n)+)`)
	if deps := depRegex.FindStringSubmatch(*data); len(deps) > 1 {
		scanner := bufio.NewScanner(strings.NewReader(deps[1]))
		for scanner.Scan() {
			line := unifyString(scanner.Text())
			if splitIndex := strings.Index(line, " "); splitIndex != -1 {
				name := line[:splitIndex]
				version := line[splitIndex:]
				result = append(result, Dependency{Name: unifyString(name), RequiredVersion: unifyString(version)})
			}
		}
	}
	return result
}

func findDependencies(parent *Dependency, yarnLock *string) []Dependency {
	libName := fmt.Sprintf("%s@%s", escape(unifyString(parent.Name)), escape(parent.RequiredVersion))
	sectionRegex := regexp.MustCompile(fmt.Sprintf(`(?ms)"?%s"?(.*?\r?\n)\r?\n`, libName))
	if section := sectionRegex.FindStringSubmatch(*yarnLock); len(section) > 1 {
		result := findDirectDependencies(&section[1])
		result = append(result, findOptionalDependencies(&section[1])...)
		return result
	}

	return nil
}
