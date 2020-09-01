package undra

import (
	"os"
	"path/filepath"

	"github.com/omm-lang/goat"
	"github.com/omm-lang/omm/lang/types"
)

var (
	requestProto  types.OmmProto
	responseProto types.OmmProto
)

func init() {
	//get the request and response protos from undrastd/

	ex, _ := os.Executable()
	os.Chdir(filepath.Dir(ex))

	lib, _ := goat.LoadLibrary("undrastd/undra.omm", types.CliParams{})

	requestProto = (*lib.Fetch("$undra_request").Value).(types.OmmProto)
	responseProto = (*lib.Fetch("$undra_response").Value).(types.OmmProto)
}
