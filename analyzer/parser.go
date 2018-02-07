package analyzer

import (
	"net/http"
	"go-crawler/base"
)

type ParseResponse func(httpResp *http.Response, respDepth uint32) ([]base.Data, []error)
