//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by iocli, run 'iocli gen' to re-generate

package call

import (
	"github.com/sirupsen/logrus"

	"github.com/alibaba/ioc-golang/aop"
	autowire "github.com/alibaba/ioc-golang/autowire"
	normal "github.com/alibaba/ioc-golang/autowire/normal"
	singleton "github.com/alibaba/ioc-golang/autowire/singleton"
	util "github.com/alibaba/ioc-golang/autowire/util"
)

func init() {
	debugLogContextStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &debugLogContext{}
		},
		ParamFactory: func() interface{} {
			var _ debugLogContextParamInterface = &debugLogContextParam{}
			return &debugLogContextParam{}
		},
		ConstructFunc: func(i interface{}, p interface{}) (interface{}, error) {
			param := p.(debugLogContextParamInterface)
			impl := i.(*debugLogContext)
			return param.init(impl)
		},
		Metadata: map[string]interface{}{
			"aop":      map[string]interface{}{},
			"autowire": map[string]interface{}{},
		},
		DisableProxy: true,
	}
	normal.RegisterStructDescriptor(debugLogContextStructDescriptor)
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &logInterceptor_{}
		},
	})
	logInterceptorStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &logInterceptor{}
		},
		ParamFactory: func() interface{} {
			var _ logInterceptorParamsInterface = &logInterceptorParams{}
			return &logInterceptorParams{}
		},
		ConstructFunc: func(i interface{}, p interface{}) (interface{}, error) {
			param := p.(logInterceptorParamsInterface)
			impl := i.(*logInterceptor)
			return param.initLogInterceptor(impl)
		},
		Metadata: map[string]interface{}{
			"aop":      map[string]interface{}{},
			"autowire": map[string]interface{}{},
		},
		DisableProxy: true,
	}
	singleton.RegisterStructDescriptor(logInterceptorStructDescriptor)
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &logGoRoutineInterceptorFacadeCtx_{}
		},
	})
	logGoRoutineInterceptorFacadeCtxStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &logGoRoutineInterceptorFacadeCtx{}
		},
		ParamFactory: func() interface{} {
			var _ logGoRoutineInterceptorFacadeCtxParamInterface = &logGoRoutineInterceptorFacadeCtxParam{}
			return &logGoRoutineInterceptorFacadeCtxParam{}
		},
		ConstructFunc: func(i interface{}, p interface{}) (interface{}, error) {
			param := p.(logGoRoutineInterceptorFacadeCtxParamInterface)
			impl := i.(*logGoRoutineInterceptorFacadeCtx)
			return param.initLogGoRoutineInterceptorFacadeCtx(impl)
		},
		Metadata: map[string]interface{}{
			"aop":      map[string]interface{}{},
			"autowire": map[string]interface{}{},
		},
	}
	normal.RegisterStructDescriptor(logGoRoutineInterceptorFacadeCtxStructDescriptor)
	invocationCtxLogsGeneratorStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &InvocationCtxLogsGenerator{}
		},
		Metadata: map[string]interface{}{
			"aop":      map[string]interface{}{},
			"autowire": map[string]interface{}{},
		},
		DisableProxy: true,
	}
	singleton.RegisterStructDescriptor(invocationCtxLogsGeneratorStructDescriptor)
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &logrusIOCCtxHook_{}
		},
	})
	logrusIOCCtxHookStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &LogrusIOCCtxHook{}
		},
		ConstructFunc: func(i interface{}, _ interface{}) (interface{}, error) {
			impl := i.(*LogrusIOCCtxHook)
			var constructFunc LogrusIOCCtxHookConstructFunc = newLogrusIOCCtxHook
			return constructFunc(impl)
		},
		Metadata: map[string]interface{}{
			"aop": map[string]interface{}{},
			"autowire": map[string]interface{}{
				"common": map[string]interface{}{
					"loadAtOnce": true,
				},
			},
		},
		DisableProxy: true,
	}
	singleton.RegisterStructDescriptor(logrusIOCCtxHookStructDescriptor)
	logServiceImplStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &logServiceImpl{}
		},
		Metadata: map[string]interface{}{
			"aop":      map[string]interface{}{},
			"autowire": map[string]interface{}{},
		},
		DisableProxy: true,
	}
	singleton.RegisterStructDescriptor(logServiceImplStructDescriptor)
}

