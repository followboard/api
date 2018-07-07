package main

import (
	"flag"

	"github.com/followboard/api/server"
)

func main() {
	flag.Parse()
	server.New().Start()
}
