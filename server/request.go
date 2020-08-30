package undra

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/omm-lang/omm/lang/types"
)

type OmmHTTPRequest struct {
	Method     types.OmmString
	Url        types.OmmString
	Protocol   types.OmmString
	Header     types.OmmHash
	Body       types.OmmString
	Host       types.OmmString
	Form       types.OmmHash
	RemoteAddr types.OmmString
	RequestURI types.OmmString
}

func (r OmmHTTPRequest) FromGoType(req http.Request) {

	var rmethod types.OmmString
	rmethod.FromGoType(req.Method)
	var rurl types.OmmString
	rurl.FromGoType(fmt.Sprint(req.URL) /* makes the url struct into a human reable url */)
	var rproto types.OmmString
	rproto.FromGoType(req.Proto)

	var rbody types.OmmString
	body, _ := ioutil.ReadAll(req.Body)
	rbody.FromGoType(string(body))

	var host types.OmmString
	host.FromGoType(req.Host)
	var remoteaddr types.OmmString
	remoteaddr.FromGoType(req.RemoteAddr)
	var requesturi types.OmmString
	requesturi.FromGoType(req.RequestURI)
}

func (r OmmHTTPRequest) Format() string {
	return "{ undra-request }"
}

func (r OmmHTTPRequest) Type() string {
	return "undra-request"
}

func (r OmmHTTPRequest) TypeOf() string {
	return r.Type()
}

func (_ OmmHTTPRequest) Deallocate() {}
