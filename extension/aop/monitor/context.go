package monitor

import (
	"sort"
	"sync"
	"time"

	"github.com/alibaba/ioc-golang/aop"
	"github.com/alibaba/ioc-golang/aop/common"
	monitorPB "github.com/alibaba/ioc-golang/extension/aop/monitor/api/ioc_golang/aop/monitor"
)

type context struct {
	SDID       string
	MethodName string
	Ch         chan *monitorPB.MonitorResponse
	inited     bool
	ticker     *time.Ticker
	stopCh     chan struct{}
	cache      sync.Map // methodUniqueName -> methodInvocationRecord
}

func newContext(sdid, method string, ch chan *monitorPB.MonitorResponse, period time.Duration) *context {
	newCtx := &context{
		SDID:       sdid,
		MethodName: method,
		Ch:         ch,
		stopCh:     make(chan struct{}),
	}
	newCtx.init(period)
	return newCtx
}

func (c *context) init(period time.Duration) {
	if !c.inited {
		c.inited = true
		c.ticker = time.NewTicker(period)
		go c.run()
	}
}

func (c *context) run() {
	for {
		select {
		case <-c.stopCh:
			return
		case <-c.ticker.C:
			// collect data
			monitorResponseItemSorter := make(monitorResponseItemsSorter, 0)
			c.cache.Range(func(key, value interface{}) bool {
				invocationMethodKey := key.(string)
				invocationMethodRecord := value.(*methodInvocationRecord)
				sdid, methodName := common.ParseSDIDAndMethodFromUniqueKey(invocationMethodKey)
				total, success, fail, agRT, failedRate := invocationMethodRecord.describeAndReset()
				if total == 0 {
					return true
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
				return true
			})
			sort.Sort(monitorResponseItemSorter)
			c.Ch <- &monitorPB.MonitorResponse{
				MonitorResponseItems: monitorResponseItemSorter,
			}
		}
	}
}

func (c *context) filterAndGetRecord(ctx *aop.InvocationContext) (*methodInvocationRecord, bool) {
	// filter invocations that is not monitored
	if c.SDID != "" && ctx.SDID != c.SDID {
		return nil, false
	} else if c.MethodName != "" && ctx.MethodName != c.MethodName {
		return nil, false
	}

	// monitor the invocation
	invocationMethodKey := common.GetMethodUniqueKey(ctx.SDID, ctx.MethodName)
	var methodRecord *methodInvocationRecord
	val, ok := c.cache.Load(invocationMethodKey)
	if !ok {
		methodRecord = newMethodInvocationRecord()
		c.cache.Store(invocationMethodKey, methodRecord)
	} else {
		methodRecord = val.(*methodInvocationRecord)
	}
	return methodRecord, true
}

func (c *context) beforeInvoke(ctx *aop.InvocationContext) {
	monitorMethodRecord, match := c.filterAndGetRecord(ctx)
	if match {
		monitorMethodRecord.beforeRequest(ctx)
	}
}

func (c *context) afterInvoke(ctx *aop.InvocationContext) {
	monitorMethodRecord, match := c.filterAndGetRecord(ctx)
	if match {
		monitorMethodRecord.afterRequest(ctx)
	}
}

func (c *context) destroy() {
	if c.inited {
		c.ticker.Stop()
		close(c.stopCh)
	}
}

type methodInvocationRecord struct {
	total   int
	success int
	fail    int
	rts     []int64

	grIDReqMap sync.Map // grID -> req before invoke time

	lock sync.RWMutex
}

func (m *methodInvocationRecord) describeAndReset() (int, int, int, float32, float32) {
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

func (m *methodInvocationRecord) beforeRequest(ctx *aop.InvocationContext) {
	m.grIDReqMap.Store(ctx.GrID, time.Now().UnixMilli())
}

func (m *methodInvocationRecord) afterRequest(ctx *aop.InvocationContext) {
	val, ok := m.grIDReqMap.LoadAndDelete(ctx.GrID)
	if !ok {
		return
	}
	startTime := val.(int64)
	duration := time.Now().UnixMilli() - startTime

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

func newMethodInvocationRecord() *methodInvocationRecord {
	return &methodInvocationRecord{
		rts: make([]int64, 0),
	}
}
