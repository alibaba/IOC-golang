package protocol_impl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"dubbo.apache.org/dubbo-go/v3/common/constant"
	dubboProtocol "dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/protocol/invocation"

	"github.com/alibaba/ioc-golang/autowire/normal"
	"github.com/alibaba/ioc-golang/extension/autowire/rpc/protocol"
	"github.com/alibaba/ioc-golang/extension/normal/http_server"
	"github.com/alibaba/ioc-golang/extension/normal/http_server/ghttp"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
// +ioc:autowire:paramType=Param
// +ioc:autowire:constructFunc=Init

// IOCProtocol fixme we can choose 'ioc:autowire:constructFunc' = Init or 'IOCGoParam.Init' some
type IOCProtocol struct {
	HttpServer http_server.HttpServer
	inited     bool
	address    string
	exportPort string
}

func (i *IOCProtocol) Invoke(invocation dubboProtocol.Invocation) dubboProtocol.Result {
	sdID, _ := invocation.GetAttachment("sdid")
	data, _ := json.Marshal(invocation.Arguments())
	rsp, err := http.Post(DefaultSchema+"://"+i.address+"/"+sdID+"/"+invocation.MethodName(), DefaultContentType, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	rspData, _ := ioutil.ReadAll(rsp.Body)
	err = json.Unmarshal(rspData, invocation.Reply())
	if err != nil {
		panic(err)
	}
	return nil
}

func (i *IOCProtocol) Export(invoker dubboProtocol.Invoker) dubboProtocol.Exporter {
	if !i.inited {
		serverImpl, err := normal.GetImpl(`github.com/alibaba/ioc-golang/extension/normal/http_server.Impl`, &http_server.HTTPServerConfig{
			Port: i.exportPort,
		})
		if err != nil {
			panic(err)
		}
		i.HttpServer = serverImpl.(http_server.HttpServer)
		go func() {
			i.HttpServer.Run(context.Background())
		}()
		i.inited = true
	}

	sdid := invoker.GetURL().GetParam(constant.InterfaceKey, "")
	clientStubFullName := invoker.GetURL().GetParam("alias", "")
	svc := ServiceMap.GetServiceByServiceKey(IOCProtocolName, sdid)
	if svc == nil {
		return nil
	}

	for methodName, methodType := range svc.Method() {
		i.HttpServer.RegisterRouter(fmt.Sprintf("/%s/%s", clientStubFullName, methodName), func(controller *ghttp.GRegisterController) error {
			reqData, err := ioutil.ReadAll(controller.R.Body)
			if err != nil {
				return err
			}
			arguments, err := ParseArgs(methodType.ArgsType(), reqData)
			if err != nil {
				return err
			}
			controller.Rsp = invoker.Invoke(context.Background(),
				invocation.NewRPCInvocation(methodName, arguments, nil)).Result()
			return nil
		}, nil, nil, http.MethodPost)
	}

	return dubboProtocol.NewBaseExporter(sdid, invoker, nil)
}

var _ protocol.Protocol = &IOCProtocol{}

// GetIOCProtocol get extended protocol from ioc-golang API
func GetIOCProtocol(param *Param) (protocol.Protocol, error) {
	iocProtocolImpl, err := normal.GetImpl("github.com/alibaba/ioc-golang/extension/autowire/rpc/protocol/protocol_impl.IOCProtocol", &Param{
		Address:    param.Address,
		ExportPort: param.ExportPort,
	})
	if err != nil {
		return nil, err
	}
	return iocProtocolImpl.(protocol.Protocol), nil
}
