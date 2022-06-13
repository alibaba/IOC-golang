//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by iocli

package main

import (
	"github.com/alibaba/ioc-golang/autowire"
	normal "github.com/alibaba/ioc-golang/autowire/normal"
	singleton "github.com/alibaba/ioc-golang/autowire/singleton"
	util "github.com/alibaba/ioc-golang/autowire/util"
)

func init() {
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &app_{}
		},
	})
	singleton.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &App{}
		},
	})
}

type app_ struct {
	Run_ func()
}

func (a *app_) Run() {
	a.Run_()
}

type AppIOCInterface interface {
	Run()
}

func GetApp() (*App, error) {
	i, err := singleton.GetImpl(util.GetSDIDByStructPtr(new(App)), nil)
	if err != nil {
		return nil, err
	}
	impl := i.(*App)
	return impl, nil
}
func GetAppIOCInterface() (AppIOCInterface, error) {
	return GetApp()
}
