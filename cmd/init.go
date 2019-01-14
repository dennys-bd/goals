package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/dennys-bd/goals/core"
	errs "github.com/dennys-bd/goals/shortcuts/errors"
	oss "github.com/dennys-bd/goals/shortcuts/os"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:     "init name",
	Aliases: []string{"initialize", "create"},
	Short:   "Initialize a Goals Application",
	Long: `Initialize (goals init) will create a new application, with 
the appropriate structure for a Go-Graphql application.

* If a name or relative path is provided, it will be created inside $GOPATH;
  (e.g. github.com/dennys-bd/goals);
* If an absolute path is provided, it will be created INSIDE $GOPATH;
* If your working directory is inside $GOPATH, it will be created on the wd;
* If the directory already exists but is empty, it will be used.

Init will not use an existing directory with contents.`,

	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Getwd()
		errs.Check(err)

		var project *core.Project
		if len(args) == 0 {
			errs.Ex("please insert the project name")
		} else if len(args) == 1 {
			project = newProject(args[0])
		} else {
			errs.Ex("please provide only one argument")
		}

		intializeProject(project)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		print("Creating your project")
		f := func() {
			for i := 0; i < 2; i++ {
				time.Sleep(500 * time.Millisecond)
				print(".")
			}
			time.Sleep(500 * time.Millisecond)
			println(".")
		}
		go f()
		go func() {
			time.Sleep(2000 * time.Millisecond)
			print("Downloading Depedences")
			f()
			println("It can take several minutes")
		}()
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		println("Done!")
	},
}

func intializeProject(project *core.Project) {
	if !oss.Exists(project.AbsPath) {
		err := os.MkdirAll(project.AbsPath, os.ModePerm)
		errs.CheckEx(err)
	} else if !oss.IsEmpty(project.AbsPath) {
		errs.Ex("Goals will not create a new project in a non empty directory: " + project.AbsPath)
	}

	initTemplates()
	basicTemplates()

	initializeDep(project)
	createStructure(project)
	createAbsFiles(project)
	// downloadDepedences(project)
	// TODO: createDatabase(project)
	// TODO: createDotEnv(project)
}

func initializeDep(project *core.Project) {
	cmd := exec.Command("dep", "init", project.AbsPath)
	err := cmd.Run()
	errs.CheckEx(err)

	str := `package main

	func main(){
	}
`
	writeStringToFile(filepath.Join(project.AbsPath, "main.go"), str)

	commands := []string{"dep ensure -add github.com/dennys-bd/goals"}

	for _, c := range commands {
		cs := strings.Split(c, " ")
		cmd := exec.Command(cs[0], cs[1:]...)
		cmd.Dir = project.AbsPath
		out, err := cmd.Output()
		errs.CheckEx(err)
		println(string(out))
	}
	removeFile(filepath.Join(project.AbsPath, "main.go"))
}

func createStructure(project *core.Project) {
	resData := map[string]interface{}{}
	resScript := executeTemplate(templates["resolver"], resData)

	schData := map[string]interface{}{"importpath": project.ImportPath}
	schScript := executeTemplate(templates["schema"], schData)

	writeStringToFile(filepath.Join(project.ResolverPath(), "resolver.go"), resScript)
	writeStringToFile(filepath.Join(project.SchemaPath(), "schema.go"), schScript)
}

func createAbsFiles(project *core.Project) {
	serverData := map[string]interface{}{"importpath": project.ImportPath}
	serverScript := executeTemplate(templates["server"], serverData)

	writeStringToFile(filepath.Join(project.AbsPath, "server.go"), serverScript)
	writeStringToFile(filepath.Join(project.AbsPath, ".gitignore"), templates["git"])
	writeStringToFile(filepath.Join(project.ConfigPath(), "Goals.toml"), project.CreateGoalsToml())
	writeStringToFile(filepath.Join(project.LibPath(), "consts.go"), templates["consts"])
	writeStringToFile(filepath.Join(project.ScalarPath(), "scalar.go"), templates["scalar"])
	writeStringToFile(filepath.Join(project.ScalarPath(), "json.go"), templates["json"])
	writeStringToFile(filepath.Join(project.ModelPath(), "helper.go"), templates["modelHelper"])
	writeStringToFile(filepath.Join(project.ResolverPath(), "helper.go"), templates["resolverHelper"])
}

func downloadDepedences(project *core.Project) {
	cmd := exec.Command("dep", "ensure")
	cmd.Dir = project.AbsPath
	err := cmd.Run()
	errs.CheckEx(err)
}
