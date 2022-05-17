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

package rocketmq

import (
	rocketmq "github.com/cinience/go_rocketmq"
	perrors "github.com/pkg/errors"
)

type Config struct {
	rocketmq.Metadata
	AccessProto string
}

func (c *Config) New(impl *Impl) (*Impl, error) {
	ok := false
	if impl.PushConsumer, ok = rocketmq.Consumers[c.AccessProto]; !ok {
		return nil, perrors.Errorf("Invalid AccessProto of rocketmq param")
	}
	if err := impl.PushConsumer.Init(&c.Metadata); err != nil {
		return impl, err
	}
	if impl.Producer, ok = rocketmq.Producers[c.AccessProto]; !ok {
		return nil, perrors.Errorf("Invalid AccessProto of rocketmq param")
	}
	if err := impl.Producer.Init(&c.Metadata); err != nil {
		return impl, err
	}
	return impl, nil
}
