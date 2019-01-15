package cmd

import (
	"os"
	"os/exec"
	"strconv"

	"github.com/dennys-bd/goals/core"
	"github.com/spf13/cobra"
)

var port string
var envPort bool
var environment string
var verbose bool

var runServerCmd = &cobra.Command{
	Use:     "runserver",
	Aliases: []string{"r"},
	Short:   "Runs your goals application",
	Run: func(cmd *cobra.Command, args []string) {
		project := recreateProjectFromGoals()

		runserver(project)
	},
}

func init() {
	runServerCmd.Flags().StringVarP(&port, "port", "p", "", "Set the port to your server.")
	runServerCmd.Flags().BoolVar(&envPort, "env-port", false, "Select the port from your environment variables (PORT)")
	runServerCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose prints more information about your Server.")
	runServerCmd.Flags().StringVarP(&environment, "env", "e", "", "Set the environments to start your Goals application.")
}

func runserver(project core.Project) {

	p := loadPort(project)
	loadDotEnv(project)

	cmd := exec.Command("go", "run", "server.go", p, strconv.FormatBool(verbose))
	cmd.Dir = project.AbsPath

	runCmd(cmd)
}

func loadPort(project core.Project) string {
	p := strconv.Itoa(project.Config.Port)

	if envPort {
		p = os.Getenv("PORT")
	}
	if port != "" {
		p = port
	}

	return p
}

func loadDotEnv(project core.Project) {
	project.LoadDotEnv()

	if os.Getenv("GOALS_ENV") == "" {
		os.Setenv("GOALS_ENV", "development")
	}

	if environment != "" {
		os.Setenv("GOALS_ENV", environment)
	}
}
