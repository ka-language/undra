package undra

import (
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/tusklang/tusk/lang/interpreter"
	"github.com/tusklang/tusk/lang/types"
)

func createResponse(res http.ResponseWriter, req *http.Request) *types.TuskType {

	var response = responseProto.New(types.Instance{})

	wrender, _ := response.Get("render", "")
	*wrender = interpreter.TuskGoFunc{
		Function: func(args []*types.TuskType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.TuskType {

			if len(args) == 0 {
				staticsend(res, req)
				var undef types.TuskType = types.TuskUndef{}
				return &undef
			} else if len(args) == 1 && (*args[0]).Type() == "hash" {

				var hash = (*args[0]).(types.TuskHash)
				var template = make(map[string]string)

				for k, v := range hash.Hash {
					var str = (*interpreter.Cast(*v, "string", stacktrace, line, file)).(types.TuskString).ToGoType()
					template[k] = str

					templated, e := templatedoc(path.Join("public", req.URL.Path), template)

					if e != nil {
						interpreter.TuskPanic("File "+req.URL.Path+" does not exist in the public directory", line, file, stacktrace)
					}

					res.Header().Set("Content-Type", "text/html")
					fmt.Fprint(res, templated)
					res.Header().Set("Content-Type", "text/plain")
				}

				var undef types.TuskType = types.TuskUndef{}
				return &undef
			}

			interpreter.TuskPanic("Function undra-response::render requires an argument count of 0 or 1 where the first argument is of type hash", line, file, stacktrace)

			var undef types.TuskType = types.TuskUndef{}
			return &undef
		},
	}

	wsend, _ := response.Get("send", "")
	*wsend = interpreter.TuskGoFunc{
		Function: func(args []*types.TuskType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.TuskType {

			for _, v := range args {
				strval := (*interpreter.Cast(*v, "string", stacktrace, line, file)).(types.TuskString).ToGoType()
				fmt.Fprintf(res, "%s", strval)
			}

			var undef types.TuskType = types.TuskUndef{}
			return &undef
		},
	}

	wsetcookie, _ := response.Get("setcookie", "")
	*wsetcookie = interpreter.TuskGoFunc{
		Function: func(args []*types.TuskType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.TuskType {

			invalidsig := "Function undra-response::setcookie requires the parameter signature: (string, hash)"

			if len(args) != 2 {
				interpreter.TuskPanic(invalidsig, line, file, stacktrace)
			}

			if (*args[0]).Type() != "string" || (*args[1]).Type() != "hash" {
				interpreter.TuskPanic(invalidsig, line, file, stacktrace)
			}

			var gocookie http.Cookie
			var kahash = (*args[1]).(types.TuskHash)
			gocookie.Name = (*args[0]).(types.TuskString).ToGoType()

			var testtype = func(kav *types.TuskType, typeof, fieldname string) bool {

				if kav == nil || (*kav).Type() == "undef" {
					return false
				}

				if (*kav).Type() != typeof {
					interpreter.TuskPanic("Expected type "+typeof+" for field "+fieldname+" in an undra-response", line, file, stacktrace)
				}

				return true
			}

			//set all of the fields
			var oval *types.TuskType

			oval = kahash.At("value")
			if testtype(oval, "string", "value") {
				gocookie.Value = (*oval).(types.TuskString).ToGoType()
			}

			oval = kahash.At("path")
			if testtype(oval, "string", "path") {
				gocookie.Path = (*oval).(types.TuskString).ToGoType()
			}

			oval = kahash.At("domain")
			if testtype(oval, "string", "domain") {
				gocookie.Domain = (*oval).(types.TuskString).ToGoType()
			}

			oval = kahash.At("expires")
			if testtype(oval, "string", "number") {
				gocookie.Expires = time.Now()                                            //set the expires to now
				gocookie.Expires.Add(time.Duration((*oval).(types.TuskNumber).ToGoType())) //and add to it
			}

			oval = kahash.At("maxage")
			if testtype(oval, "string", "number") {
				gocookie.MaxAge = int((*oval).(types.TuskNumber).ToGoType())
			}

			oval = kahash.At("secure")
			if testtype(oval, "string", "bool") {
				gocookie.Secure = (*oval).(types.TuskBool).ToGoType()
			}

			oval = kahash.At("httponly")
			if testtype(oval, "string", "httponly") {
				gocookie.HttpOnly = (*oval).(types.TuskBool).ToGoType()
			}

			oval = kahash.At("samesite")
			if testtype(oval, "string", "number") {
				gocookie.SameSite = http.SameSite((*oval).(types.TuskNumber).ToGoType())
			}
			///////////////////////

			http.SetCookie(res, &gocookie)

			var undef types.TuskType = types.TuskUndef{}
			return &undef
		},
	}

	wclearcookie, _ := response.Get("clearcookie", "")
	*wclearcookie = interpreter.TuskGoFunc{
		Function: func(args []*types.TuskType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.TuskType {

			if len(args) != 1 || (*args[0]).Type() != "string" {
				interpreter.TuskPanic("Function undra-response::clearcookie requires a argument count of 1 with the type of string", line, file, stacktrace)
			}

			var name = (*args[0]).(types.TuskString).ToGoType()
			http.SetCookie(res, &http.Cookie{ //set the expires to 1970, Jan 1 (unix epoch)
				Name:    name,
				Value:   "",
				Expires: time.Unix(0, 0),
			})

			var undef types.TuskType = types.TuskUndef{}
			return &undef
		},
	}

	wredirect, _ := response.Get("redirect", "")
	*wredirect = interpreter.TuskGoFunc{
		Function: func(args []*types.TuskType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.TuskType {
			if len(args) != 1 || (*args[0]).Type() != "string" {
				interpreter.TuskPanic("Function undra-response::redirect requires a argument count of 1 with the type of string", line, file, stacktrace)
			}

			nurl := (*args[0]).(types.TuskString).ToGoType()
			http.Redirect(res, req, nurl, http.StatusSeeOther)

			var undef types.TuskType = types.TuskUndef{}
			return &undef
		},
	}

	werror, _ := response.Get("error", "")
	*werror = interpreter.TuskGoFunc{
		Function: func(args []*types.TuskType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.TuskType {

			if len(args) != 2 || (*args[0]).Type() != "string" || (*args[1]).Type() != "number" {
				interpreter.TuskPanic("Function undra-response::error requires the parameter signature: (string, number)", line, file, stacktrace)
			}

			msg := (*args[0]).(types.TuskString).ToGoType()
			err := int((*args[1]).(types.TuskNumber).ToGoType())

			http.Error(res, msg, err)

			var undef types.TuskType = types.TuskUndef{}
			return &undef
		},
	}

	wheader, _ := response.Get("header", "")
	*wheader = interpreter.TuskGoFunc{
		Function: func(args []*types.TuskType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.TuskType {

			if len(args) != 2 || (*args[0]).Type() != "string" || (*args[1]).Type() != "string" {
				interpreter.TuskPanic("Function undra-response::header requires the parameter signature: (string, string)", line, file, stacktrace)
			}

			//get the name and value as go strings
			name := (*args[0]).(types.TuskString).ToGoType()
			value := (*args[1]).(types.TuskString).ToGoType()
			//////////////////////////////////////

			res.Header().Set(name, value) //set the header

			var undef types.TuskType = types.TuskUndef{}
			return &undef
		},
	}

	var tusktype types.TuskType = response
	return &tusktype
}
