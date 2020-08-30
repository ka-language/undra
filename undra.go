package main

import (
	"fmt"
	"os"

	undra "github.com/omm-lang/undra/server"
)

func main() {

	dirname, _ := os.Getwd() //get the working directory
	os.Chdir(dirname)        //and change to it

	var host = "localhost:80"

	if len(os.Args) > 1 {
		host = os.Args[1]
	}

	if len(os.Args) > 2 {
		switch os.Args[2] {
		case "-v":
			fmt.Printf("%d.%d.%d", undra.UNDRA_MAJOR, undra.UNDRA_MINOR, undra.UNDRA_BUG)
			os.Exit(0)
		}
	}

	undra.StartServer(host)
}
