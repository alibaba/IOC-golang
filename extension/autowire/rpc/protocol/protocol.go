package protocol

import "dubbo.apache.org/dubbo-go/v3/protocol"

type Protocol interface {
	Invoke(invocation protocol.Invocation) protocol.Result
	Export(invoker protocol.Invoker) protocol.Exporter
}
