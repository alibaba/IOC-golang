package protocol_impl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	tracer2 "github.com/alibaba/ioc-golang/debug/interceptor/trace"

	"github.com/opentracing/opentracing-go"

	"github.com/alibaba/ioc-golang/autowire/util"
	"github.com/alibaba/ioc-golang/debug/interceptor/trace"

	"github.com/fatih/color"

	"dubbo.apache.org/dubbo-go/v3/common/constant"
	dubboProtocol "dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/protocol/invocation"

	"github.com/alibaba/ioc-golang/common"
	"github.com/alibaba/ioc-golang/extension/autowire/rpc/protocol"
	"github.com/alibaba/ioc-golang/extension/normal/http_server"
	"github.com/alibaba/ioc-golang/extension/normal/http_server/ghttp"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
// +ioc:autowire:paramType=Param
// +ioc:autowire:constructFunc=Init

// IOCProtocol is ioc protocol impl
type IOCProtocol struct {
	httpServer http_server.ImplIOCInterface
	address    string
	exportPort string
	timeout    string
}

func (i *IOCProtocol) Invoke(invocation dubboProtocol.Invocation) dubboProtocol.Result {
	sdID, _ := invocation.GetAttachment("sdid")
	data, _ := json.Marshal(invocation.Arguments())
	invokeURL := DefaultSchema + "://" + i.address + "/" + sdID + "/" + invocation.MethodName()

	timeoutDuration, err := time.ParseDuration(i.timeout)
	if err != nil {
		timeoutDuration = time.Second * 3
	}
	requestCtx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()
	req, err := http.NewRequestWithContext(requestCtx, http.MethodPost, invokeURL, bytes.NewBuffer(data))
	if err != nil {
		return nil
	}

	// inject tracing context if necessary
	if currentSpan := tracer2.GetTraceInterceptor().GetCurrentSpan(); currentSpan != nil {
		// current rpc invocation is in tracing link
		carrier := opentracing.HTTPHeadersCarrier(req.Header)
		_ = trace.GetGlobalTracer().Inject(currentSpan.Context(), opentracing.HTTPHeaders, carrier)
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		color.Red("[IOC Protocol] Invoke %s with error = %s", invokeURL, err)
		return &dubboProtocol.RPCResult{
			Err: err,
		}
	}
	rspData, _ := ioutil.ReadAll(rsp.Body)
	replyList := invocation.Reply().(*[]interface{})
	finalIsError := false
	finalErrorNotNil := false
	if length := len(*replyList); length > 0 {
		_, ok := (*replyList)[length-1].(*error)
		if ok {
			finalIsError = true
		}
	}
	err = json.Unmarshal(rspData, replyList)
	if err != nil && finalIsError {
		// try to recover unmarshal failed caused by error not empty, first try to unmarshal to string
		(*replyList)[len(*replyList)-1] = ""
		err = json.Unmarshal(rspData, replyList)
		if err == nil {
			// previous unmarshal failed is caused by error not empty, mark final error not nil
			finalErrorNotNil = true
		}
	}
	if err != nil {
		return &dubboProtocol.RPCResult{
			Err: err,
		}
	}
	if finalErrorNotNil {
		realErr := fmt.Errorf((*replyList)[len(*replyList)-1].(string))
		(*replyList)[len(*replyList)-1] = &realErr
	}
	return nil
}

func (i *IOCProtocol) Export(invoker dubboProtocol.Invoker) dubboProtocol.Exporter {
	if i.httpServer == nil {
		i.httpServer = getHTTPServerSingleton(i.exportPort)
	}

	sdid := invoker.GetURL().GetParam(constant.InterfaceKey, "")
	clientStubFullName := invoker.GetURL().GetParam(common.AliasKey, "")
	svc := ServiceMap.GetServiceByServiceKey(IOCProtocolName, sdid)
	if svc == nil {
		return nil
	}

	for methodName, methodType := range svc.Method() {
		argsType := methodType.ArgsType()
		tempMethod := methodName
		i.httpServer.RegisterRouter(fmt.Sprintf("/%s/%s", clientStubFullName, tempMethod), func(controller *ghttp.GRegisterController) error {
			reqData, err := ioutil.ReadAll(controller.R.Body)
			if err != nil {
				return err
			}
			arguments, err := ParseArgs(argsType, reqData)
			if err != nil {
				return err
			}

			carrier := opentracing.HTTPHeadersCarrier(controller.R.Header)
			clientContext, err := trace.GetGlobalTracer().Extract(opentracing.HTTPHeaders, carrier)
			if err == nil {
				traceCtx := &trace.Context{
					SDID:              util.ToRPCServiceSDID(clientStubFullName),
					MethodName:        tempMethod,
					ClientSpanContext: clientContext,
				}
				trace.GetTraceInterceptor().TraceThisGR(traceCtx)
				defer trace.GetTraceInterceptor().UnTrace(traceCtx)
			}
			controller.Rsp = invoker.Invoke(context.Background(),
				invocation.NewRPCInvocation(tempMethod, arguments, nil)).Result()
			return nil
		}, nil, nil, http.MethodPost)
	}

	return dubboProtocol.NewBaseExporter(sdid, invoker, nil)
}

var _ protocol.Protocol = &IOCProtocol{}
