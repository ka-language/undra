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
	rmethod, _ := reqobj.Get("method", "")
	*rmethod = method

	var url types.TuskString
	url.FromGoType(fmt.Sprint(req.URL) /* makes the url struct into a human reable url */)
	rurl, _ := reqobj.Get("url", "")
	*rurl = url

	var proto types.TuskString
	proto.FromGoType(req.Proto)
	rproto, _ := reqobj.Get("protocol", "")
	*rproto = proto

	var body types.TuskString
	_body, _ := ioutil.ReadAll(req.Body)
	body.FromGoType(string(_body))
	rbody, _ := reqobj.Get("body", "")
	*rbody = body

	var host types.TuskString
	host.FromGoType(req.Host)
	rhost, _ := reqobj.Get("host", "")
	*rhost = host

	var form types.TuskHash
	rform, _ := reqobj.Get("form", "")
	*rform = form

	var remoteaddr types.TuskString
	remoteaddr.FromGoType(req.RemoteAddr)
	rremoteaddr, _ := reqobj.Get("remoteaddr", "")
	*rremoteaddr = remoteaddr

	var requesturi types.TuskString
	requesturi.FromGoType(req.RequestURI)
	rrequesturi, _ := reqobj.Get("requesturi", "")
	*rrequesturi = requesturi
	//////////////////////

	var tusktype types.TuskType = reqobj
	return &tusktype
}
