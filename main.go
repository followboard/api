package main

import (
	"flag"

	"github.com/followboard/api/config"
	"github.com/followboard/api/server"
)

const root = "./"

func main() {
	flag.Parse()
	c := config.New(root)
	server.New(c).Start()
}
