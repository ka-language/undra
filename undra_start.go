package main

import (
	"flag"
	"os"

	"github.com/omm-lang/omm/lang/types"
	undra "github.com/omm-lang/undra/server"
)

var cwd = flag.String("cwd", "", "set cwd")
var host = flag.String("host", "localhost:80", "set cwd")

func main() {

	basedir, _ := os.Executable()
	var params types.CliParams
	params.OmmDirname = basedir

	os.Chdir(*cwd) //change to the cwd

	undra.StartServer(*host, params)
}
