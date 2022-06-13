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
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
// +ioc:autowire:paramType=Config
// +ioc:autowire:constructFunc=New

type Impl struct {
	config_client.IConfigClient
	naming_client.INamingClient
}

func (i *Impl) GetConfig(param vo.ConfigParam) (string, error) {
	return i.IConfigClient.GetConfig(param)
}

func (i *Impl) PublishConfig(param vo.ConfigParam) (bool, error) {
	return i.IConfigClient.PublishConfig(param)
}

func (i *Impl) DeleteConfig(param vo.ConfigParam) (bool, error) {
	return i.IConfigClient.DeleteConfig(param)
}

func (i *Impl) ListenConfig(params vo.ConfigParam) (err error) {
	return i.IConfigClient.ListenConfig(params)
}

func (i *Impl) CancelListenConfig(params vo.ConfigParam) (err error) {
	return i.IConfigClient.CancelListenConfig(params)
}

func (i *Impl) SearchConfig(param vo.SearchConfigParm) (*model.ConfigPage, error) {
	return i.IConfigClient.SearchConfig(param)
}

func (i *Impl) RegisterInstance(param vo.RegisterInstanceParam) (bool, error) {
	return i.INamingClient.RegisterInstance(param)
}

func (i *Impl) DeregisterInstance(param vo.DeregisterInstanceParam) (bool, error) {
	return i.INamingClient.DeregisterInstance(param)
}

func (i *Impl) UpdateInstance(param vo.UpdateInstanceParam) (bool, error) {
	return i.INamingClient.UpdateInstance(param)
}

func (i *Impl) GetService(param vo.GetServiceParam) (service model.Service, err error) {
	return i.INamingClient.GetService(param)
}

func (i *Impl) GetAllServicesInfo(param vo.GetAllServiceInfoParam) (model.ServiceList, error) {
	return i.INamingClient.GetAllServicesInfo(param)
}

func (i *Impl) SelectAllInstances(param vo.SelectAllInstancesParam) ([]model.Instance, error) {
	return i.INamingClient.SelectAllInstances(param)
}

func (i *Impl) SelectInstances(param vo.SelectInstancesParam) ([]model.Instance, error) {
	return i.INamingClient.SelectInstances(param)
}

func (i *Impl) SelectOneHealthyInstance(param vo.SelectOneHealthInstanceParam) (*model.Instance, error) {
	return i.INamingClient.SelectOneHealthyInstance(param)
}

func (i *Impl) Subscribe(param *vo.SubscribeParam) error {
	return i.INamingClient.Subscribe(param)
}

func (i *Impl) Unsubscribe(param *vo.SubscribeParam) (err error) {
	return i.INamingClient.Unsubscribe(param)
}
