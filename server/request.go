package undra

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"omm/lang/types"
)

func createRequest(req http.Request) *types.OmmType {

	reqobj := requestProto.New(types.Instance{}) //create a new object with an empty instance

	//convert to omm types
	var method types.OmmString
	method.FromGoType(req.Method)
	rmethod, _ := reqobj.Get("method", "")
	*rmethod = method

	var url types.OmmString
	url.FromGoType(fmt.Sprint(req.URL) /* makes the url struct into a human reable url */)
	rurl, _ := reqobj.Get("url", "")
	*rurl = url

	var proto types.OmmString
	proto.FromGoType(req.Proto)
	rproto, _ := reqobj.Get("protocol", "")
	*rproto = proto

	var body types.OmmString
	_body, _ := ioutil.ReadAll(req.Body)
	body.FromGoType(string(_body))
	rbody, _ := reqobj.Get("body", "")
	*rbody = body

	var host types.OmmString
	host.FromGoType(req.Host)
	rhost, _ := reqobj.Get("host", "")
	*rhost = host

	var form types.OmmHash
	rform, _ := reqobj.Get("form", "")
	*rform = form

	var remoteaddr types.OmmString
	remoteaddr.FromGoType(req.RemoteAddr)
	rremoteaddr, _ := reqobj.Get("remoteaddr", "")
	*rremoteaddr = remoteaddr

	var requesturi types.OmmString
	requesturi.FromGoType(req.RequestURI)
	rrequesturi, _ := reqobj.Get("requesturi", "")
	*rrequesturi = requesturi
	//////////////////////

	var ommtype types.OmmType = reqobj
	return &ommtype
}
