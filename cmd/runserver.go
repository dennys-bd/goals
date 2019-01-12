package cmd

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var port int
var environment string

var runServerCmd = &cobra.Command{
	Use:     "runserver",
	Aliases: []string{"r"},
	Short:   "Runs your goals application",
	Run: func(cmd *cobra.Command, args []string) {
		project := recreateProjectFromGoals()

		runserver(project)
	},
}

func runserver(project Project) {
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := exec.Command("go", "run", "server.go", "8080")
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
