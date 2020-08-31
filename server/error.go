package undra

import (
	"net/http"
	"os"
	"path"
)

func serve_errs(res http.ResponseWriter, req *http.Request) bool {

	//detect all 4** errors
	//for now it only detect a 404

	//404
	htmfile := path.Join("./public", req.URL.Path)
	if _, e := os.Stat(htmfile); os.IsNotExist(e) { //if the file does not exist, send a 404
		http.ServeFile(res, req, "status404.html")
		return true
	}

	return false
}
