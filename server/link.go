package undra

import (
	"fmt"
	"net/http"
	"os"

	"ka/lang/types"
)

var params types.CliParams

//StartServer starts the undra server in the current working directory
func StartServer(host string, _params types.CliParams) {

	params = _params

	http.HandleFunc("/", handle)

	if e := http.ListenAndServe(host, nil); e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}
