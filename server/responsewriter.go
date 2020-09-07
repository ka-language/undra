package undra

import (
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/omm-lang/omm/lang/interpreter"
	"github.com/omm-lang/omm/lang/types"
	"github.com/omm-lang/omm/native"
)

func createResponse(res http.ResponseWriter, req *http.Request) *types.OmmType {

	var response = responseProto.New(types.Instance{})

	wrender, _ := response.Get("render", "")
	*wrender = native.OmmGoFunc{
		Function: func(args []*types.OmmType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.OmmType {

			if len(args) == 0 {
				staticsend(res, req)
				var undef types.OmmType = types.OmmUndef{}
				return &undef
			} else if len(args) == 1 && (*args[0]).Type() == "hash" {

				var hash = (*args[0]).(types.OmmHash)
				var template = make(map[string]string)

				for k, v := range hash.Hash {
					var str = (*interpreter.Cast(*v, "string", stacktrace, line, file)).(types.OmmString).ToGoType()
					template[k] = str

					templated, e := templatedoc(path.Join("public", req.URL.Path), template)

					if e != nil {
						native.OmmPanic("File "+req.URL.Path+" does not exist in the public directory", line, file, stacktrace)
					}

					res.Header().Set("Content-Type", "text/html")
					fmt.Fprint(res, templated)
					res.Header().Set("Content-Type", "text/plain")
				}

				var undef types.OmmType = types.OmmUndef{}
				return &undef
			}

			native.OmmPanic("Function undra-response::render requires an argument count of 0 or 1 where the first argument is of type hash", line, file, stacktrace)

			var undef types.OmmType = types.OmmUndef{}
			return &undef
		},
	}

	wsend, _ := response.Get("send", "")
	*wsend = native.OmmGoFunc{
		Function: func(args []*types.OmmType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.OmmType {

			for _, v := range args {
				strval := (*interpreter.Cast(*v, "string", stacktrace, line, file)).(types.OmmString).ToGoType()
				fmt.Fprintf(res, "%s", strval)
			}

			var undef types.OmmType = types.OmmUndef{}
			return &undef
		},
	}

	wsetcookie, _ := response.Get("setcookie", "")
	*wsetcookie = native.OmmGoFunc{
		Function: func(args []*types.OmmType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.OmmType {

			invalidsig := "Function undra-response::setcookie requires the parameter signature: (string, hash)"

			if len(args) != 2 {
				native.OmmPanic(invalidsig, line, file, stacktrace)
			}

			if (*args[0]).Type() != "string" || (*args[1]).Type() != "hash" {
				native.OmmPanic(invalidsig, line, file, stacktrace)
			}

			var gocookie http.Cookie
			var ommhash = (*args[1]).(types.OmmHash)
			gocookie.Name = (*args[0]).(types.OmmString).ToGoType()

			var testtype = func(ommv *types.OmmType, typeof, fieldname string) bool {

				if ommv == nil || (*ommv).Type() == "undef" {
					return false
				}

				if (*ommv).Type() != typeof {
					native.OmmPanic("Expected type "+typeof+" for field "+fieldname+" in an undra-response", line, file, stacktrace)
				}

				return true
			}

			//set all of the fields
			var oval *types.OmmType

			oval = ommhash.At("value")
			if testtype(oval, "string", "value") {
				gocookie.Value = (*oval).(types.OmmString).ToGoType()
			}

			oval = ommhash.At("path")
			if testtype(oval, "string", "path") {
				gocookie.Path = (*oval).(types.OmmString).ToGoType()
			}

			oval = ommhash.At("domain")
			if testtype(oval, "string", "domain") {
				gocookie.Domain = (*oval).(types.OmmString).ToGoType()
			}

			oval = ommhash.At("expires")
			if testtype(oval, "string", "number") {
				gocookie.Expires = time.Now()                                             //set the expires to now
				gocookie.Expires.Add(time.Duration((*oval).(types.OmmNumber).ToGoType())) //and add to it
			}

			oval = ommhash.At("maxage")
			if testtype(oval, "string", "number") {
				gocookie.MaxAge = int((*oval).(types.OmmNumber).ToGoType())
			}

			oval = ommhash.At("secure")
			if testtype(oval, "string", "bool") {
				gocookie.Secure = (*oval).(types.OmmBool).ToGoType()
			}

			oval = ommhash.At("httponly")
			if testtype(oval, "string", "httponly") {
				gocookie.HttpOnly = (*oval).(types.OmmBool).ToGoType()
			}

			oval = ommhash.At("samesite")
			if testtype(oval, "string", "number") {
				gocookie.SameSite = http.SameSite((*oval).(types.OmmNumber).ToGoType())
			}
			///////////////////////

			http.SetCookie(res, &gocookie)

			var undef types.OmmType = types.OmmUndef{}
			return &undef
		},
	}

	wclearcookie, _ := response.Get("clearcookie", "")
	*wclearcookie = native.OmmGoFunc{
		Function: func(args []*types.OmmType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.OmmType {

			if len(args) != 1 || (*args[0]).Type() != "string" {
				native.OmmPanic("Function undra-response::clearcookie requires a argument count of 1 with the type of string", line, file, stacktrace)
			}

			var name = (*args[0]).(types.OmmString).ToGoType()
			http.SetCookie(res, &http.Cookie{ //set the expires to 1970, Jan 1 (unix epoch)
				Name:    name,
				Value:   "",
				Expires: time.Unix(0, 0),
			})

			var undef types.OmmType = types.OmmUndef{}
			return &undef
		},
	}

	wredirect, _ := response.Get("redirect", "")
	*wredirect = native.OmmGoFunc{
		Function: func(args []*types.OmmType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.OmmType {
			if len(args) != 1 || (*args[0]).Type() != "string" {
				native.OmmPanic("Function undra-response::redirect requires a argument count of 1 with the type of string", line, file, stacktrace)
			}

			nurl := (*args[0]).(types.OmmString).ToGoType()
			http.Redirect(res, req, nurl, http.StatusSeeOther)

			var undef types.OmmType = types.OmmUndef{}
			return &undef
		},
	}

	werror, _ := response.Get("error", "")
	*werror = native.OmmGoFunc{
		Function: func(args []*types.OmmType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.OmmType {

			if len(args) != 2 || (*args[0]).Type() != "string" || (*args[1]).Type() != "number" {
				native.OmmPanic("Function undra-response::error requires the parameter signature: (string, number)", line, file, stacktrace)
			}

			msg := (*args[0]).(types.OmmString).ToGoType()
			err := int((*args[1]).(types.OmmNumber).ToGoType())

			http.Error(res, msg, err)

			var undef types.OmmType = types.OmmUndef{}
			return &undef
		},
	}

	wheader, _ := response.Get("header", "")
	*wheader = native.OmmGoFunc{
		Function: func(args []*types.OmmType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.OmmType {

			if len(args) != 2 || (*args[0]).Type() != "string" || (*args[1]).Type() != "string" {
				native.OmmPanic("Function undra-response::header requires the parameter signature: (string, string)", line, file, stacktrace)
			}

			//get the name and value as go strings
			name := (*args[0]).(types.OmmString).ToGoType()
			value := (*args[1]).(types.OmmString).ToGoType()
			//////////////////////////////////////

			res.Header().Set(name, value) //set the header

			var undef types.OmmType = types.OmmUndef{}
			return &undef
		},
	}

	var ommtype types.OmmType = response
	return &ommtype
}
