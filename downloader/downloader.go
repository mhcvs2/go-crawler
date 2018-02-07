package downloader

import (
	log "github.com/Sirupsen/logrus"
	"go-crawler/base"
	"go-crawler/middleware"
	"net/http"
)

var downloaderIdGennerator = middleware.NewIdGenerator()

func genDownloaderId() uint32 {
	return downloaderIdGennerator.GetUint32()
}

// ---------------------------------------------------

type PageDownloader interface {
	Id() uint32
	Download(req base.Request) (*base.Response, error)
}

func NewPageDownloader(client *http.Client) PageDownloader {
	id := genDownloaderId()
	if client == nil {
		client = &http.Client{}
	}
	return &myPageDownloader{
		id:         id,
		httpClient: *client,
	}
}

type myPageDownloader struct {
	id         uint32
	httpClient http.Client
}

func (dl *myPageDownloader) Id() uint32 {
	return dl.id
}

func (dl *myPageDownloader) Download(req base.Request) (*base.Response, error) {
	httpReq := req.HttpReq()
	log.Infof("Do the request (url=%s)...\n", httpReq.URL)
	httpResp, err := dl.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	return base.NewResponse(httpResp, req.Depth()), nil
}
