package undra

import (
	"fmt"
	"net/http"

	"github.com/omm-lang/omm/lang/interpreter"
	"github.com/omm-lang/omm/lang/types"
	"github.com/omm-lang/omm/stdlib/native"
)

//OmmHTTPResponseWriter represents an http response writer in Omm
type OmmHTTPResponseWriter struct {
	Render,
	Send,
	JSON,
	Cookie,
	ClearCookie,
	Redirect,
	Header,
	Status native.OmmGoFunc
}

func (r *OmmHTTPResponseWriter) FromGoType(res http.ResponseWriter, req *http.Request) {
	r.Render = native.OmmGoFunc{
		Function: func(args []*types.OmmType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.OmmType {

			if len(args) == 0 {
				staticsend(res, req)
				var undef types.OmmType = types.OmmUndef{}
				return &undef
			} else if len(args) == 1 && (*args[0]).Type() == "hash" {
				var undef types.OmmType = types.OmmUndef{}
				return &undef
			}

			native.OmmPanic("Function undra-response::Render requires an argument count of 0 or 1 where the first argument is of type hash", line, file, stacktrace)

			var undef types.OmmType = types.OmmUndef{}
			return &undef
		},
	}
	r.Send = native.OmmGoFunc{
		Function: func(args []*types.OmmType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.OmmType {

			for _, v := range args {
				strval := (*interpreter.Cast(*v, "string", stacktrace, line, file)).(types.OmmString).ToGoType()
				fmt.Fprintf(res, "%s", strval)
			}

			var undef types.OmmType = types.OmmUndef{}
			return &undef
		},
	}
}

func (r OmmHTTPResponseWriter) Format() string {
	return "{ undra-response }"
}

func (r OmmHTTPResponseWriter) Type() string {
	return "undra-response"
}

func (r OmmHTTPResponseWriter) TypeOf() string {
	return r.Type()
}

func (_ OmmHTTPResponseWriter) Deallocate() {}
