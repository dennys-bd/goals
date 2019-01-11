package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
)

type goalsToml struct {
	Project Project
}

// Project model
type Project struct {
	AbsPath    string `toml:"abs_path"`
	ImportPath string `toml:"import_path"`
	Name       string
	GoVersion  string `toml:"go_version"`
	AppMode    string `toml:"app_mode"`
}

// NewProject returns Project with specified project name.
func NewProject(projectName string) *Project {
	if projectName == "" {
		er("can't create project with blank name")
	}

	p := new(Project)

	p.GoVersion = getGoVersion()
	p.AppMode = "gateway"

	// 1. Find already created protect.
	p.AbsPath = findPackage(projectName)

	// 2. If there are no created project with this path, and user is in GOPATH/src,
	// then use working directory.
	if p.AbsPath == "" {
		wd, err := os.Getwd()
		check(err)

		for _, srcPath := range srcPaths {
			goPath := filepath.Dir(srcPath)
			if filepathHasPrefix(wd, goPath) {
				p.AbsPath = filepath.Join(wd, projectName)
				break
			}
		}
	}

	// 3. If user is not in GOPATH, then use (first GOPATH)/src/projectName.
	if p.AbsPath == "" {
		p.AbsPath = filepath.Join(srcPaths[0], projectName)
	}

	p.Name = filepath.Base(p.AbsPath)

	goPath := os.Getenv("GOPATH") + "/src/"
	p.ImportPath = strings.TrimPrefix(p.AbsPath, goPath)

	return p
}

// findPackage returns full path to existing go package in GOPATHs.
func findPackage(packageName string) string {
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

// recreateProjectFromGoals return the project configs from Goals.toml
func recreateProjectFromGoals() Project {
	wd, err := os.Getwd()
	check(err)

	data, err := ioutil.ReadFile(filepath.Join(wd, "lib/Goals.toml"))
	if err != nil {
		er("This is not a goals project")
	}
	check(err)

	p, err := RecreateProject(string(data))
	check(err)

	p.AbsPath = wd

	return p
}

// RecreateProject returns the project configs based on a Project String
func RecreateProject(projectString string) (Project, error) {
	var m goalsToml
	_, err := toml.Decode(projectString, &m)
	return m.Project, err
}

// CreateGoalsToml create the file Goals.Toml
// in which we save some of the project attributes
func (p Project) CreateGoalsToml() string {
	return fmt.Sprintf(`[project]
	name = "%s"
	import_path = "%s"
	go_version = "%s"
	app_mode = "%s"`, p.Name, p.ImportPath, p.GoVersion, p.AppMode)
}

//ResolverPath is the path to package resolver
func (p Project) ResolverPath() string {
	if p.AbsPath == "" {
		return ""
	}
	return filepath.Join(p.AbsPath, "app/resolver")
}

//ScalarPath is the path to package scalar
func (p Project) ScalarPath() string {
	if p.AbsPath == "" {
		return ""
	}
	return filepath.Join(p.AbsPath, "app/scalar")
}

//ModelPath is the path to package model
func (p Project) ModelPath() string {
	if p.AbsPath == "" {
		return ""
	}
	return filepath.Join(p.AbsPath, "app/model")
}

//SchemaPath is the path to package schema
func (p Project) SchemaPath() string {
	if p.AbsPath == "" {
		return ""
	}
	return filepath.Join(p.AbsPath, "app/schema")
}

//LibPath is the path to package lib
func (p Project) LibPath() string {
	if p.AbsPath == "" {
		return ""
	}
	return filepath.Join(p.AbsPath, "lib")
}
