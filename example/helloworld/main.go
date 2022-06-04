package main

import (
	"fmt"
	"time"

	"github.com/alibaba/ioc-golang/example/helloworld/substruct"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/autowire/singleton"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type App struct {
	ServiceImpl1  Service        `singleton:"main.ServiceImpl1"` // inject Service 's ServiceImpl1 implementation
	ServiceImpl2  Service        `singleton:"main.ServiceImpl2"` // inject Service 's ServiceImpl2 implementation
	ServiceStruct *ServiceStruct `singleton:""`                  // inject ServiceStruct struct pointer
}

func (a *App) Run() {
	for {
		time.Sleep(time.Second * 3)
		a.ServiceImpl1.Hello()
		a.ServiceImpl2.Hello()

		fmt.Println(a.ServiceStruct.GetString("laurence", "", nil))
	}
}

type Service interface {
	Hello()
}

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type ServiceImpl1 struct {
}

func (s *ServiceImpl1) Hello() {
	fmt.Println("This is ServiceImpl1, hello world")
}

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type ServiceImpl2 struct {
}

func (s *ServiceImpl2) Hello() {
	fmt.Println("This is ServiceImpl2, hello world")
}

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:isRPCService=true

type ServiceStruct struct {
	//MyName string
}

func (s *ServiceStruct) GetString(name, name2 string, param *substruct.Param) (string, error) {
	return fmt.Sprintf("Hello %s", name), nil
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
