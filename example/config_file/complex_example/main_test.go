package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/config"
	"github.com/alibaba/ioc-golang/test/docker_compose"
)

func (a *App) TestRun(t *testing.T) {
	assert.Equal(t, "myValue", a.ConfigValue.Value())
	assert.Equal(t, "myEnvValue", a.ConfigValueFromEnv.Value())
	assert.Equal(t, "myValue", a.NestedConfigValue.Value())
	assert.Equal(t, "myEnvValue", a.NestedConfigValueFromEnv.Value())
	assert.Nil(t, a.RedisClient.Ping().Err())
}

func TestSetConfigName(t *testing.T) {
	assert.Nil(t, docker_compose.DockerComposeUp("./docker-compose/docker-compose.yaml", 5))
	// start
	os.Setenv("MY_CONFIG_ENV_KEY", "myEnvValue")
	os.Setenv("REDIS_ADDRESS", "127.0.0.1:6379")
	if err := ioc.Load(
		config.WithSearchPath("./redis_config", "./my_custom_config_path"),
		config.WithProfilesActive("pro")); err != nil {
		panic(err)
	}

	app, err := GetAppSingleton()
	if err != nil {
		panic(err)
	}
	app.TestRun(t)
	assert.Nil(t, docker_compose.DockerComposeDown("./docker-compose/docker-compose.yaml"))
}
