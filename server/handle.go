package undra

import (
	"fmt"
	"io/ioutil"
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
	htmfile := path.Join("./public", req.URL.Path)
	http.ServeFile(res, req, htmfile)
}

func getfmt(fpath string) string {
	file, e := os.Open(path.Join("./public", fpath))
	if e != nil {
		return ".oat"
	}
	read, e := ioutil.ReadAll(file)
	if e != nil {
		return ".oat"
	}

	if strings.HasPrefix(string(read), "<!--fmt:omm-->") {
		return ".omm"
	} else if strings.HasPrefix(string(read), "<!--fmt:klr-->") {
		return ".klr"
	}

	return ".oat"
}

func handle(res http.ResponseWriter, req *http.Request) {

	//remove the extension, and replace it with .oat (or .omm or .klr)
	oatname := strings.TrimSuffix(req.URL.Path, filepath.Ext(req.URL.Path)) + getfmt(req.URL.Path)

	//prepend the server path
	oatf := path.Join("server", oatname)

	if _, f := os.Stat(oatf); !os.IsNotExist(f) {

		var tmp = params
		tmp.Name = oatname

		//load the oat (or omm or kayl) file using goat
		lib, e := goat.LoadLibrary(oatf, tmp)

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
