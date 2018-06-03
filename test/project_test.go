package test

import (
	"os"
	"testing"

	"github.com/dennys-bd/goals/cmd"
)

type ProjectTest struct {
	folder      string
	projectName string
	result      cmd.Project
}

var projectTest = []ProjectTest{
	{".", "letest/goals", cmd.Project{
		AbsPath:    os.Getenv("GOPATH") + "/src/github.com/dennys-bd/goals/test/letest/goals",
		ImportPath: "github.com/dennys-bd/goals/test/letest/goals",
		Name:       "goals"},
	},
	{os.Getenv("GOPATH") + "/src", "legoals", cmd.Project{
		AbsPath:    os.Getenv("GOPATH") + "/src/legoals",
		ImportPath: "legoals",
		Name:       "legoals"},
	},
	{os.Getenv("HOME") + "/Desktop", "github.com/goals/outhergoals", cmd.Project{
		AbsPath:    os.Getenv("GOPATH") + "/src/github.com/goals/outhergoals",
		ImportPath: "github.com/goals/outhergoals",
		Name:       "outhergoals"},
	},
	{os.Getenv("HOME"), "onemoregoals", cmd.Project{
		AbsPath:    os.Getenv("GOPATH") + "/src/onemoregoals",
		ImportPath: "onemoregoals",
		Name:       "onemoregoals"},
	},
}

func TestNewProject(t *testing.T) {
	for _, test := range projectTest {
		err := os.Chdir(test.folder)
		if err != nil {
			t.Error(err)
		}
		project := cmd.NewProject(test.projectName)
		if project.Name != test.result.Name {
			t.Errorf(`NewProject(%s) in %s, Name wanted "%s", got "%s"`,
				test.projectName, test.folder, test.result.Name, project.Name)
		}
		if project.AbsPath != test.result.AbsPath {
			t.Errorf(`NewProject(%s) in %s, AbsPath wanted "%s", got "%s"`,
				test.projectName, test.folder, test.result.AbsPath, project.AbsPath)
		}
		if project.ImportPath != test.result.ImportPath {
			t.Errorf(`NewProject(%s) in %s, ImportPath wanted "%s", got "%s"`,
				test.projectName, test.folder, test.result.ImportPath, project.ImportPath)
		}
	}
}

type RecreateProjectTest struct {
	projectString string
	result        cmd.Project
}

var recreateProjectTest = []RecreateProjectTest{
	{`[project]
	name = "goals"
	import_path = "github.com/dennys-bd/goals"
	go_version = "go1.8"
	app_mode = "gateway"`,
		cmd.Project{Name: "goals", ImportPath: "github.com/dennys-bd/goals",
			GoVersion: "go1.8", AppMode: "gateway"},
	},
	{`[project]
	name = "othergoals"
	import_path = "othergoals"
	go_version = "go1.10"
	app_mode = "webapp"`,
		cmd.Project{Name: "othergoals", ImportPath: "othergoals",
			GoVersion: "go1.10", AppMode: "webapp"},
	},
}

func TestRecriateProject(t *testing.T) {
	for _, test := range recreateProjectTest {
		project, err := cmd.RecreateProject(test.projectString)
		if err != nil {
			t.Error(err)
		}
		if project.Name != test.result.Name {
			t.Errorf(`RecreateProject(%s), Name wanted "%s", got "%s"`,
				test.projectString, test.result.Name, project.Name)
		}
		if project.GoVersion != test.result.GoVersion {
			t.Errorf(`RecreateProject(%s), GoVersion wanted "%s", got "%s"`,
				test.projectString, test.result.GoVersion, project.GoVersion)
		}
		if project.ImportPath != test.result.ImportPath {
			t.Errorf(`RecreateProject(%s), ImportPath wanted "%s", got "%s"`,
				test.projectString, test.result.ImportPath, project.ImportPath)
		}
		if project.AppMode != test.result.AppMode {
			t.Errorf(`RecreateProject(%s), AppMode wanted "%s", got "%s"`,
				test.projectString, test.result.AppMode, project.AppMode)
		}
	}
}

func TestCreateGoalsToml(t *testing.T) {
	for _, test := range recreateProjectTest {
		str := test.result.CreateGoalsToml()
		if str != test.projectString {
			t.Errorf(`CreateGoalsToml() of %s,
String wanted:
%s 
got:
%s`,
				test.result.Name, test.projectString, str)
		}
	}
}
