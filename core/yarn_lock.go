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

func escapePackageVersion(version string) string {
	re := regexp.MustCompile("[.~<>^]")
	return re.ReplaceAllString(version, "\\$0")
}

func findVersion(data *string) string {
	versionRegex := regexp.MustCompile(`(?m)^  version "(\S+)"`)
	if actualVersion := versionRegex.FindStringSubmatch(*data); actualVersion != nil && len(actualVersion) > 1 {
		return actualVersion[1]
	}
	return ""
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
					if packageData := packageRegex.FindStringSubmatch(packageRawData); packageData != nil && len(packageData) == 3 {
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
	depRegex := regexp.MustCompile(`(?m)  dependencies:\r?\n+?((?:    \S+ \S+\r?\n)+)`)
	if deps := depRegex.FindStringSubmatch(*data); deps != nil && len(deps) > 1 {
		scanner := bufio.NewScanner(strings.NewReader(deps[1]))
		for scanner.Scan() {
			line := unifyString(scanner.Text())
			dep := strings.Split(line, " ")
			result = append(result, Dependency{Name: unifyString(dep[0]), RequiredVersion: unifyString(dep[1])})
		}
	}
	return result
}

func findOptionalDependencies(data *string) []Dependency {
	var result []Dependency
	depRegex := regexp.MustCompile(`(?m)  optionalDependencies:\r?\n+?((?:    \S+ \S+\r?\n)+)`)
	if deps := depRegex.FindStringSubmatch(*data); deps != nil && len(deps) > 1 {
		scanner := bufio.NewScanner(strings.NewReader(deps[1]))
		for scanner.Scan() {
			line := unifyString(scanner.Text())
			dep := strings.Split(line, " ")
			result = append(result, Dependency{Name: unifyString(dep[0]), RequiredVersion: unifyString(dep[1])})
		}
	}
	return result
}

func findDependencies(parent *Dependency, yarnLock *string) []Dependency {
	libName := fmt.Sprintf("%s@%s", unifyString(parent.Name), escapePackageVersion(parent.RequiredVersion))
	sectionRegex := regexp.MustCompile(fmt.Sprintf(`(?ms)"?%s"?(.*?\r?\n)\r?\n`, libName))
	if section := sectionRegex.FindStringSubmatch(*yarnLock); section != nil && len(section) > 1 {
		result := findDirectDependencies(&section[1])
		result = append(result, findOptionalDependencies(&section[1])...)
		return result
	}

	return nil
}
