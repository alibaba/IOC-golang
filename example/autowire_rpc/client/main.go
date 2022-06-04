package main

import (
	"fmt"
	"time"

	"github.com/alibaba/ioc-golang/example/autowire_rpc/server/pkg/dto"
	"github.com/alibaba/ioc-golang/example/autowire_rpc/server/pkg/service/api"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/autowire/singleton"
	_ "github.com/alibaba/ioc-golang/extension/autowire/rpc/rpc_client"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type App struct {
	ServiceStruct *api.ServiceStructIOCRPCClient `rpc-client:",address=localhost:2022"`
}

func (a *App) Run() {
	for {
		time.Sleep(time.Second * 3)
		age := 23
		age32 := int32(23)
		age64 := int64(23)
		ageF32 := float32(23)
		ageF64 := float64(23)
		usr, paramUser, err := a.ServiceStruct.GetUser("laurence", 23, 23, 23, 23, 23, &age, &age32, &age64, &ageF32, &ageF64, &dto.User{
			Name: "laurence",
			Age:  18,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(usr, usr.Age)
		fmt.Println(paramUser, paramUser.Age)
	}
}

func main() {
	// start
	if err := ioc.Load(); err != nil {
		panic(err)
	}

	// 'App' is alias name
	// We can get instance by ths id
	appInterface, err := singleton.GetImpl("main.App")
	if err != nil {
		panic(err)
	}
	app := appInterface.(*App)
	app.Run()
}
