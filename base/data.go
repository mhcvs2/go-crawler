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

func (req *Request) Valid() bool {
	return req.httpReq != nil && req.httpReq.URL != nil
}

//-------------------------------------------------------------------
//Response is http response
type Response struct {
	httpResp *http.Response
	depth    uint32
}

//NewResponse return a Response instance
func NewResponse(httpResp *http.Response, depth uint32) *Response {
	return &Response{httpResp: httpResp, depth: depth}
}

func (resp *Response) HttpResp() *http.Response {
	return resp.httpResp
}

func (resp *Response) Depth() uint32 {
	return resp.depth
}

func (resp *Response) Valid() bool {
	return resp.httpResp != nil && resp.httpResp.Body != nil
}

//-------------------------------------------------------------------
//Item is item data
type Item map[string]interface{}

func (item Item) Valid() bool {
	return item != nil
}

//-------------------------------------------------------------------
//Data valid data
type Data interface {
	Valid() bool
}
