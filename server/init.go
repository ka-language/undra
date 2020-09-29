package undra

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/tusklang/goat"

	"github.com/tusklang/tusk/lang/types"
)

var (
	requestProto  types.TuskProto
	responseProto types.TuskProto
)

func init() {
	//get the request and response protos from undrastd/

	_, ex, _, _ := runtime.Caller(0)
	os.Chdir(filepath.Dir(ex))
	os.Chdir("..")

	lib, _ := goat.LoadLibrary("undrastd", types.CliParams{})

	requestProto = (*lib.Fetch("$undra_request").Value).(types.TuskProto)
	responseProto = (*lib.Fetch("$undra_response").Value).(types.TuskProto)
}
