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

	wrender, _ := response.Get("render", "", "global")
	*wrender = native.TuskGoFunc{
		Function: func(args []*types.TuskType, stacktrace []string, line uint64, file string, instance *types.Instance) (*types.TuskType, *types.TuskError) {

			if len(args) == 0 {
				staticsend(res, req)
				var undef types.TuskType = types.TuskUndef{}
				return &undef, nil
			} else if len(args) == 1 && (*args[0]).Type() == "hash" {

				var hash = (*args[0]).(types.TuskHash)
				var template = make(map[string]string)

				hash.Range(func(k, v *types.TuskType) (types.Returner, *types.TuskError) {
					var strk = (*k).Format()
					casted, e := interpreter.Cast(*v, "string", stacktrace, line, file)

					if e != nil {
						e.Print()
					}

					var strv = (*casted).(types.TuskString).ToGoType()
					template[strk] = strv

					templated, e2 := templatedoc(path.Join("public", req.URL.Path), template)

					if e2 != nil {
						native.TuskPanic("File "+req.URL.Path+" does not exist in the public directory", line, file, stacktrace, native.ErrCodes["FILENOTFOUND"]).Print()
					}

					res.Header().Set("Content-Type", "text/html")
					fmt.Fprint(res, templated)
					res.Header().Set("Content-Type", "text/plain")

					return types.Returner{}, nil
				})

				var undef types.TuskType = types.TuskUndef{}
				return &undef, nil
			}

			return nil, native.TuskPanic("Invalid signature given to render", line, file, stacktrace, native.ErrCodes["SIGNOMATCH"])
		},
	}

	wsend, _ := response.Get("send", "", "global")
	*wsend = native.TuskGoFunc{
		Function: func(args []*types.TuskType, stacktrace []string, line uint64, file string, instance *types.Instance) (*types.TuskType, *types.TuskError) {

			for _, v := range args {
				casted, e := interpreter.Cast(*v, "string", stacktrace, line, file)
				if e != nil {
					e.Print()
				}
				strval := (*casted).(types.TuskString).ToGoType()
				fmt.Fprintf(res, "%s", strval)
			}

			var undef types.TuskType = types.TuskUndef{}
			return &undef, nil
		},
	}

	wsetcookie, _ := response.Get("setcookie", "", "global")
	*wsetcookie = native.TuskGoFunc{
		Function: func(args []*types.TuskType, stacktrace []string, line uint64, file string, instance *types.Instance) (*types.TuskType, *types.TuskError) {

			invalidsig := "Invalid signature given to setcookie"
			invalidsigErrCode := native.ErrCodes["SIGNOMATCH"]

			if len(args) != 2 {
				native.TuskPanic(invalidsig, line, file, stacktrace, invalidsigErrCode).Print()
			}

			if (*args[0]).Type() != "string" || (*args[1]).Type() != "hash" {
				native.TuskPanic(invalidsig, line, file, stacktrace, invalidsigErrCode)
			}

			var gocookie http.Cookie
			var kahash = (*args[1]).(types.TuskHash)
			gocookie.Name = (*args[0]).(types.TuskString).ToGoType()

			var testtype = func(kav *types.TuskType, typeof, fieldname string) (bool, *types.TuskError) {

				if kav == nil || (*kav).Type() == "undef" {
					return false, nil
				}

				if (*kav).Type() != typeof {
					return false, native.TuskPanic("Expected type "+typeof+" for field "+fieldname+" in an undra_response", line, file, stacktrace, native.ErrCodes["INVALIDARG"])
				}

				return true, nil
			}

			//set all of the fields
			var tval *types.TuskType

			tval = kahash.AtStr("value")
			if b, e := testtype(tval, "string", "value"); b {
				if e != nil {
					return nil, e
				}
				gocookie.Value = (*tval).(types.TuskString).ToGoType()
			}

			tval = kahash.AtStr("path")
			if b, e := testtype(tval, "string", "path"); b {
				if e != nil {
					return nil, e
				}
				gocookie.Path = (*tval).(types.TuskString).ToGoType()
			}

			tval = kahash.AtStr("domain")
			if b, e := testtype(tval, "string", "value"); b {
				if e != nil {
					return nil, e
				}
				gocookie.Domain = (*tval).(types.TuskString).ToGoType()
			}

			tval = kahash.AtStr("expires")
			if b, e := testtype(tval, "string", "expires"); b {
				if e != nil {
					return nil, e
				}
				gocookie.Expires = time.Now()                                              //set the expires to now
				gocookie.Expires.Add(time.Duration((*tval).(types.TuskNumber).ToGoType())) //and add to it
			}

			tval = kahash.AtStr("maxage")
			if b, e := testtype(tval, "string", "maxage"); b {
				if e != nil {
					return nil, e
				}
				gocookie.MaxAge = int((*tval).(types.TuskNumber).ToGoType())
			}

			tval = kahash.AtStr("secure")
			if b, e := testtype(tval, "string", "secure"); b {
				if e != nil {
					return nil, e
				}
				gocookie.Secure = (*tval).(types.TuskBool).ToGoType()
			}

			tval = kahash.AtStr("httponly")
			if b, e := testtype(tval, "string", "httponly"); b {
				if e != nil {
					return nil, e
				}
				gocookie.HttpOnly = (*tval).(types.TuskBool).ToGoType()
			}

			tval = kahash.AtStr("samesite")
			if b, e := testtype(tval, "string", "samesite"); b {
				if e != nil {
					return nil, e
				}
				gocookie.SameSite = http.SameSite((*tval).(types.TuskNumber).ToGoType())
			}
			///////////////////////

			http.SetCookie(res, &gocookie)

			var undef types.TuskType = types.TuskUndef{}
			return &undef, nil
		},
	}

	wclearcookie, _ := response.Get("clearcookie", "", "global")
	*wclearcookie = native.TuskGoFunc{
		Function: func(args []*types.TuskType, stacktrace []string, line uint64, file string, instance *types.Instance) (*types.TuskType, *types.TuskError) {

			if len(args) != 1 || (*args[0]).Type() != "string" {
				return nil, native.TuskPanic("Invalid signature given to clearcookie", line, file, stacktrace, native.ErrCodes["SIGNOMATCH"])
			}

			var name = (*args[0]).(types.TuskString).ToGoType()
			http.SetCookie(res, &http.Cookie{ //set the expires to 1970, Jan 1 (unix epoch)
				Name:    name,
				Value:   "",
				Expires: time.Unix(0, 0),
			})

			var undef types.TuskType = types.TuskUndef{}
			return &undef, nil
		},
	}

	wredirect, _ := response.Get("redirect", "", "global")
	*wredirect = native.TuskGoFunc{
		Function: func(args []*types.TuskType, stacktrace []string, line uint64, file string, instance *types.Instance) (*types.TuskType, *types.TuskError) {
			if len(args) != 1 || (*args[0]).Type() != "string" {
				return nil, native.TuskPanic("Invalid signature given to redirect", line, file, stacktrace, native.ErrCodes["SIGNOMATCH"])
			}

			nurl := (*args[0]).(types.TuskString).ToGoType()
			http.Redirect(res, req, nurl, http.StatusSeeOther)

			var undef types.TuskType = types.TuskUndef{}
			return &undef, nil
		},
	}

	werror, _ := response.Get("error", "", "global")
	*werror = native.TuskGoFunc{
		Function: func(args []*types.TuskType, stacktrace []string, line uint64, file string, instance *types.Instance) (*types.TuskType, *types.TuskError) {

			if len(args) != 2 || (*args[0]).Type() != "string" || (*args[1]).Type() != "number" {
				return nil, native.TuskPanic("Invalid signature given to error", line, file, stacktrace, native.ErrCodes["SIGNOMATCH"])
			}

			msg := (*args[0]).(types.TuskString).ToGoType()
			err := int((*args[1]).(types.TuskNumber).ToGoType())

			http.Error(res, msg, err)

			var undef types.TuskType = types.TuskUndef{}
			return &undef, nil
		},
	}

	wheader, _ := response.Get("header", "", "global")
	*wheader = native.TuskGoFunc{
		Function: func(args []*types.TuskType, stacktrace []string, line uint64, file string, instance *types.Instance) (*types.TuskType, *types.TuskError) {

			if len(args) != 2 || (*args[0]).Type() != "string" || (*args[1]).Type() != "string" {
				return nil, native.TuskPanic("Invalid signature given to header", line, file, stacktrace, native.ErrCodes["SIGNOMATCH"])
			}

			//get the name and value as go strings
			name := (*args[0]).(types.TuskString).ToGoType()
			value := (*args[1]).(types.TuskString).ToGoType()
			//////////////////////////////////////

			res.Header().Set(name, value) //set the header

			var undef types.TuskType = types.TuskUndef{}
			return &undef, nil
		},
	}

	var tusktype types.TuskType = response
	return &tusktype
}
