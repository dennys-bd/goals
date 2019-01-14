package core

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	errs "github.com/dennys-bd/goals/shortcuts/errors"
)

func loadDotEnv(project Project) {
	file, err := os.Open(filepath.Join(project.AbsPath, ".env"))
	if err != nil && !os.IsNotExist(err) {
		errs.Ex(err.Error())
	} else {
		defer file.Close()
		if err == nil {
			scanner := bufio.NewScanner(file)
			line := 1

			for scanner.Scan() {
				s := strings.SplitN(scanner.Text(), "=", 2)
				if len(s) == 2 {
					os.Setenv(s[0], s[1])
				} else {
					errs.Ex(fmt.Sprintf("Syntax error in .env file. Line %d: \"%s\"\n", line, scanner.Text()))
				}
				line++
			}
		}
	}
}
