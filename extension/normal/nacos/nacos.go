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

package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"github.com/alibaba/IOC-Golang/autowire/normal"
)

const SDID = "NacosClient-Impl"

type NacosClient interface {
	GetConfig(param vo.ConfigParam) (string, error)
	PublishConfig(param vo.ConfigParam) (bool, error)
	DeleteConfig(param vo.ConfigParam) (bool, error)
	ListenConfig(params vo.ConfigParam) (err error)
	CancelListenConfig(params vo.ConfigParam) (err error)
	SearchConfig(param vo.SearchConfigParm) (*model.ConfigPage, error)

	RegisterInstance(param vo.RegisterInstanceParam) (bool, error)
	DeregisterInstance(param vo.DeregisterInstanceParam) (bool, error)
	UpdateInstance(param vo.UpdateInstanceParam) (bool, error)
	GetService(param vo.GetServiceParam) (service model.Service, err error)
	GetAllServicesInfo(param vo.GetAllServiceInfoParam) (model.ServiceList, error)
	SelectAllInstances(param vo.SelectAllInstancesParam) ([]model.Instance, error)
	SelectInstances(param vo.SelectInstancesParam) ([]model.Instance, error)
	SelectOneHealthyInstance(param vo.SelectOneHealthInstanceParam) (*model.Instance, error)
	Subscribe(param *vo.SubscribeParam) error
	Unsubscribe(param *vo.SubscribeParam) (err error)
}

// +ioc:autowire=true
// +ioc:autowire:interface=NacosClient
// +ioc:autowire:type=normal
// +ioc:autowire:paramType=Config
// +ioc:autowire:constructFunc=New

type Impl struct {
	config_client.IConfigClient
	naming_client.INamingClient
}

var _ NacosClient = &Impl{}

func GetNacosClient(config *Config) (NacosClient, error) {
	nacosClientImpl, err := normal.GetImpl(SDID, config)
	if err != nil {
		return nil, err
	}
	return nacosClientImpl.(NacosClient), nil
}
