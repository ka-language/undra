package main

import (
	"flag"
	"os"

	suite "github.com/omm-lang/omm-suite"
	"github.com/omm-lang/omm/lang/types"
	undra "github.com/omm-lang/undra/server"
)

var cwd = flag.String("cwd", "", "set the current working directory (automatically placed by the shell/pwsh script)")
var host = flag.String("host", "localhost:80", "set ")

func init() {
	flag.Usage = suite.Usagef("Undra")
}

func main() {
	flag.Parse()

	basedir, _ := os.Executable()
	var params types.CliParams
	params.OmmDirname = basedir

	os.Chdir(*cwd) //change to the cwd

	undra.StartServer(*host, params)
}
