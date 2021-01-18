package undra

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tusklang/tusk/lang/types"
)

func createRequest(req http.Request) *types.TuskType {

	reqobj := requestProto.New(types.Instance{}) //create a new object with an empty instance

	//convert to tusk types
	var method types.TuskString
	method.FromGoType(req.Method)
	rmethod, _ := reqobj.Get("method", "", "global")
	*rmethod = method

	var url types.TuskString
	url.FromGoType(fmt.Sprint(req.URL) /* makes the url struct into a human reable url */)
	rurl, _ := reqobj.Get("url", "", "global")
	*rurl = url

	var proto types.TuskString
	proto.FromGoType(req.Proto)
	rproto, _ := reqobj.Get("protocol", "", "global")
	*rproto = proto

	var body types.TuskString
	_body, _ := ioutil.ReadAll(req.Body)
	body.FromGoType(string(_body))
	rbody, _ := reqobj.Get("body", "", "global")
	*rbody = body

	var host types.TuskString
	host.FromGoType(req.Host)
	rhost, _ := reqobj.Get("host", "", "global")
	*rhost = host

	var form types.TuskHash
	rform, _ := reqobj.Get("form", "", "global")
	*rform = form

	var remoteaddr types.TuskString
	remoteaddr.FromGoType(req.RemoteAddr)
	rremoteaddr, _ := reqobj.Get("remoteaddr", "", "global")
	*rremoteaddr = remoteaddr

	var requesturi types.TuskString
	requesturi.FromGoType(req.RequestURI)
	rrequesturi, _ := reqobj.Get("requesturi", "", "global")
	*rrequesturi = requesturi
	//////////////////////

	var tusktype types.TuskType = reqobj
	return &tusktype
}
