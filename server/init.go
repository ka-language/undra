package undra

import (
	"os"
	"path/filepath"
	"runtime"

	"goat"
	"omm/lang/types"
)

var (
	requestProto  types.OmmProto
	responseProto types.OmmProto
)

func init() {
	//get the request and response protos from undrastd/

	_, ex, _, _ := runtime.Caller(0)
	os.Chdir(filepath.Dir(ex))
	os.Chdir("..")

	lib, _ := goat.LoadLibrary("undrastd", types.CliParams{})

	requestProto = (*lib.Fetch("$undra_request").Value).(types.OmmProto)
	responseProto = (*lib.Fetch("$undra_response").Value).(types.OmmProto)
}
