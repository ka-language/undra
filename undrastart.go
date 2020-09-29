package main

import (
	"flag"
	"os"

	"github.com/tusklang/tools"

	"github.com/tusklang/tusk/lang/types"
	undra "github.com/tusklang/undra/server"
)

var host = flag.String("host", "localhost:80", "Set the host:port of the undra instance")

func init() {
	flag.Usage = tools.Usagef("Undra")
}

func main() {
	flag.Parse()

	basedir, _ := os.Executable()
	var params types.CliParams
	params.TuskDirname = basedir

	undra.StartServer(*host, params)
}
