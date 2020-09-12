package main

import (
	"flag"
	"os"

	suite "omm-suite"
	"omm/lang/types"
	undra "undra/server"
)

var host = flag.String("host", "localhost:80", "Set the host:port of the undra instance")

func init() {
	flag.Usage = suite.Usagef("Undra")
}

func main() {
	flag.Parse()

	basedir, _ := os.Executable()
	var params types.CliParams
	params.OmmDirname = basedir

	undra.StartServer(*host, params)
}
