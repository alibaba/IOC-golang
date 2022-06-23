/*
 * Copyright (c) 2022, Alibaba Group;
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,马克
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
// +ioc:autowire:type=singleton
// +ioc:autowire:paramType=Param
// +ioc:autowire:constructFunc=New

type ConfigClient struct {
	config_client.IConfigClient
}

func (c *ConfigClient) GetConfig(param vo.ConfigParam) (string, error) {
	return c.IConfigClient.GetConfig(param)
}

func (c *ConfigClient) PublishConfig(param vo.ConfigParam) (bool, error) {
	return c.IConfigClient.PublishConfig(param)
}

func (c *ConfigClient) DeleteConfig(param vo.ConfigParam) (bool, error) {
	return c.IConfigClient.DeleteConfig(param)
}

func (c *ConfigClient) ListenConfig(params vo.ConfigParam) (err error) {
	return c.IConfigClient.ListenConfig(params)
}

func (c *ConfigClient) CancelListenConfig(params vo.ConfigParam) (err error) {
	return c.IConfigClient.CancelListenConfig(params)
}

func (c *ConfigClient) SearchConfig(param vo.SearchConfigParm) (*model.ConfigPage, error) {
	return c.IConfigClient.SearchConfig(param)
}
