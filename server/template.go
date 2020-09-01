package undra

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func templatedoc(filename string, template map[string]string) (string, error) {

	_f, e := ioutil.ReadFile(filename)

	if e != nil {
		return "", e
	}

	f := string(_f)

	f += "<script type=\"text/javascript\">"
	for k, v := range template {
		f += fmt.Sprintf("document.querySelector(\"%s\").innerHTML=\"%s\";", strings.ReplaceAll(k, "\\", "\\\\"), strings.ReplaceAll(v, "\\", "\\\\"))
	}
	f += "</script>"

	/*
		Example:
		<html>
			<p id="test"></p>
		</html>

		Render with {
			"#test": "test text"
		}

		Becomes:
		<html>
			<p id="test"></p>
		</html>
		<script>document.querySelector("#test").innerHTML = "test text";</script>
	*/

	return f, nil //for now
}
