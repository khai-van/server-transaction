package main

import (
	"fmt"
	"os"

	"app/api"
)

// file config default
const defaultConfigFilename = "server.yml"

var serverConfig api.Configuration

func main() {
	// read config
	err := serverConfig.BindFile(defaultConfigFilename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// create server
	srv := api.NewServer(serverConfig)
	// start server
	if err := srv.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
