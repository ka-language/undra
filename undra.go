package main

import (
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
		}
	}

	undra.StartServer(host)
}
