package main

import (
	"log"
	"runtime"
	"testing"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/autowire/singleton"
	normalNacos "github.com/alibaba/ioc-golang/extension/normal/nacos"
	"github.com/alibaba/ioc-golang/test/docker_compose"
)

const (
	ipFoo   = "127.0.0.1"
	portFoo = 1999
)

func (a *App) TestRun(t *testing.T) {
	testGetAndSetService(t, a.NormalNacosClient, "normal-autowire-client-ioc-golang-debug-service")
	testGetAndSetService(t, a.NormalNacosClient2, "normal-autowire-client2-ioc-golang-debug-service")
	testGetAndSetService(t, a.createByAPINacosClient, "createByAPINacosClient-ioc-golang-debug-service")
}

func testGetAndSetService(t *testing.T, client normalNacos.NacosClient, serviceName string) {
	_, err := client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "127.0.0.1",
		Port:        1999,
		ServiceName: serviceName,
	})
	if err != nil {
		panic(err)
	}

	service, err := client.GetService(vo.GetServiceParam{
		ServiceName: serviceName,
	})
	if err != nil {
		return
	}
	assert.Equal(t, serviceName, service.Name)
	assert.Equal(t, 1, len(service.Hosts))
	assert.Equal(t, ipFoo, service.Hosts[0].Ip)
	assert.Equal(t, uint64(portFoo), service.Hosts[0].Port)
}

func TestNacosClient(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		log.Println("Warning: Nacos image only support amd arch. Skip integration test")
		return
	}
	assert.Nil(t, docker_compose.DockerComposeUp("../docker-compose/docker-compose.yaml", time.Second*10))
	err := ioc.Load()
	retries := 0
	for err != nil && retries < 3 {
		time.Sleep(time.Second * 10)
		err = ioc.Load()
		retries++
	}
	assert.Nil(t, err)
	appInterface, err := singleton.GetImpl("App-App")
	assert.Nil(t, err)
	app := appInterface.(*App)
	app.TestRun(t)

	assert.Nil(t, docker_compose.DockerComposeDown("../docker-compose/docker-compose.yaml"))
}
