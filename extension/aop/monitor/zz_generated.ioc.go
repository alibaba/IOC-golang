//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by iocli, run 'iocli gen' to re-generate

package monitor

import (
	"github.com/alibaba/ioc-golang/aop"
	autowire "github.com/alibaba/ioc-golang/autowire"
	normal "github.com/alibaba/ioc-golang/autowire/normal"
	singleton "github.com/alibaba/ioc-golang/autowire/singleton"
	util "github.com/alibaba/ioc-golang/autowire/util"
	aopmonitor "github.com/alibaba/ioc-golang/extension/aop/monitor/api/ioc_golang/aop/monitor"
)

func init() {
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &context_{}
		},
	})
	contextStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &context{}
		},
		ParamFactory: func() interface{} {
			var _ contextParamInterface = &contextParam{}
			return &contextParam{}
		},
		ConstructFunc: func(i interface{}, p interface{}) (interface{}, error) {
			param := p.(contextParamInterface)
			impl := i.(*context)
			return param.init(impl)
		},
		Metadata: map[string]interface{}{
			"aop":      map[string]interface{}{},
			"autowire": map[string]interface{}{},
		},
		DisableProxy: true,
	}
	normal.RegisterStructDescriptor(contextStructDescriptor)
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &methodInvocationRecord_{}
		},
	})
	methodInvocationRecordStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &methodInvocationRecord{}
		},
		ConstructFunc: func(i interface{}, _ interface{}) (interface{}, error) {
			impl := i.(*methodInvocationRecord)
			var constructFunc methodInvocationRecordConstructFunc = newMethodInvocationRecord
			return constructFunc(impl)
		},
		Metadata: map[string]interface{}{
			"aop":      map[string]interface{}{},
			"autowire": map[string]interface{}{},
		},
		DisableProxy: true,
	}
	normal.RegisterStructDescriptor(methodInvocationRecordStructDescriptor)
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &interceptorImpl_{}
		},
	})
	interceptorImplStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &interceptorImpl{}
		},
		Metadata: map[string]interface{}{
			"aop":      map[string]interface{}{},
			"autowire": map[string]interface{}{},
		},
		DisableProxy: true,
	}
	singleton.RegisterStructDescriptor(interceptorImplStructDescriptor)
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &monitorService_{}
		},
	})
	monitorServiceStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &monitorService{}
		},
		Metadata: map[string]interface{}{
			"aop":      map[string]interface{}{},
			"autowire": map[string]interface{}{},
		},
		DisableProxy: true,
	}
	singleton.RegisterStructDescriptor(monitorServiceStructDescriptor)
}

type contextParamInterface interface {
	init(impl *context) (*context, error)
}
type methodInvocationRecordConstructFunc func(impl *methodInvocationRecord) (*methodInvocationRecord, error)
type context_ struct {
	BeforeInvoke_ func(ctx *aop.InvocationContext)
	AfterInvoke_  func(ctx *aop.InvocationContext)
	Destroy_      func()
}

func (c *context_) BeforeInvoke(ctx *aop.InvocationContext) {
	c.BeforeInvoke_(ctx)
}

func (c *context_) AfterInvoke(ctx *aop.InvocationContext) {
	c.AfterInvoke_(ctx)
}

func (c *context_) Destroy() {
	c.Destroy_()
}

type methodInvocationRecord_ struct {
	DescribeAndReset_ func() (int, int, int, float32, float32)
	BeforeRequest_    func(ctx *aop.InvocationContext)
	AfterRequest_     func(ctx *aop.InvocationContext)
}

func (m *methodInvocationRecord_) DescribeAndReset() (int, int, int, float32, float32) {
	return m.DescribeAndReset_()
}

func (m *methodInvocationRecord_) BeforeRequest(ctx *aop.InvocationContext) {
	m.BeforeRequest_(ctx)
}

func (m *methodInvocationRecord_) AfterRequest(ctx *aop.InvocationContext) {
	m.AfterRequest_(ctx)
}

type interceptorImpl_ struct {
	BeforeInvoke_ func(ctx *aop.InvocationContext)
	AfterInvoke_  func(ctx *aop.InvocationContext)
	Monitor_      func(monitorCtx contextIOCInterface)
	StopMonitor_  func()
}

func (i *interceptorImpl_) BeforeInvoke(ctx *aop.InvocationContext) {
	i.BeforeInvoke_(ctx)
}

func (i *interceptorImpl_) AfterInvoke(ctx *aop.InvocationContext) {
	i.AfterInvoke_(ctx)
}

func (i *interceptorImpl_) Monitor(monitorCtx contextIOCInterface) {
	i.Monitor_(monitorCtx)
}

func (i *interceptorImpl_) StopMonitor() {
	i.StopMonitor_()
}

type monitorService_ struct {
	Monitor_ func(req *aopmonitor.MonitorRequest, svr aopmonitor.MonitorService_MonitorServer) error
}

