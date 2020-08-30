package undra

import (
	"reflect"

	"github.com/omm-lang/goat"
	"github.com/omm-lang/omm/lang/types"
	"github.com/omm-lang/omm/stdlib/native"
)

func init() {

	//define the operations to access go struct fields in omm

	goat.DefOp("undra-request", "string", "::", func(val1, val2 types.OmmType, instance *types.Instance, stacktrace []string, line uint64, file string, stacksize uint) *types.OmmType {
		asserted := val1.(OmmHTTPRequest)
		return native.GoatProtoIndex(reflect.ValueOf(&asserted), val2.(types.OmmString), stacktrace, line, file)
	})
	goat.DefOp("undra-response", "string", "::", func(val1, val2 types.OmmType, instance *types.Instance, stacktrace []string, line uint64, file string, stacksize uint) *types.OmmType {
		asserted := val1.(types.OmmString)
		return native.GoatProtoIndex(reflect.ValueOf(&asserted), val2.(types.OmmString), stacktrace, line, file)
	})

	/////////////////////////////////////////////////////////
}
