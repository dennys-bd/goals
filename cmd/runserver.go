package cmd

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"

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
}

func runserver(project Project) {

	p := strconv.Itoa(project.Config.Port)

	if envPort {
		p = os.Getenv("PORT")
	}
	if port != "" {
		p = port
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := exec.Command("go", "run", "server.go", p, strconv.FormatBool(verbose))
	cmd.Dir = project.AbsPath

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	err := cmd.Start()

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()

	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()

	err = cmd.Wait()
	check(err)
	if errStdout != nil || errStderr != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}
}
