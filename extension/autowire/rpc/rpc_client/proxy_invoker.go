package rpc_client

import (
	"context"

	"dubbo.apache.org/dubbo-go/v3/common"
	dubboProtocol "dubbo.apache.org/dubbo-go/v3/protocol"

	"github.com/alibaba/ioc-golang/extension/autowire/rpc/protocol"
)

type proxyInvoker struct {
	protocol protocol.Protocol
	sdid     string
}

func (i *proxyInvoker) GetURL() *common.URL {
	return nil
}

func (i *proxyInvoker) IsAvailable() bool {
	return true
}

func (i *proxyInvoker) Destroy() {

}

func (i *proxyInvoker) Invoke(ctx context.Context, invocation dubboProtocol.Invocation) dubboProtocol.Result {
	invocation.SetAttachment("sdid", i.sdid)
	return i.protocol.Invoke(invocation)
}

func newProxyInvoker(protocol protocol.Protocol, sdid string) dubboProtocol.Invoker {
	return &proxyInvoker{protocol: protocol, sdid: sdid}
}
