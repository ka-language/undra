package undra

import (
	"fmt"
	"net/http"
	"os"
	"unsafe"
)

//Goatv is the location of the oat shared library (or dynamic library) "CallFunc" function
var Goatv unsafe.Pointer

//StartServer starts the undra server in the current working directory
func StartServer(host string) {

	http.HandleFunc("/", handle)

	if e := http.ListenAndServe(host, nil); e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}
