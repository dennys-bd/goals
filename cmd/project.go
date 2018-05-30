package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Project model
type Project struct {
	absPath string
	srcPath string
	name    string
}

// NewProject returns Project with specified project name.
func NewProject(projectName string) *Project {
	if projectName == "" {
		er("can't create project with blank name")
	}

	p := new(Project)
	p.name = projectName

	// 1. Find already created protect.
	p.absPath = findPackage(projectName)

	// 2. If there are no created project with this path, and user is in GOPATH,
	// then use GOPATH/src/projectName.
	if p.absPath == "" {
		wd, err := os.Getwd()
		check(err)

		for _, srcPath := range srcPaths {
			goPath := filepath.Dir(srcPath)
			if filepathHasPrefix(wd, goPath) {
				p.absPath = filepath.Join(srcPath, projectName)
				break
			}
		}
	}

	// 3. If user is not in GOPATH, then use (first GOPATH)/src/projectName.
	if p.absPath == "" {
		p.absPath = filepath.Join(srcPaths[0], projectName)
	}

	return p
}

// findPackage returns full path to existing go package in GOPATHs.
func findPackage(packageName string) string {
	if packageName == "" {
		return ""
	}

	for _, srcPath := range srcPaths {
		packagePath := filepath.Join(srcPath, packageName)
		if exists(packagePath) {
			return packagePath
		}
	}

	return ""
}

func filepathHasPrefix(path string, prefix string) bool {
	if len(path) <= len(prefix) {
		return false
	}
	if runtime.GOOS == "windows" {
		// Paths in windows are case-insensitive.
		return strings.EqualFold(path[0:len(prefix)], prefix)
	}
	return path[0:len(prefix)] == prefix

}

func (p Project) Name() string {
	return p.name
}
func (p Project) AbsPath() string {
	return p.absPath
}
func (p Project) GqlPath() string {
	if p.absPath == "" {
		return ""
	}
	return filepath.Join(p.absPath, "app/gqltype")
}
func (p Project) ResolverPath() string {
	if p.absPath == "" {
		return ""
	}
	return filepath.Join(p.absPath, "app/resolver")
}
func (p Project) ScalarPath() string {
	if p.absPath == "" {
		return ""
	}
	return filepath.Join(p.absPath, "app/scalar")
}
func (p Project) ModelPath() string {
	if p.absPath == "" {
		return ""
	}
	return filepath.Join(p.absPath, "app/model")
}
func (p Project) SchemaPath() string {
	if p.absPath == "" {
		return ""
	}
	return filepath.Join(p.absPath, "app/schema")
}
func (p Project) ImportPath() string {
	if p.absPath == "" {
		return ""
	}
	return filepath.Base(p.absPath)
}
func (p Project) LibPath() string {
	if p.absPath == "" {
		return ""
	}
	return filepath.Join(p.absPath, "lib")
}
