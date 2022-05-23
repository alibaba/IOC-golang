package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/autowire/singleton"
	"github.com/alibaba/ioc-golang/test/docker_compose"
)

func (a *App) TestRun(t *testing.T) {
	_, err := a.NormalRedis.Set("mykey", "db0", -1)
	assert.Nil(t, err)

	_, err = a.NormalDB1Redis.Set("mykey", "db1", -1)
	assert.Nil(t, err)

	_, err = a.NormalDB2Redis.Set("mykey", "db2", -1)
	assert.Nil(t, err)

	_, err = a.NormalDB3Redis.Set("mykey", "db3", -1)
	assert.Nil(t, err)

	val1, err := a.NormalRedis.Get("mykey")
	assert.Nil(t, err)
	assert.Equal(t, "db0", val1)

	val2, err := a.NormalDB1Redis.Get("mykey")
	assert.Nil(t, err)
	assert.Equal(t, "db1", val2)

	val3, err := a.NormalDB2Redis.Get("mykey")
	assert.Nil(t, err)
	assert.Equal(t, "db2", val3)

	val4, err := a.NormalDB3Redis.Get("mykey")
	assert.Nil(t, err)
	assert.Equal(t, "db3", val4)
}

func TestRedisClient(t *testing.T) {
	assert.Nil(t, docker_compose.DockerComposeUp("../docker-compose/docker-compose.yaml", 0))
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