func (m *monitorService_) Monitor(req *aopmonitor.MonitorRequest, svr aopmonitor.MonitorService_MonitorServer) error {
	return m.Monitor_(req, svr)
}

type contextIOCInterface interface {
	BeforeInvoke(ctx *aop.InvocationContext)
	AfterInvoke(ctx *aop.InvocationContext)
	Destroy()
}

type methodInvocationRecordIOCInterface interface {
	DescribeAndReset() (int, int, int, float32, float32)
	BeforeRequest(ctx *aop.InvocationContext)
	AfterRequest(ctx *aop.InvocationContext)
}

type interceptorImplIOCInterface interface {
	BeforeInvoke(ctx *aop.InvocationContext)
	AfterInvoke(ctx *aop.InvocationContext)
	Monitor(monitorCtx contextIOCInterface)
	StopMonitor()
}

type monitorServiceIOCInterface interface {
	Monitor(req *aopmonitor.MonitorRequest, svr aopmonitor.MonitorService_MonitorServer) error
}

var _contextSDID string

func Getcontext(p *contextParam) (*context, error) {
	if _contextSDID == "" {
		_contextSDID = util.GetSDIDByStructPtr(new(context))
	}
	i, err := normal.GetImpl(_contextSDID, p)
	if err != nil {
		return nil, err
	}
	impl := i.(*context)
	return impl, nil
}

func GetcontextIOCInterface(p *contextParam) (contextIOCInterface, error) {
	if _contextSDID == "" {
		_contextSDID = util.GetSDIDByStructPtr(new(context))
	}
	i, err := normal.GetImplWithProxy(_contextSDID, p)
	if err != nil {
		return nil, err
	}
	impl := i.(contextIOCInterface)
	return impl, nil
}

var _methodInvocationRecordSDID string

func GetmethodInvocationRecord() (*methodInvocationRecord, error) {
	if _methodInvocationRecordSDID == "" {
		_methodInvocationRecordSDID = util.GetSDIDByStructPtr(new(methodInvocationRecord))
	}
	i, err := normal.GetImpl(_methodInvocationRecordSDID, nil)
	if err != nil {
		return nil, err
	}
	impl := i.(*methodInvocationRecord)
	return impl, nil
}

func GetmethodInvocationRecordIOCInterface() (methodInvocationRecordIOCInterface, error) {
	if _methodInvocationRecordSDID == "" {
		_methodInvocationRecordSDID = util.GetSDIDByStructPtr(new(methodInvocationRecord))
	}
	i, err := normal.GetImplWithProxy(_methodInvocationRecordSDID, nil)
	if err != nil {
		return nil, err
	}
	impl := i.(methodInvocationRecordIOCInterface)
	return impl, nil
}

var _interceptorImplSDID string

func GetinterceptorImplSingleton() (*interceptorImpl, error) {
	if _interceptorImplSDID == "" {
		_interceptorImplSDID = util.GetSDIDByStructPtr(new(interceptorImpl))
	}
	i, err := singleton.GetImpl(_interceptorImplSDID, nil)
	if err != nil {
		return nil, err
	}
	impl := i.(*interceptorImpl)
	return impl, nil
}

func GetinterceptorImplIOCInterfaceSingleton() (interceptorImplIOCInterface, error) {
	if _interceptorImplSDID == "" {
		_interceptorImplSDID = util.GetSDIDByStructPtr(new(interceptorImpl))
	}
	i, err := singleton.GetImplWithProxy(_interceptorImplSDID, nil)
	if err != nil {
		return nil, err
	}
	impl := i.(interceptorImplIOCInterface)
	return impl, nil
}

type ThisinterceptorImpl struct {
}

func (t *ThisinterceptorImpl) This() interceptorImplIOCInterface {
	thisPtr, _ := GetinterceptorImplIOCInterfaceSingleton()
	return thisPtr
}

var _monitorServiceSDID string

func GetmonitorServiceSingleton() (*monitorService, error) {
	if _monitorServiceSDID == "" {
		_monitorServiceSDID = util.GetSDIDByStructPtr(new(monitorService))
	}
	i, err := singleton.GetImpl(_monitorServiceSDID, nil)
	if err != nil {
		return nil, err
	}
	impl := i.(*monitorService)
	return impl, nil
}

func GetmonitorServiceIOCInterfaceSingleton() (monitorServiceIOCInterface, error) {
	if _monitorServiceSDID == "" {
		_monitorServiceSDID = util.GetSDIDByStructPtr(new(monitorService))
	}
	i, err := singleton.GetImplWithProxy(_monitorServiceSDID, nil)
	if err != nil {
		return nil, err
	}
	impl := i.(monitorServiceIOCInterface)
	return impl, nil
}

type ThismonitorService struct {
}

func (t *ThismonitorService) This() monitorServiceIOCInterface {
	thisPtr, _ := GetmonitorServiceIOCInterfaceSingleton()
	return thisPtr
}
