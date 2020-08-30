package undra

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func handle(res http.ResponseWriter, req *http.Request) {
	//remove the extension, and replace it with .oat
	oatname := strings.TrimSuffix(req.URL.Path, filepath.Ext(req.URL.Path)) + ".oat"
	//load the oat file using goat
	oatf := path.Join("server", oatname)

	if _, f := os.Stat(oatf); !os.IsNotExist(f) {

	} else {
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
}
