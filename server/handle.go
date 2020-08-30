package undra

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/omm-lang/goat"
	"github.com/omm-lang/omm/lang/types"
)

func staticsend(res http.ResponseWriter, req *http.Request) {
	//otherwise, just render the file with no handling (static)
	htmfile := "./public" + req.URL.Path

	if _, f := os.Stat(htmfile); !os.IsNotExist(f) { //only if it exists
		http.ServeFile(res, req, htmfile)
	} else {
		//otherwise serve the 404 not exists file
		_404file := "./public" + "notfound.html"

		if _, f := os.Stat(_404file); !os.IsNotExist(f) { //only if it exists
			http.ServeFile(res, req, _404file)
		} else {
			//otherwise just write 404 path not found
			http.Error(res, "404 path not found", http.StatusNotFound)
		}
	}
}

func handle(res http.ResponseWriter, req *http.Request) {

	//remove the extension, and replace it with .oat
	oatname := strings.TrimSuffix(req.URL.Path, filepath.Ext(req.URL.Path)) + ".oat"

	//prepend the server path
	oatf := path.Join("server", oatname)

	if _, f := os.Stat(oatf); !os.IsNotExist(f) {

		//load the oat file using goat
		lib, e := goat.LoadLibrary(oatf)

		if e != nil {
			fmt.Println(e)
			os.Exit(1)
		}

		var instance = goat.NewInstance(lib)

		var ommreq OmmHTTPRequest
		ommreq.FromGoType(*req)

		var ommres OmmHTTPResponseWriter
		ommres.FromGoType(res, req)

		var reqommtype types.OmmType = ommreq
		var resommtype types.OmmType = ommres
		goat.CallOatFunc(instance, "handle", &reqommtype, &resommtype)

	} else {
		staticsend(res, req)
	}
}
