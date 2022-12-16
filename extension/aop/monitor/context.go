package monitor

import (
	"sort"
	"sync"
	"time"

	"github.com/alibaba/ioc-golang/aop"
	"github.com/alibaba/ioc-golang/aop/common"
	monitorPB "github.com/alibaba/ioc-golang/extension/aop/monitor/api/ioc_golang/aop/monitor"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
// +ioc:autowire:paramType=contextParam
// +ioc:autowire:constructFunc=init
// +ioc:autowire:proxy:autoInjection=false

type context struct {
	sdid                                    string
	methodName                              string
	ch                                      chan *monitorPB.MonitorResponse
	ticker                                  *time.Ticker
	stopCh                                  chan struct{}
	methodUniqueNameInvocationRecordMap     map[string]methodInvocationRecordIOCInterface // methodUniqueName -> methodInvocationRecord
	methodUniqueNameInvocationRecordMapLock sync.Mutex
	destroyed                               bool
}

type contextParam struct {
	SDID       string
	MethodName string
	Ch         chan *monitorPB.MonitorResponse
	Period     time.Duration
}

func (p *contextParam) init(c *context) (*context, error) {
	c.ch = p.Ch
	c.methodName = p.MethodName
	c.sdid = p.SDID

	c.methodUniqueNameInvocationRecordMap = make(map[string]methodInvocationRecordIOCInterface)
	c.stopCh = make(chan struct{})
	c.destroyed = false
	c.ticker = time.NewTicker(p.Period)
	go c.run()
	return c, nil
}

func (c *context) run() {
	for {
		select {
		case <-c.stopCh:
			return
		case <-c.ticker.C:
			// collect data
			monitorResponseItemSorter := make(monitorResponseItemsSorter, 0)
			c.methodUniqueNameInvocationRecordMapLock.Lock()
			for invocationMethodKey, invocationMethodRecord := range c.methodUniqueNameInvocationRecordMap {
				sdid, methodName := common.ParseSDIDAndMethodFromUniqueKey(invocationMethodKey)
				total, success, fail, agRT, failedRate := invocationMethodRecord.DescribeAndReset()
				if total == 0 {
					continue
				}
				monitorResponseItemSorter = append(monitorResponseItemSorter, &monitorPB.MonitorResponseItem{
					Sdid:     sdid,
					Method:   methodName,
					Total:    int64(total),
					Success:  int64(success),
					Fail:     int64(fail),
					AvgRT:    agRT,
					FailRate: failedRate,
				})
			}
			c.methodUniqueNameInvocationRecordMapLock.Unlock()
			sort.Sort(monitorResponseItemSorter)
			c.ch <- &monitorPB.MonitorResponse{
				MonitorResponseItems: monitorResponseItemSorter,
			}
		}
	}
}

func (c *context) filterAndGetRecord(ctx *aop.InvocationContext) (methodInvocationRecordIOCInterface, bool) {
	// filter invocations that is not monitored
	if c.sdid != "" && ctx.SDID != c.sdid {
		return nil, false
	} else if c.methodName != "" && ctx.MethodName != c.methodName {
		return nil, false
	}

	// monitor the invocation
	invocationMethodKey := common.GetMethodUniqueKey(ctx.SDID, ctx.MethodName)
	c.methodUniqueNameInvocationRecordMapLock.Lock()
	defer c.methodUniqueNameInvocationRecordMapLock.Unlock()
	methodRecord, ok := c.methodUniqueNameInvocationRecordMap[invocationMethodKey]
	if !ok {
		methodRecord, _ = GetmethodInvocationRecordIOCInterface()
		c.methodUniqueNameInvocationRecordMap[invocationMethodKey] = methodRecord
	}
	return methodRecord, true
}

func (c *context) BeforeInvoke(ctx *aop.InvocationContext) {
	monitorMethodRecord, match := c.filterAndGetRecord(ctx)
	if match {
		monitorMethodRecord.BeforeRequest(ctx)
	}
}

func (c *context) AfterInvoke(ctx *aop.InvocationContext) {
	monitorMethodRecord, match := c.filterAndGetRecord(ctx)
	if match {
		monitorMethodRecord.AfterRequest(ctx)
	}
}

func (c *context) Destroy() {
	if !c.destroyed {
		c.destroyed = true
		c.ticker.Stop()
		close(c.stopCh)
	}
}

// +ioc:autowire=true
// +ioc:autowire:type=normal
// +ioc:autowire:constructFunc=newMethodInvocationRecord
// +ioc:autowire:proxy:autoInjection=false

type methodInvocationRecord struct {
	total   int
	success int
	fail    int
	rts     []int64

	grIDReqMap sync.Map // grID -> req before invoke time

	lock sync.RWMutex
}

func newMethodInvocationRecord(record *methodInvocationRecord) (*methodInvocationRecord, error) {
	record.rts = make([]int64, 0)
	return record, nil
}

func (m *methodInvocationRecord) DescribeAndReset() (int, int, int, float32, float32) {
	m.lock.Lock()
	defer m.lock.Unlock()

	agRT := getAverageInt64(m.rts)
	failedRate := float32(0)
	if m.fail > 0 {
		failedRate = float32(m.fail) / float32(m.total)
	}

	total := m.total
	success := m.success
	fail := m.fail

	m.total = 0
	m.success = 0
	m.fail = 0
	m.rts = make([]int64, 0)

	return total, success, fail, agRT, failedRate
}

func (m *methodInvocationRecord) BeforeRequest(ctx *aop.InvocationContext) {
	m.grIDReqMap.Store(ctx.ID, time.Now().UnixMicro())
}

func (m *methodInvocationRecord) AfterRequest(ctx *aop.InvocationContext) {
	val, ok := m.grIDReqMap.LoadAndDelete(ctx.ID)
	if !ok {
		return
	}
	startTime := val.(int64)
	duration := time.Now().UnixMicro() - startTime

	m.lock.Lock()
	defer m.lock.Unlock()

	m.rts = append(m.rts, duration)
	m.total += 1
	if isFailed, _ := common.IsInvocationFailed(ctx.ReturnValues); isFailed {
		m.fail += 1
	} else {
		m.success += 1
	}
}
