package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/autowire/normal"
	"github.com/alibaba/ioc-golang/autowire/singleton"
	"github.com/alibaba/ioc-golang/extension/normal/redis"
	"github.com/alibaba/ioc-golang/test/docker_compose"
)

func (a *App) TestRun(t *testing.T) {
	redisClientGetByNormalAPI, err := normal.GetImpl("Redis-Impl", &redis.Config{
		Address: "localhost:6379",
		DB:      "0",
	})
	assert.Nil(t, err)
	redisClientGetByNormalAPIImpl := redisClientGetByNormalAPI.(redis.Redis)
	_, err = redisClientGetByNormalAPIImpl.Set("myKey", "myValue", -1)
	assert.Nil(t, err)

	redisClientGetByRedisExtension, err := redis.GetRedis(&redis.Config{
		Address: "localhost:6379",
		DB:      "0",
	})
	assert.Nil(t, err)
	val, err := redisClientGetByRedisExtension.Get("myKey")
	assert.Nil(t, err)
	assert.Equal(t, "myValue", val)
}

func TestGetAPI(t *testing.T) {
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
