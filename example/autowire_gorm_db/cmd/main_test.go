package main

import (
	"log"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang/test/docker_compose"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/autowire/singleton"
)

func (a *App) TestRun(t *testing.T) {
	// create table
	assert.Nil(t, a.MyDataTable.GetDB().Model(&MyDataDO{}).AutoMigrate(&MyDataDO{}))
	toInsertMyData := &MyDataDO{
		Value: "first value",
	}
	assert.Nil(t, a.MyDataTable.Insert(toInsertMyData))
	myDataDOs := make([]MyDataDO, 0)
	assert.Nil(t, a.MyDataTable.SelectWhere("id = ?", &myDataDOs, 1))
	assert.Equal(t, 1, len(myDataDOs))
	assert.Equal(t, int32(1), myDataDOs[0].Id)
	assert.Equal(t, "first value", myDataDOs[0].Value)
}

func TestGORM(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		log.Println("Warning: Mysql image only support amd arch. Skip integration test")
		return
	}
	assert.Nil(t, docker_compose.DockerComposeUp("../docker-compose/docker-compose.yaml", time.Second*10))
	if err := ioc.Load(); err != nil {
		panic(err)
	}
	appInterface, err := singleton.GetImpl("App-App")
	if err != nil {
		panic(err)
	}
	app := appInterface.(*App)
	app.TestRun(t)
	assert.Nil(t, docker_compose.DockerComposeDown("../docker-compose/docker-compose.yaml"))
}
