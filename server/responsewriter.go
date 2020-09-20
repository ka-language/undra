package undra

import (
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/tusklang/tusk/lang/interpreter"
	"github.com/tusklang/tusk/lang/types"
	"github.com/tusklang/tusk/native"
)

func createResponse(res http.ResponseWriter, req *http.Request) *types.TuskType {

	var response = responseProto.New(types.Instance{})

	wrender, _ := response.Get("render", "")
	*wrender = native.TuskGoFunc{
		Function: func(args []*types.TuskType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.TuskType {

			if len(args) == 0 {
				staticsend(res, req)
				var undef types.TuskType = types.TuskUndef{}
				return &undef
			} else if len(args) == 1 && (*args[0]).Type() == "hash" {

				var hash = (*args[0]).(types.TuskHash)
				var template = make(map[string]string)

				for k, v := range hash.Hash {
					var strk = k
					var strv = (*interpreter.Cast(*v, "string", stacktrace, line, file)).(types.TuskString).ToGoType()
					template[strk] = strv

					templated, e := templatedoc(path.Join("public", req.URL.Path), template)

					if e != nil {
						native.TuskPanic("File "+req.URL.Path+" does not exist in the public directory", line, file, stacktrace)
					}

					res.Header().Set("Content-Type", "text/html")
					fmt.Fprint(res, templated)
					res.Header().Set("Content-Type", "text/plain")
				}

				var undef types.TuskType = types.TuskUndef{}
				return &undef
			}

			native.TuskPanic("Function undra-response::render requires an argument count of 0 or 1 where the first argument is of type hash", line, file, stacktrace)

			var undef types.TuskType = types.TuskUndef{}
			return &undef
		},
	}

	wsend, _ := response.Get("send", "")
	*wsend = native.TuskGoFunc{
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
	*wsetcookie = native.TuskGoFunc{
		Function: func(args []*types.TuskType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.TuskType {

			invalidsig := "Function undra-response::setcookie requires the parameter signature: (string, hash)"

			if len(args) != 2 {
				native.TuskPanic(invalidsig, line, file, stacktrace)
			}

			if (*args[0]).Type() != "string" || (*args[1]).Type() != "hash" {
				native.TuskPanic(invalidsig, line, file, stacktrace)
			}

			var gocookie http.Cookie
			var kahash = (*args[1]).(types.TuskHash)
			gocookie.Name = (*args[0]).(types.TuskString).ToGoType()

			var testtype = func(kav *types.TuskType, typeof, fieldname string) bool {

				if kav == nil || (*kav).Type() == "undef" {
					return false
				}

				if (*kav).Type() != typeof {
					native.TuskPanic("Expected type "+typeof+" for field "+fieldname+" in an undra-response", line, file, stacktrace)
				}

				return true
			}

			//set all of the fields
			var tval *types.TuskType

			tval = kahash.AtStr("value")
			if testtype(tval, "string", "value") {
				gocookie.Value = (*tval).(types.TuskString).ToGoType()
			}

			tval = kahash.AtStr("path")
			if testtype(tval, "string", "path") {
				gocookie.Path = (*tval).(types.TuskString).ToGoType()
			}

			tval = kahash.AtStr("domain")
			if testtype(tval, "string", "domain") {
				gocookie.Domain = (*tval).(types.TuskString).ToGoType()
			}

			tval = kahash.AtStr("expires")
			if testtype(tval, "string", "number") {
				gocookie.Expires = time.Now()                                              //set the expires to now
				gocookie.Expires.Add(time.Duration((*tval).(types.TuskNumber).ToGoType())) //and add to it
			}

			tval = kahash.AtStr("maxage")
			if testtype(tval, "string", "number") {
				gocookie.MaxAge = int((*tval).(types.TuskNumber).ToGoType())
			}

			tval = kahash.AtStr("secure")
			if testtype(tval, "string", "bool") {
				gocookie.Secure = (*tval).(types.TuskBool).ToGoType()
			}

			tval = kahash.AtStr("httponly")
			if testtype(tval, "string", "httponly") {
				gocookie.HttpOnly = (*tval).(types.TuskBool).ToGoType()
			}

			tval = kahash.AtStr("samesite")
			if testtype(tval, "string", "number") {
				gocookie.SameSite = http.SameSite((*tval).(types.TuskNumber).ToGoType())
			}
			///////////////////////

			http.SetCookie(res, &gocookie)

			var undef types.TuskType = types.TuskUndef{}
			return &undef
		},
	}

	wclearcookie, _ := response.Get("clearcookie", "")
	*wclearcookie = native.TuskGoFunc{
		Function: func(args []*types.TuskType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.TuskType {

			if len(args) != 1 || (*args[0]).Type() != "string" {
				native.TuskPanic("Function undra-response::clearcookie requires a argument count of 1 with the type of string", line, file, stacktrace)
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
	*wredirect = native.TuskGoFunc{
		Function: func(args []*types.TuskType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.TuskType {
			if len(args) != 1 || (*args[0]).Type() != "string" {
				native.TuskPanic("Function undra-response::redirect requires a argument count of 1 with the type of string", line, file, stacktrace)
			}

			nurl := (*args[0]).(types.TuskString).ToGoType()
			http.Redirect(res, req, nurl, http.StatusSeeOther)

			var undef types.TuskType = types.TuskUndef{}
			return &undef
		},
	}

	werror, _ := response.Get("error", "")
	*werror = native.TuskGoFunc{
		Function: func(args []*types.TuskType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.TuskType {

			if len(args) != 2 || (*args[0]).Type() != "string" || (*args[1]).Type() != "number" {
				native.TuskPanic("Function undra-response::error requires the parameter signature: (string, number)", line, file, stacktrace)
			}

			msg := (*args[0]).(types.TuskString).ToGoType()
			err := int((*args[1]).(types.TuskNumber).ToGoType())

			http.Error(res, msg, err)

			var undef types.TuskType = types.TuskUndef{}
			return &undef
		},
	}

	wheader, _ := response.Get("header", "")
	*wheader = native.TuskGoFunc{
		Function: func(args []*types.TuskType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.TuskType {

			if len(args) != 2 || (*args[0]).Type() != "string" || (*args[1]).Type() != "string" {
				native.TuskPanic("Function undra-response::header requires the parameter signature: (string, string)", line, file, stacktrace)
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
