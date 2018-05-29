package cmd

func initConfig() {
	Templates["config"] = `package lib

import "os"

// GetPort returns the port to start the server
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
`
}
