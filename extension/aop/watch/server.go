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

package watch

import (
	"github.com/alibaba/ioc-golang/aop/common"
	"github.com/alibaba/ioc-golang/extension/aop/watch/api/ioc_golang/aop/watch"
	"github.com/alibaba/ioc-golang/logger"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:proxy=false

type watchService struct {
	watch.UnimplementedWatchServiceServer
	WatchInterceptor interceptorImplIOCInterface `singleton:""`
}

func (w *watchService) Watch(req *watch.WatchRequest, svr watch.WatchService_WatchServer) error {
	logger.Red("[Debug Server] Receive watch request %+v\n", req.String())
	defer logger.Red("[Debug Server] Watch request %s finished \n", req.String())
	sdid := req.GetSdid()
	method := req.GetMethod()
	sendCh := make(chan *watch.WatchResponse)
	maxDepth := 5
	maxLength := 1000
	if req.GetMaxDepth() != 0 {
		maxDepth = int(req.GetMaxDepth())
	}

	if req.GetMaxLength() != 0 {
		maxLength = int(req.GetMaxLength())
	}

	var fieldMatcher *common.FieldMatcher
	for _, matcher := range req.GetMatchers() {
		// todo multi match support
		fieldMatcher = &common.FieldMatcher{
			FieldIndex: int(matcher.Index),
			MatchRule:  matcher.GetMatchPath() + "=" + matcher.GetMatchValue(),
		}
	}

	watchCtx, err := GetcontextIOCInterfaceSingleton(&contextParam{
		SDID:         sdid,
		MethodName:   method,
		MaxLength:    maxLength,
		MaxDepth:     maxDepth,
		Ch:           sendCh,
		FieldMatcher: fieldMatcher,
	})
	if err != nil {
		return err
	}
	w.WatchInterceptor.Watch(watchCtx)

	done := svr.Context().Done()
	for {
		select {
		case <-done:
			// watch stop
			w.WatchInterceptor.UnWatch(watchCtx)
			return nil
		case watchRsp := <-sendCh:
			if err := svr.Send(watchRsp); err != nil {
				return err
			}
		}
	}
}
