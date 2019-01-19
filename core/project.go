package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	errs "github.com/dennys-bd/goals/shortcuts/errors"
	"github.com/rs/cors"
)

// Project model
type Project struct {
	AbsPath    string `toml:"abs_path"`
	ImportPath string `toml:"import_path"`
	Name       string
	GoVersion  string `toml:"go_version"`
	AppMode    string `toml:"app_mode"`
	Config     config `toml:"config"`
}

type config struct {
	Port     int  `toml:"port"`
	Graphiql bool `toml:"graphilql"`
	verbose  bool `toml:"-"`
}

// CreateGoalsToml create the file Goals.Toml
// in which we save some of the project attributes
func (p Project) CreateGoalsToml() string {
	return fmt.Sprintf(`[project]
	name = "%s"
	import_path = "%s"
	go_version = "%s"
	app_mode = "%s"
	
[project.config]
	port = 8080`, p.Name, p.ImportPath, p.GoVersion, p.AppMode)
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

//ConfigPath is the path to package lib
func (p Project) ConfigPath() string {
	if p.AbsPath == "" {
		return ""
	}
	return filepath.Join(p.AbsPath, "config")
}

func (p Project) LoadDotEnv() {
	loadDotEnv(p)
}

type goalsToml struct {
	Project Project
	Cors    *cors.Options
}

func recreateProjectFromGoals() (Project, *cors.Options) {
	wd, err := os.Getwd()
	errs.CheckEx(err)

	data, err := ioutil.ReadFile(filepath.Join(wd, "config/Goals.toml"))
	if err != nil {
		errs.Ex("This is not a goals project")
	}

	p, c, err := recreateProject(string(data))
	errs.CheckEx(err)

	p.AbsPath = wd

	return p, c
}

func recreateProject(projectString string) (Project, *cors.Options, error) {
	var m goalsToml
	_, err := toml.Decode(projectString, &m)
	return m.Project, m.Cors, err
}
