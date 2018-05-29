package cmd

func initHelpers() {
	Templates["resolverHelper"] = `package resolver

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// ConnectToAPI Connect to a url service on a given method with body params and inject the object response in v
func ConnectToAPI(method string, url string, body *map[string]interface{}, header *map[string]string, v interface{}) (int, error) {

	client := &http.Client{}

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
		for k, v := range *header {
			req.Header.Add(k, v)
		}
	}

	// 4. Run the request
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

func get(url string, body *map[string]interface{}, header *map[string]string, responseBody interface{}) (int, error) {
	return ConnectToAPI("GET", url, body, header, responseBody)
}
func post(url string, body *map[string]interface{}, header *map[string]string, responseBody interface{}) (int, error) {
	return ConnectToAPI("POST", url, body, header, responseBody)
}
func put(url string, body *map[string]interface{}, header *map[string]string, responseBody interface{}) (int, error) {
	return ConnectToAPI("PUT", url, body, header, responseBody)
}
func delete(url string, body *map[string]interface{}, header *map[string]string, responseBody interface{}) (int, error) {
	return ConnectToAPI("DELETE", url, body, header, responseBody)
}

func getHeaders(ctx *context.Context) *map[string]string {
	m := make(map[string]string)
	m["Content-Type"] = "application/json"
	m["Accept"] = "application/json, text/plain, */*"
	// TODO: Change the client_id
	// m["client_id"] = "4"
	if ctx != nil {
		// TODO: APIGATEWAY MODE: Create here the types comming from context that will proceed as
		// header in the request, applied in every single request e.g:
		// con := *ctx
		// m["access-token"] = con.Value(ContextKeyAuth).(string)

	}
	return &m
}

func getFeedParams(Cursor *int32) (int32, int32) {
	var cursor int32
	var limit int32
	limit = 10

	if Cursor != nil {
		cursor = *Cursor
		cursor = (cursor - 1) * 10
	}

	return cursor, limit
}
`
	Templates["modelHelper"] = `package model

import (
	"bytes"
	"errors"
	"reflect"
	"regexp"
	"strings"
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
}
