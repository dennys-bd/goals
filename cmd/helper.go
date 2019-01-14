package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"

	errs "github.com/dennys-bd/goals/shortcuts/errors"
	oss "github.com/dennys-bd/goals/shortcuts/os"
)

var srcPaths []string

func init() {
	// Initialize srcPaths.
	envGoPath := os.Getenv("GOPATH")
	goPaths := filepath.SplitList(envGoPath)
	if len(goPaths) == 0 {
		// Adapted from https://github.com/Masterminds/glide/pull/798/files.
		// As of Go 1.8 the GOPATH is no longer required to be set. Instead there
		// is a default value. If there is no GOPATH check for the default value.
		// Note, checking the GOPATH first to avoid invoking the go toolchain if
		// possible.

		out, err := exec.Command("go", "env", "GOPATH").Output()
		errs.CheckEx(err)

		toolchainGoPath := strings.TrimSpace(string(out))
		goPaths = filepath.SplitList(toolchainGoPath)
		if len(goPaths) == 0 {
			errs.Ex("$GOPATH is not set")
		}
	}
	srcPaths = make([]string, 0, len(goPaths))
	for _, goPath := range goPaths {
		srcPaths = append(srcPaths, filepath.Join(goPath, "src"))
	}
}

func executeTemplate(tmplStr string, data interface{}) string {
	tmpl, err := template.New("").Funcs(template.FuncMap{"comment": commentifyString}).Parse(tmplStr)
	if err != nil {
		errs.Ex(err)
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, data)
	if err != nil {
		errs.Ex(err)
	}
	return buf.String()
}

func replaceTemplate(tmplStr string, data map[string]string) string {
	for k, v := range data {
		tmplStr = strings.Replace(tmplStr, fmt.Sprintf("{{.%s}}", k), v, -1)
	}
	return tmplStr
}

func commentifyString(in string) string {
	var newlines []string
	lines := strings.Split(in, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "//") {
			newlines = append(newlines, line)
		} else {
			if line == "" {
				newlines = append(newlines, "//")
			} else {
				newlines = append(newlines, "// "+line)
			}
		}
	}
	return strings.Join(newlines, "\n")
}

func writeToFile(path string, r io.Reader) error {
	if oss.Exists(path) {
		return fmt.Errorf("%v already exists", path)
	}

	dir := filepath.Dir(path)
	if dir != "" {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return err
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	return err
}

func writeStringToFile(path string, s string) {
	err := writeToFile(path, strings.NewReader(s))
	if err != nil {
		errs.Ex(err)
	}
}

func removeFile(path string) error {
	if !oss.Exists(path) {
		return fmt.Errorf("%v doesnt exists", path)
	}
	err := os.Remove(path)
	return err
}

func toSnake(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}

func toAbbreviation(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	out = append(out, unicode.ToLower(runes[0]))
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, unicode.ToLower(runes[i]))
		}
	}

	return string(out)
}

func getGoVersion() string {
	v, err := exec.Command("go", "version").Output()
	errs.CheckEx(err)
	vl := strings.Split(string(v), " ")
	return vl[2]
}
