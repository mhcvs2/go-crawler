package base

import "net/http"

//Request structure include http request pointer and request depth
type Request struct {
	httpReq *http.Request
	depth   uint32
}

//NewRequest create a new request
func NewRequest(httpReq *http.Request, depth uint32) *Request {
	return &Request{httpReq: httpReq, depth: depth}
}

//HttpReq get http request
func (req *Request) HttpReq() *http.Request {
	return req.httpReq
}

//Depth get value of depth
func (req *Request) Depth() uint32 {
	return req.depth
}
