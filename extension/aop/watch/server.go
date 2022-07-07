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
	"github.com/fatih/color"

	"github.com/alibaba/ioc-golang/aop/common"
	"github.com/alibaba/ioc-golang/extension/aop/watch/api/ioc_golang/aop/watch"
)

type watchService struct {
	watch.UnimplementedWatchServiceServer
	watchInterceptor *interceptorImpl
}

func getWatchService() *watchService {
	return &watchService{
		watchInterceptor: getWatchInterceptorSingleton(),
	}
}

func (w *watchService) Watch(req *watch.WatchRequest, svr watch.WatchService_WatchServer) error {
	color.Red("[Debug Server] Receive watch request %+v\n", req.String())
	defer color.Red("[Debug Server] Watch request %s finished \n", req.String())
	sdid := req.GetSdid()
	method := req.GetMethod()
	sendCh := make(chan *watch.WatchResponse)
	var fieldMatcher *common.FieldMatcher
	for _, matcher := range req.GetMatchers() {
		// todo multi match support
		fieldMatcher = &common.FieldMatcher{
			FieldIndex: int(matcher.Index),
			MatchRule:  matcher.GetMatchPath() + "=" + matcher.GetMatchValue(),
		}
	}

	watchCtx := &context{
		SDID:         sdid,
		MethodName:   method,
		Ch:           sendCh,
		FieldMatcher: fieldMatcher,
	}
	w.watchInterceptor.Watch(watchCtx)

	done := svr.Context().Done()
	for {
		select {
		case <-done:
			// watch stop
			w.watchInterceptor.UnWatch(watchCtx)
			return nil
		case watchRsp := <-sendCh:
			if err := svr.Send(watchRsp); err != nil {
				return err
			}
		}
	}
}
