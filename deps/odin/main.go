package main

import "github.com/Sam-Izdat/pogo/deps/odin/cli"

// VERSION is the odin version
var VERSION = "1.4.0"

var app = cli.New(VERSION, "a command line DSL for go-lang", cli.ShowUsage)

func main() {
	app.Start()
}