type debugLogContextParamInterface interface {
	init(impl *debugLogContext) (*debugLogContext, error)
}
type logInterceptorParamsInterface interface {
	initLogInterceptor(impl *logInterceptor) (*logInterceptor, error)
}
type logGoRoutineInterceptorFacadeCtxParamInterface interface {
	initLogGoRoutineInterceptorFacadeCtx(impl *logGoRoutineInterceptorFacadeCtx) (*logGoRoutineInterceptorFacadeCtx, error)
}
type LogrusIOCCtxHookConstructFunc func(impl *LogrusIOCCtxHook) (*LogrusIOCCtxHook, error)
type logInterceptor_ struct {
	BeforeInvoke_ func(ctx *aop.InvocationContext)
	AfterInvoke_  func(ctx *aop.InvocationContext)
	WatchLogs_    func(logCtx *debugLogContext)
	StopWatch_    func()
	NotifyLogs_   func(content string)
}

func (l *logInterceptor_) BeforeInvoke(ctx *aop.InvocationContext) {
	l.BeforeInvoke_(ctx)
}

func (l *logInterceptor_) AfterInvoke(ctx *aop.InvocationContext) {
	l.AfterInvoke_(ctx)
}

func (l *logInterceptor_) WatchLogs(logCtx *debugLogContext) {
	l.WatchLogs_(logCtx)
}

func (l *logInterceptor_) StopWatch() {
	l.StopWatch_()
}

func (l *logInterceptor_) NotifyLogs(content string) {
	l.NotifyLogs_(content)
}

type logGoRoutineInterceptorFacadeCtx_ struct {
	pushContent_  func(content string)
	BeforeInvoke_ func(ctx *aop.InvocationContext)
	AfterInvoke_  func(ctx *aop.InvocationContext)
	Type_         func() string
}

func (l *logGoRoutineInterceptorFacadeCtx_) pushContent(content string) {
	l.pushContent_(content)
}

func (l *logGoRoutineInterceptorFacadeCtx_) BeforeInvoke(ctx *aop.InvocationContext) {
	l.BeforeInvoke_(ctx)
}

func (l *logGoRoutineInterceptorFacadeCtx_) AfterInvoke(ctx *aop.InvocationContext) {
	l.AfterInvoke_(ctx)
}

func (l *logGoRoutineInterceptorFacadeCtx_) Type() string {
	return l.Type_()
}

type logrusIOCCtxHook_ struct {
	Levels_      func() []logrus.Level
	Fire_        func(entry *logrus.Entry) error
	SetLogLevel_ func(level uint32)
}

func (l *logrusIOCCtxHook_) Levels() []logrus.Level {
	return l.Levels_()
}

func (l *logrusIOCCtxHook_) Fire(entry *logrus.Entry) error {
	return l.Fire_(entry)
}

func (l *logrusIOCCtxHook_) SetLogLevel(level uint32) {
	l.SetLogLevel_(level)
}

type logInterceptorIOCInterface interface {
	BeforeInvoke(ctx *aop.InvocationContext)
	AfterInvoke(ctx *aop.InvocationContext)
	WatchLogs(logCtx *debugLogContext)
	StopWatch()
	NotifyLogs(content string)
}

type logGoRoutineInterceptorFacadeCtxIOCInterface interface {
	pushContent(content string)
	BeforeInvoke(ctx *aop.InvocationContext)
	AfterInvoke(ctx *aop.InvocationContext)
	Type() string
}

type LogrusIOCCtxHookIOCInterface interface {
	Levels() []logrus.Level
	Fire(entry *logrus.Entry) error
	SetLogLevel(level uint32)
}

var _debugLogContextSDID string

func GetdebugLogContext(p *debugLogContextParam) (*debugLogContext, error) {
	if _debugLogContextSDID == "" {
		_debugLogContextSDID = util.GetSDIDByStructPtr(new(debugLogContext))
	}
	i, err := normal.GetImpl(_debugLogContextSDID, p)
	if err != nil {
		return nil, err
	}
	impl := i.(*debugLogContext)
	return impl, nil
}

var _logInterceptorSDID string

func GetlogInterceptorSingleton(p *logInterceptorParams) (*logInterceptor, error) {
	if _logInterceptorSDID == "" {
		_logInterceptorSDID = util.GetSDIDByStructPtr(new(logInterceptor))
	}
	i, err := singleton.GetImpl(_logInterceptorSDID, p)
	if err != nil {
		return nil, err
	}
	impl := i.(*logInterceptor)
	return impl, nil
}

