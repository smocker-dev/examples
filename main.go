package main

import (
	"os"

	"example/server"
	"example/types"

	"github.com/namsral/flag"
)

var config types.Config

func main() {
	flag.CommandLine = flag.NewFlagSetWithEnvPrefix(os.Args[0], "RES", flag.ExitOnError)
	flag.IntVar(&config.Port, "port", 8080, "Port that server will listen on")
	flag.StringVar(&config.HotelsURL, "hotels-url", "http://localhost:8081", "URL of the hotels uservice")
	flag.StringVar(&config.UsersURL, "users-url", "http://localhost:8082", "URL of the users uservice")
	flag.StringVar(&config.DSN, "dsn", "postgres://postgres:@localhost:8083/test?sslmode=disable", "databse DSN")
	flag.Parse()

	server.Serve(config)
}
