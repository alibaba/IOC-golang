/*
 * Copyright (c) 2022, Alibaba Group;
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
