package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/dennys-bd/goals/core"
	errs "github.com/dennys-bd/goals/shortcuts/errors"
	oss "github.com/dennys-bd/goals/shortcuts/os"
)

type goalsToml struct {
	Project core.Project
}

// newProject returns Project with specified project name.
func newProject(projectName string) *core.Project {
	if projectName == "" {
		errs.Ex("can't create project with blank name")
	}

	p := new(core.Project)

	p.GoVersion = getGoVersion()
	p.AppMode = "gateway"

	// 1. Find already created protect.
	p.AbsPath = findPackage(projectName)

	// 2. If there are no created project with this path, and user is in GOPATH/src,
	// then use working directory.
	if p.AbsPath == "" {
		wd, err := os.Getwd()
		errs.CheckEx(err)

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

// recreateProjectFromGoals return the project configs from Goals.toml
func recreateProjectFromGoals() core.Project {
	wd, err := os.Getwd()
	errs.CheckEx(err)

	data, err := ioutil.ReadFile(filepath.Join(wd, "config/Goals.toml"))
	if err != nil {
		errs.Ex("This is not a goals project")
	}

	p, err := recreateProject(string(data))
	errs.CheckEx(err)

	p.AbsPath = wd

	return p
}

// recreateProject returns the project configs based on a Project String
func recreateProject(projectString string) (core.Project, error) {
	var m goalsToml
	_, err := toml.Decode(projectString, &m)
	return m.Project, err
}

// findPackage returns full path to existing go package in GOPATHs.
func findPackage(packageName string) string {
	for _, srcPath := range srcPaths {
		packagePath := filepath.Join(srcPath, packageName)
		if oss.Exists(packagePath) {
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
