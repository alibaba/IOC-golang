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
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
// +ioc:autowire:type=singleton
// +ioc:autowire:paramType=Param
// +ioc:autowire:constructFunc=New

type NamingClient struct {
	naming_client.INamingClient
}

func (n *NamingClient) RegisterInstance(param vo.RegisterInstanceParam) (bool, error) {
	return n.INamingClient.RegisterInstance(param)
}

func (n *NamingClient) DeregisterInstance(param vo.DeregisterInstanceParam) (bool, error) {
	return n.INamingClient.DeregisterInstance(param)
}

func (n *NamingClient) UpdateInstance(param vo.UpdateInstanceParam) (bool, error) {
	return n.INamingClient.UpdateInstance(param)
}

func (n *NamingClient) GetService(param vo.GetServiceParam) (service model.Service, err error) {
	return n.INamingClient.GetService(param)
}

func (n *NamingClient) GetAllServicesInfo(param vo.GetAllServiceInfoParam) (model.ServiceList, error) {
	return n.INamingClient.GetAllServicesInfo(param)
}

func (n *NamingClient) SelectAllInstances(param vo.SelectAllInstancesParam) ([]model.Instance, error) {
	return n.INamingClient.SelectAllInstances(param)
}

func (n *NamingClient) SelectInstances(param vo.SelectInstancesParam) ([]model.Instance, error) {
	return n.INamingClient.SelectInstances(param)
}

func (n *NamingClient) SelectOneHealthyInstance(param vo.SelectOneHealthInstanceParam) (*model.Instance, error) {
	return n.INamingClient.SelectOneHealthyInstance(param)
}

func (n *NamingClient) Subscribe(param *vo.SubscribeParam) error {
	return n.INamingClient.Subscribe(param)
}

func (n *NamingClient) Unsubscribe(param *vo.SubscribeParam) (err error) {
	return n.INamingClient.Unsubscribe(param)
}
