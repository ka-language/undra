package undra

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"ka/lang/types"
)

func createRequest(req http.Request) *types.KaType {

	reqobj := requestProto.New(types.Instance{}) //create a new object with an empty instance

	//convert to ka types
	var method types.KaString
	method.FromGoType(req.Method)
	rmethod, _ := reqobj.Get("method", "")
	*rmethod = method

	var url types.KaString
	url.FromGoType(fmt.Sprint(req.URL) /* makes the url struct into a human reable url */)
	rurl, _ := reqobj.Get("url", "")
	*rurl = url

	var proto types.KaString
	proto.FromGoType(req.Proto)
	rproto, _ := reqobj.Get("protocol", "")
	*rproto = proto

	var body types.KaString
	_body, _ := ioutil.ReadAll(req.Body)
	body.FromGoType(string(_body))
	rbody, _ := reqobj.Get("body", "")
	*rbody = body

	var host types.KaString
	host.FromGoType(req.Host)
	rhost, _ := reqobj.Get("host", "")
	*rhost = host

	var form types.KaHash
	rform, _ := reqobj.Get("form", "")
	*rform = form

	var remoteaddr types.KaString
	remoteaddr.FromGoType(req.RemoteAddr)
	rremoteaddr, _ := reqobj.Get("remoteaddr", "")
	*rremoteaddr = remoteaddr

	var requesturi types.KaString
	requesturi.FromGoType(req.RequestURI)
	rrequesturi, _ := reqobj.Get("requesturi", "")
	*rrequesturi = requesturi
	//////////////////////

	var katype types.KaType = reqobj
	return &katype
}
