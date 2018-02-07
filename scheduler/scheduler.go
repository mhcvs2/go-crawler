package scheduler

import (
	"errors"
	"fmt"
	"github.com/prometheus/common/log"
	anlz "go-crawler/analyzer"
	"go-crawler/base"
	dl "go-crawler/downloader"
	ipl "go-crawler/itempipeline"
	mdw "go-crawler/middleware"
	"net/http"
	"sync/atomic"
)

const (
	DOWNLOADER_CODE   = "downloader"
	ANALYZER_CODE     = "analyzer"
	ITEMPIPELINE_CODE = "item_pipeline"
	SECHDULER_CODE    = "scheduler"
)

type GenHttpClient func() *http.Client

type Scheduler interface {
	Start(
		channelArgs base.ChannelArgs,
		poolBaseArgs base.PoolBaseArgs,
		crawlDepth uint32,
		httpClientGenerator GenHttpClient,
		respParsers []anlz.ParseResponse,
		itemProcessors []ipl.ProcessItem,
		firstHrrpReq *http.Request,
	) (err error)
	Stop() bool
	Running() bool
	ErrorChan() <-chan error
	Idle() bool
	Summary(prefix string) SchedSummary
}

func NewScheduler() Scheduler {
	return &myScheduler{}
}

type myScheduler struct {
	channelArgs   base.ChannelArgs
	poolBaseArgs  base.PoolBaseArgs
	crawlDepth    uint32
	primaryDomain string
	chanman       mdw.ChannelManager
	stopSign      mdw.StopSign
	dlpool        dl.PageDownloaderPool
	analyzerPool  anlz.AnalyzerPool
	itemPipeline  ipl.ItemPipeline
	reqCache      requestCache
	urlMap        map[string]bool
	running       uint32
}

func (sched *myScheduler) Start(
	channelArgs base.ChannelArgs,
	poolBaseArgs base.PoolBaseArgs,
	crawlDepth uint32,
	httClientGenerator GenHttpClient,
	respParsers []anlz.ParseResponse,
	itemProcessors []ipl.ProcessItem,
	firstHrrpReq *http.Request,
) (err error) {
	defer func() {
		if p := recover(); p != nil {
			errMsg := fmt.Sprintf("Fatal Scheduler Error: %s\n", p)
			log.Fatal(errMsg)
			err = errors.New(errMsg)
		}
	}()
	if atomic.LoadUint32(&sched.running) == 1 {
		return errors.New("the scheduler has been started\n")
	}
	atomic.StoreUint32(&sched.running, 1)

	if err := channelArgs.Check(); err != nil {
		return err
	}
	sched.channelArgs = channelArgs
	if err := poolBaseArgs.Check(); err != nil {
		return err
	}
	sched.poolBaseArgs = poolBaseArgs
	sched.crawlDepth = crawlDepth

	sched.chanman = generateChannelManager(sched.channelArgs)
	if httClientGenerator == nil {
		return errors.New("the http client generator list is invalid")
	}
	dlpool, err :=
		generatePageDownloaderPool(
			sched.poolBaseArgs.PageDownloaderPoolSize(),
			httClientGenerator,
		)
	if err != nil {
		errMsg :=
			fmt.Sprintf("Occur error when get page downloader pool: %s\n", err)
		return errors.New(errMsg)
	}
	sched.dlpool = dlpool
	analyzerPool, err := generateAnalyzerPool(sched.poolBaseArgs.AnalyzerPoolSize())
	if err != nil {
		errMsg :=
			fmt.Sprintf("Occur error when get analyzer pool: %s\n", err)
		return errors.New(errMsg)
	}
	sched.analyzerPool = analyzerPool

	if itemProcessors == nil {
		return errors.New("The item processor list is invalid!")
	}
	for i, ip := range itemProcessors {
		if ip == nil {
			return errors.New(fmt.Sprintf("The %dth item processor is invalid!", i))
		}
	}
	sched.itemPipeline = generateItemPipeline(itemProcessors)

	if sched.stopSign == nil {
		sched.stopSign = mdw.NewStopSign()
	} else {
		sched.stopSign.Reset()
	}

	sched.reqCache = newRequestCache()
	sched.urlMap = make(map[string]bool)




}
