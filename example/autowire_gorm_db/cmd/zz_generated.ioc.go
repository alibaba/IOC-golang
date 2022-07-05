//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by iocli, run 'iocli gen' to re-generate

package main

import (
	autowire "github.com/alibaba/ioc-golang/autowire"
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
	appStructDescriptor := &autowire.StructDescriptor{
		Alias: "AppAlias",
		Factory: func() interface{} {
			return &App{}
		},
	}
	singleton.RegisterStructDescriptor(appStructDescriptor)
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

func GetAppSingleton() (*App, error) {
	i, err := singleton.GetImpl(util.GetSDIDByStructPtr(new(App)), nil)
	if err != nil {
		return nil, err
	}
	impl := i.(*App)
	return impl, nil
}

func GetAppIOCInterfaceSingleton() (AppIOCInterface, error) {
	i, err := singleton.GetImplWithProxy(util.GetSDIDByStructPtr(new(App)), nil)
	if err != nil {
		return nil, err
	}
	impl := i.(AppIOCInterface)
	return impl, nil
}
