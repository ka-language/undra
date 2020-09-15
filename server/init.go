package undra

import (
	"goat"
	"os"
	"path/filepath"
	"runtime"

	"ka/lang/types"
)

var (
	requestProto  types.KaProto
	responseProto types.KaProto
)

func init() {
	//get the request and response protos from undrastd/

	_, ex, _, _ := runtime.Caller(0)
	os.Chdir(filepath.Dir(ex))
	os.Chdir("..")

	lib, _ := goat.LoadLibrary("undrastd", types.CliParams{})

	requestProto = (*lib.Fetch("$undra_request").Value).(types.KaProto)
	responseProto = (*lib.Fetch("$undra_response").Value).(types.KaProto)
}
