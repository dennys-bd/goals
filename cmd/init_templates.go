package cmd

func initTemplates() {
	templates["server"] = `package main
	
import (
	"github.com/dennys-bd/goals/core"
	"{{.importpath}}/app/resolver"
	"{{.importpath}}/app/schema"
)

func main() {
	opts := core.GetOpts()

	core.RegisterSchema("/public", schema.GetSchema(), &resolver.Resolver{})

	core.Server(opts)
}
`

	templates["git"] = `### Go ###
# Binaries for programs and plugins
*.exe
*.dll
*.so
*.dylib

# Test binary, build with ` + "`" + `go test -c` + "`" + `
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Project-local glide cache, RE: https://github.com/Masterminds/glide/issues/736
.glide/

# Dep Vendors
vendor/*/**

# dotenv
.env

### macOS ###
*.DS_Store
.AppleDouble
.LSOverride

# Icon must end with two \r
Icon

# Thumbnails
._*

# Files that might appear in the root of a volume
.DocumentRevisions-V100
.fseventsd
.Spotlight-V100
.TemporaryItems
.Trashes
.VolumeIcon.icns
.com.apple.timemachine.donotpresent

# Directories potentially created on remote AFP share
.AppleDB
.AppleDesktop
Network Trash Folder
Temporary Items
.apdisk

#vscode
.vscode/*
`

	templates["consts"] = `package lib

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

const (
	// ContextKeyAuth const for Context authorization
	ContextKeyAuth = contextKey("Authorization")
)
`

	templates["resolverHelper"] = `package resolver

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// ConnectToAPI Connect to a url service on a given method with body params and inject the object response in v
func ConnectToAPI(method string, url string, header map[string]string, body map[string]interface{}, v interface{}) (int, error) {
	// 1. Create Body
	var b io.Reader
	if body != nil {
		jsonStr, err := json.Marshal(body)
		if err != nil {
			return 0, err
		}
		b = bytes.NewBuffer(jsonStr)
	}

	// 2. Create the Request
	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return 0, err
	}

	// 3. Put headers on request
	if header != nil {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}

	// 4. Run the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return res.StatusCode, err
	}
	defer res.Body.Close()

	// DEBUGGING: To print response body just uncomment the two lines below
	// ATTENTION, THIS WILL BREAK YOUR CODE, JUST FOR DEBUGGING!
	// bodyBytes, _ := ioutil.ReadAll(res.Body)
	// fmt.Println(string(bodyBytes))

	// 5. Decode the resultBody into interface
	if v != nil && res.StatusCode != 204 {
		err = json.NewDecoder(res.Body).Decode(v)
		if err != nil {
			return res.StatusCode, err
		}
	}

	return res.StatusCode, nil
}

func get(url string, header map[string]string, body map[string]interface{}, responseBody interface{}) (int, error) {
	return ConnectToAPI(http.MethodGet, url, header, body, responseBody)
}
func post(url string, header map[string]string, body map[string]interface{}, responseBody interface{}) (int, error) {
	return ConnectToAPI(http.MethodPost, url, header, body, responseBody)
}
func put(url string, header map[string]string, body map[string]interface{}, responseBody interface{}) (int, error) {
	return ConnectToAPI(http.MethodPut, url, header, body, responseBody)
}
func delete(url string, header map[string]string, body map[string]interface{}, responseBody interface{}) (int, error) {
	return ConnectToAPI(http.MethodDelete, url, header, body, responseBody)
}

`
	templates["modelHelper"] = `package model

import (
	"bytes"
	"errors"
	"reflect"
	"regexp"
	"strings"
	"time"
)

func cleanNomNumbers(phone *string) {
	if phone != nil {
		reg := regexp.MustCompile("[^0-9]+")
		*phone = reg.ReplaceAllString(*phone, "")
	}
}

func mask(s string, pattern string) string {
	re := regexp.MustCompile(` + "`" + `^[0-9]+$` + "`" + `)

	if !re.MatchString(s) {
		panic(errors.New("String should contain only numbers"))
	}

	len := len(s)
	count := 0
	var buff []byte
	pb := bytes.NewBufferString(pattern).Bytes()

	for _, i := range pb {
		if string(i) == "0" {
			if len > count {
				buff = append(buff, s[count])
				count++
			}
		} else {
			buff = append(buff, i)
		}
	}

	return string(buff)
}

func dateFormat(s string) string {
	s = strings.Replace(s, "aaaa", "2006", -1)
	s = strings.Replace(s, "aa", "06", -1)
	s = strings.Replace(s, "Mon", "Jan", -1)
	s = strings.Replace(s, "MM", "01", -1)
	s = strings.Replace(s, "M", "1", -1)
	s = strings.Replace(s, "dd", "02", -1)
	s = strings.Replace(s, "d", "2", -1)
	s = strings.Replace(s, "HH", "15", -1)
	s = strings.Replace(s, "hh", "03", -1)
	s = strings.Replace(s, "h", "3", -1)
	s = strings.Replace(s, "mm", "04", -1)
	s = strings.Replace(s, "m", "4", -1)
	s = strings.Replace(s, "ss", "05", -1)
	s = strings.Replace(s, "s", "5", -1)
	s = strings.Replace(s, "p4", "pm", -1)
	return s
}

func getDateInFormat(date *time.Time, format *string) *string {
	if date != nil {
		if format != nil {
			str := date.Format(dateFormat(*format))
			return &str
		}
		str := date.Format(time.RFC3339)
		return &str
	}
	return nil
}

func in(obj interface{}, s ...interface{}) bool {
	if obj != nil && s != nil && len(s) > 0 {
		objType := reflect.TypeOf(obj).Elem()
		objValue := reflect.ValueOf(obj).Elem()
		for i := range s {
			if s[i] != nil && reflect.TypeOf(s[i]).Elem() == objType && reflect.ValueOf(s[i]).Elem() == objValue {
				return true
			}
		}
	}
	return false
}
`

	templates["json"] = `package scalar

import (
	"encoding/json"
	"errors"
)

// Json represents GraphQL's "Json" scalar type.
type Json map[string]interface{}

func (Json) ImplementsGraphQLType(name string) bool {
	return name == "Json"
}

func (j *Json) UnmarshalGraphQL(input interface{}) error {
	switch input := input.(type) {
	case string:
		return json.Unmarshal([]byte(input), j)
	case map[string]interface{}:
		*j = Json(input)
		return nil
	default:
		return errors.New("wrong type")
	}
}
`
	templates["scalar"] = `package scalar

// Scalars for graphql definition
const Scalars = ` + "`" + `
scalar Time
scalar Json
` + "`" + `
`
}
