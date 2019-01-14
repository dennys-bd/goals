package errors

import (
	"fmt"
	"os"
)

// CheckEx function verify if error is nil and calls Ex() on error
func CheckEx(err error) {
	if err != nil {
		Ex(err.Error())
	}
}

// Ex function print error and exits
func Ex(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

// Check function verify if error is nil and panics on error
func Check(err error) {
	if err != nil {
		panic(err.Error())
	}
}
