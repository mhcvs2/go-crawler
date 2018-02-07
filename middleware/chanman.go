package middleware

import (
	"errors"
	"fmt"
	"go-crawler/base"
	"sync"
)

type ChannelManagerStatus uint8

const (
	CHANNEL_MANAGER_STATUS_UNINITIALEZED ChannelManagerStatus = iota
	CHANNEL_MANAGER_STATUS_INITIALEZED
	CHANNEL_MANAGER_STATUS_CLOSED
)

var statusNameMap = map[ChannelManagerStatus]string{
	CHANNEL_MANAGER_STATUS_UNINITIALEZED: "uninitialized",
	CHANNEL_MANAGER_STATUS_INITIALEZED:   "initialized",
	CHANNEL_MANAGER_STATUS_CLOSED:        "closed",
}

type ChannelManager interface {
	Init(channelArgs base.ChannelArgs, reset bool) bool
	Close() bool
	ReqChan() (chan base.Request, error)
	RespChan() (chan base.Response, error)
	ItemChan() (chan base.Item, error)
	ErrorChan() (chan error, error)
	Status() ChannelManagerStatus
	Summary() string
}

func NewChannelManager(args base.ChannelArgs) ChannelManager {
	chanman := &myChannelManager{}
	chanman.Init(args, true)
	return chanman
}

// 通道管理器的实现类型。
type myChannelManager struct {
	channelArgs base.ChannelArgs     // 通道参数的容器。
	reqCh       chan base.Request    // 请求通道。
	respCh      chan base.Response   // 响应通道。
	itemCh      chan base.Item       // 条目通道。
	errorCh     chan error           // 错误通道。
	status      ChannelManagerStatus // 通道管理器的状态。
	rwmutex     sync.RWMutex         // 读写锁。
}

func (chanman *myChannelManager) Init(channelArgs base.ChannelArgs, reset bool) bool {
	if err := channelArgs.Check(); err != nil {
		panic(err)
	}
	chanman.rwmutex.Lock()
	defer chanman.rwmutex.Unlock()
	if chanman.status == CHANNEL_MANAGER_STATUS_INITIALEZED && !reset {
		return false
	}
	chanman.channelArgs = channelArgs
	chanman.reqCh = make(chan base.Request, channelArgs.ReqChanLen())
	chanman.respCh = make(chan base.Response, channelArgs.RespChanLen())
	chanman.itemCh = make(chan base.Item, channelArgs.ItemChanLen())
	chanman.errorCh = make(chan error, channelArgs.ErrorChanLen())
	chanman.status = CHANNEL_MANAGER_STATUS_INITIALEZED
	return true
}

func (chanman *myChannelManager) Close() bool {
	chanman.rwmutex.Lock()
	defer chanman.rwmutex.Unlock()
	if chanman.status != CHANNEL_MANAGER_STATUS_INITIALEZED {
		return false
	}
	close(chanman.reqCh)
	close(chanman.respCh)
	close(chanman.itemCh)
	close(chanman.errorCh)
	chanman.status = CHANNEL_MANAGER_STATUS_CLOSED
	return true
}

func (chanman *myChannelManager) ReqChan() (chan base.Request, error) {
	chanman.rwmutex.RLock()
	defer chanman.rwmutex.RUnlock()
	if err := chanman.checkStatus(); err != nil {
		return nil, err
	}
	return chanman.reqCh, nil
}

func (chanman *myChannelManager) RespChan() (chan base.Response, error) {
	chanman.rwmutex.RLock()
	defer chanman.rwmutex.RUnlock()
	if err := chanman.checkStatus(); err != nil {
		return nil, err
	}
	return chanman.respCh, nil
}

func (chanman *myChannelManager) ItemChan() (chan base.Item, error) {
	chanman.rwmutex.RLock()
	defer chanman.rwmutex.RUnlock()
	if err := chanman.checkStatus(); err != nil {
		return nil, err
	}
	return chanman.itemCh, nil
}

func (chanman *myChannelManager) ErrorChan() (chan error, error) {
	chanman.rwmutex.RLock()
	defer chanman.rwmutex.RUnlock()
	if err := chanman.checkStatus(); err != nil {
		return nil, err
	}
	return chanman.errorCh, nil
}

// 检查状态。在获取通道的时候，通道管理器应处于已初始化状态。
// 如果通道管理器未处于已初始化状态，那么本方法将会返回一个非nil的错误值。
func (chanman *myChannelManager) checkStatus() error {
	if chanman.status == CHANNEL_MANAGER_STATUS_INITIALEZED {
		return nil
	}
	statusName, ok := statusNameMap[chanman.status]
	if !ok {
		statusName = fmt.Sprintf("%d", chanman.status)
	}
	errMsg :=
		fmt.Sprintf("The undesirable status of channel manager: %s!\n",
			statusName)
	return errors.New(errMsg)
}

func (chanman *myChannelManager) Status() ChannelManagerStatus {
	return chanman.status
}

var chanmanSummaryTemplate = "status: %s, " +
	"requestChannel: %d/%d, " +
	"responseChannel: %d/%d, " +
	"itemChannel: %d/%d, " +
	"errorChannel: %d/%d"

func (chanman *myChannelManager) Summary() string {
	summary := fmt.Sprintf(chanmanSummaryTemplate,
		statusNameMap[chanman.status],
		len(chanman.reqCh), cap(chanman.reqCh),
		len(chanman.respCh), cap(chanman.respCh),
		len(chanman.itemCh), cap(chanman.itemCh),
		len(chanman.errorCh), cap(chanman.errorCh))
	return summary
}
