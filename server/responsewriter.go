package undra

import (
	"fmt"
	"net/http"
	"path"
	"time"

	"ka/lang/interpreter"
	"ka/lang/types"
)

func createResponse(res http.ResponseWriter, req *http.Request) *types.KaType {

	var response = responseProto.New(types.Instance{})

	wrender, _ := response.Get("render", "")
	*wrender = interpreter.KaGoFunc{
		Function: func(args []*types.KaType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.KaType {

			if len(args) == 0 {
				staticsend(res, req)
				var undef types.KaType = types.KaUndef{}
				return &undef
			} else if len(args) == 1 && (*args[0]).Type() == "hash" {

				var hash = (*args[0]).(types.KaHash)
				var template = make(map[string]string)

				for k, v := range hash.Hash {
					var str = (*interpreter.Cast(*v, "string", stacktrace, line, file)).(types.KaString).ToGoType()
					template[k] = str

					templated, e := templatedoc(path.Join("public", req.URL.Path), template)

					if e != nil {
						interpreter.KaPanic("File "+req.URL.Path+" does not exist in the public directory", line, file, stacktrace)
					}

					res.Header().Set("Content-Type", "text/html")
					fmt.Fprint(res, templated)
					res.Header().Set("Content-Type", "text/plain")
				}

				var undef types.KaType = types.KaUndef{}
				return &undef
			}

			interpreter.KaPanic("Function undra-response::render requires an argument count of 0 or 1 where the first argument is of type hash", line, file, stacktrace)

			var undef types.KaType = types.KaUndef{}
			return &undef
		},
	}

	wsend, _ := response.Get("send", "")
	*wsend = interpreter.KaGoFunc{
		Function: func(args []*types.KaType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.KaType {

			for _, v := range args {
				strval := (*interpreter.Cast(*v, "string", stacktrace, line, file)).(types.KaString).ToGoType()
				fmt.Fprintf(res, "%s", strval)
			}

			var undef types.KaType = types.KaUndef{}
			return &undef
		},
	}

	wsetcookie, _ := response.Get("setcookie", "")
	*wsetcookie = interpreter.KaGoFunc{
		Function: func(args []*types.KaType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.KaType {

			invalidsig := "Function undra-response::setcookie requires the parameter signature: (string, hash)"

			if len(args) != 2 {
				interpreter.KaPanic(invalidsig, line, file, stacktrace)
			}

			if (*args[0]).Type() != "string" || (*args[1]).Type() != "hash" {
				interpreter.KaPanic(invalidsig, line, file, stacktrace)
			}

			var gocookie http.Cookie
			var kahash = (*args[1]).(types.KaHash)
			gocookie.Name = (*args[0]).(types.KaString).ToGoType()

			var testtype = func(kav *types.KaType, typeof, fieldname string) bool {

				if kav == nil || (*kav).Type() == "undef" {
					return false
				}

				if (*kav).Type() != typeof {
					interpreter.KaPanic("Expected type "+typeof+" for field "+fieldname+" in an undra-response", line, file, stacktrace)
				}

				return true
			}

			//set all of the fields
			var oval *types.KaType

			oval = kahash.At("value")
			if testtype(oval, "string", "value") {
				gocookie.Value = (*oval).(types.KaString).ToGoType()
			}

			oval = kahash.At("path")
			if testtype(oval, "string", "path") {
				gocookie.Path = (*oval).(types.KaString).ToGoType()
			}

			oval = kahash.At("domain")
			if testtype(oval, "string", "domain") {
				gocookie.Domain = (*oval).(types.KaString).ToGoType()
			}

			oval = kahash.At("expires")
			if testtype(oval, "string", "number") {
				gocookie.Expires = time.Now()                                            //set the expires to now
				gocookie.Expires.Add(time.Duration((*oval).(types.KaNumber).ToGoType())) //and add to it
			}

			oval = kahash.At("maxage")
			if testtype(oval, "string", "number") {
				gocookie.MaxAge = int((*oval).(types.KaNumber).ToGoType())
			}

			oval = kahash.At("secure")
			if testtype(oval, "string", "bool") {
				gocookie.Secure = (*oval).(types.KaBool).ToGoType()
			}

			oval = kahash.At("httponly")
			if testtype(oval, "string", "httponly") {
				gocookie.HttpOnly = (*oval).(types.KaBool).ToGoType()
			}

			oval = kahash.At("samesite")
			if testtype(oval, "string", "number") {
				gocookie.SameSite = http.SameSite((*oval).(types.KaNumber).ToGoType())
			}
			///////////////////////

			http.SetCookie(res, &gocookie)

			var undef types.KaType = types.KaUndef{}
			return &undef
		},
	}

	wclearcookie, _ := response.Get("clearcookie", "")
	*wclearcookie = interpreter.KaGoFunc{
		Function: func(args []*types.KaType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.KaType {

			if len(args) != 1 || (*args[0]).Type() != "string" {
				interpreter.KaPanic("Function undra-response::clearcookie requires a argument count of 1 with the type of string", line, file, stacktrace)
			}

			var name = (*args[0]).(types.KaString).ToGoType()
			http.SetCookie(res, &http.Cookie{ //set the expires to 1970, Jan 1 (unix epoch)
				Name:    name,
				Value:   "",
				Expires: time.Unix(0, 0),
			})

			var undef types.KaType = types.KaUndef{}
			return &undef
		},
	}

	wredirect, _ := response.Get("redirect", "")
	*wredirect = interpreter.KaGoFunc{
		Function: func(args []*types.KaType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.KaType {
			if len(args) != 1 || (*args[0]).Type() != "string" {
				interpreter.KaPanic("Function undra-response::redirect requires a argument count of 1 with the type of string", line, file, stacktrace)
			}

			nurl := (*args[0]).(types.KaString).ToGoType()
			http.Redirect(res, req, nurl, http.StatusSeeOther)

			var undef types.KaType = types.KaUndef{}
			return &undef
		},
	}

	werror, _ := response.Get("error", "")
	*werror = interpreter.KaGoFunc{
		Function: func(args []*types.KaType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.KaType {

			if len(args) != 2 || (*args[0]).Type() != "string" || (*args[1]).Type() != "number" {
				interpreter.KaPanic("Function undra-response::error requires the parameter signature: (string, number)", line, file, stacktrace)
			}

			msg := (*args[0]).(types.KaString).ToGoType()
			err := int((*args[1]).(types.KaNumber).ToGoType())

			http.Error(res, msg, err)

			var undef types.KaType = types.KaUndef{}
			return &undef
		},
	}

	wheader, _ := response.Get("header", "")
	*wheader = interpreter.KaGoFunc{
		Function: func(args []*types.KaType, stacktrace []string, line uint64, file string, instance *types.Instance) *types.KaType {

			if len(args) != 2 || (*args[0]).Type() != "string" || (*args[1]).Type() != "string" {
				interpreter.KaPanic("Function undra-response::header requires the parameter signature: (string, string)", line, file, stacktrace)
			}

			//get the name and value as go strings
			name := (*args[0]).(types.KaString).ToGoType()
			value := (*args[1]).(types.KaString).ToGoType()
			//////////////////////////////////////

			res.Header().Set(name, value) //set the header

			var undef types.KaType = types.KaUndef{}
			return &undef
		},
	}

	var katype types.KaType = response
	return &katype
}