func GetlogInterceptorIOCInterfaceSingleton(p *logInterceptorParams) (logInterceptorIOCInterface, error) {
	if _logInterceptorSDID == "" {
		_logInterceptorSDID = util.GetSDIDByStructPtr(new(logInterceptor))
	}
	i, err := singleton.GetImplWithProxy(_logInterceptorSDID, p)
	if err != nil {
		return nil, err
	}
	impl := i.(logInterceptorIOCInterface)
	return impl, nil
}

type ThislogInterceptor struct {
}

func (t *ThislogInterceptor) This() logInterceptorIOCInterface {
	thisPtr, _ := GetlogInterceptorIOCInterfaceSingleton(nil)
	return thisPtr
}

var _logGoRoutineInterceptorFacadeCtxSDID string

func GetlogGoRoutineInterceptorFacadeCtx(p *logGoRoutineInterceptorFacadeCtxParam) (*logGoRoutineInterceptorFacadeCtx, error) {
	if _logGoRoutineInterceptorFacadeCtxSDID == "" {
		_logGoRoutineInterceptorFacadeCtxSDID = util.GetSDIDByStructPtr(new(logGoRoutineInterceptorFacadeCtx))
	}
	i, err := normal.GetImpl(_logGoRoutineInterceptorFacadeCtxSDID, p)
	if err != nil {
		return nil, err
	}
	impl := i.(*logGoRoutineInterceptorFacadeCtx)
	return impl, nil
}

func GetlogGoRoutineInterceptorFacadeCtxIOCInterface(p *logGoRoutineInterceptorFacadeCtxParam) (logGoRoutineInterceptorFacadeCtxIOCInterface, error) {
	if _logGoRoutineInterceptorFacadeCtxSDID == "" {
		_logGoRoutineInterceptorFacadeCtxSDID = util.GetSDIDByStructPtr(new(logGoRoutineInterceptorFacadeCtx))
	}
	i, err := normal.GetImplWithProxy(_logGoRoutineInterceptorFacadeCtxSDID, p)
	if err != nil {
		return nil, err
	}
	impl := i.(logGoRoutineInterceptorFacadeCtxIOCInterface)
	return impl, nil
}

var _invocationCtxLogsGeneratorSDID string

func GetInvocationCtxLogsGeneratorSingleton() (*InvocationCtxLogsGenerator, error) {
	if _invocationCtxLogsGeneratorSDID == "" {
		_invocationCtxLogsGeneratorSDID = util.GetSDIDByStructPtr(new(InvocationCtxLogsGenerator))
	}
	i, err := singleton.GetImpl(_invocationCtxLogsGeneratorSDID, nil)
	if err != nil {
		return nil, err
	}
	impl := i.(*InvocationCtxLogsGenerator)
	return impl, nil
}

var _logrusIOCCtxHookSDID string

func GetLogrusIOCCtxHookSingleton() (*LogrusIOCCtxHook, error) {
	if _logrusIOCCtxHookSDID == "" {
		_logrusIOCCtxHookSDID = util.GetSDIDByStructPtr(new(LogrusIOCCtxHook))
	}
	i, err := singleton.GetImpl(_logrusIOCCtxHookSDID, nil)
	if err != nil {
		return nil, err
	}
	impl := i.(*LogrusIOCCtxHook)
	return impl, nil
}

func GetLogrusIOCCtxHookIOCInterfaceSingleton() (LogrusIOCCtxHookIOCInterface, error) {
	if _logrusIOCCtxHookSDID == "" {
		_logrusIOCCtxHookSDID = util.GetSDIDByStructPtr(new(LogrusIOCCtxHook))
	}
	i, err := singleton.GetImplWithProxy(_logrusIOCCtxHookSDID, nil)
	if err != nil {
		return nil, err
	}
	impl := i.(LogrusIOCCtxHookIOCInterface)
	return impl, nil
}

type ThisLogrusIOCCtxHook struct {
}

func (t *ThisLogrusIOCCtxHook) This() LogrusIOCCtxHookIOCInterface {
	thisPtr, _ := GetLogrusIOCCtxHookIOCInterfaceSingleton()
	return thisPtr
}

var _logServiceImplSDID string

func GetlogServiceImplSingleton() (*logServiceImpl, error) {
	if _logServiceImplSDID == "" {
		_logServiceImplSDID = util.GetSDIDByStructPtr(new(logServiceImpl))
	}
	i, err := singleton.GetImpl(_logServiceImplSDID, nil)
	if err != nil {
		return nil, err
	}
	impl := i.(*logServiceImpl)
	return impl, nil
}
