package analyzer

import (
	log "github.com/Sirupsen/logrus"
	"go-crawler/base"
	"go-crawler/middleware"
	"errors"
	"net/url"
	"fmt"
)

var analyzerIdGennerator = middleware.NewIdGenerator()

func genAnalyzerId() uint32 {
	return analyzerIdGennerator.GetUint32()
}

//----------------------------------------------------

type Analyzer interface {
	Id() uint32
	Analyze(
		respParsers []ParseResponse,
		resp base.Response,
	) ([]base.Data, []error)
}

func NewAnalyzer() Analyzer {
	return &myAnalyzer{id: genAnalyzerId()}
}

type myAnalyzer struct {
	id uint32
}

func (analyzer *myAnalyzer) Id() uint32 {
	return analyzer.id
}

func (analyzer *myAnalyzer) Analyze(
	respParsers []ParseResponse,
	resp base.Response,
) (dataList []base.Data, errorList []error) {
	if respParsers == nil {
		err := errors.New("the response parser list is invalid")
		return nil, []error{err}
	}
	httpResp := resp.HttpResp()
	if httpResp == nil {
		err := errors.New("the http response is invalid")
		return nil, []error{err}
	}
	var reqUrl *url.URL = httpResp.Request.URL
	log.Infof("Parse the response (reqUrl=%s)... \n", reqUrl)
	respDepth := resp.Depth()

	dataList = make([]base.Data, 0)
	errorList = make([]error, 0)

	for i, respParser := range respParsers {
		if respParser == nil {
			err := errors.New(fmt.Sprintf("the document parser [%d] id invalid", i))
			errorList = append(errorList, err)
			continue
		}
		pDataList, pErrorList := respParser(httpResp, respDepth)
		if pDataList != nil {
			for _, pdata := range pDataList {
				dataList = appendDataList(dataList, pdata, respDepth)
			}
		}
		if pErrorList != nil {
			for _, pError := range pErrorList {
				errorList = appendErrorList(errorList, pError)
			}
		}
	}

	return dataList, errorList
}

func appendDataList(dataList []base.Data, data base.Data, respDepth uint32) []base.Data {
	if data == nil {
		return dataList
	}
	req, ok := data.(*base.Request)
	if !ok {
		return append(dataList, data)
	}
	newDepth := respDepth + 1
	if req.Depth() != newDepth {
		req = base.NewRequest(req.HttpReq(), newDepth)
	}
	return append(dataList, req)
}

func appendErrorList(errorList []error, err error) []error {
	if err == nil {
		return errorList
	}
	return append(errorList, err)
}
