//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by iocli, run 'iocli gen' to re-generate

package nacos

import (
	autowire "github.com/alibaba/ioc-golang/autowire"
	normal "github.com/alibaba/ioc-golang/autowire/normal"
	singleton "github.com/alibaba/ioc-golang/autowire/singleton"
	util "github.com/alibaba/ioc-golang/autowire/util"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

func init() {
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &configClient_{}
		},
	})
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &ConfigClient{}
		},
		ParamFactory: func() interface{} {
			var _ paramInterface = &Param{}
			return &Param{}
		},
		ConstructFunc: func(i interface{}, p interface{}) (interface{}, error) {
			param := p.(paramInterface)
			impl := i.(*ConfigClient)
			return param.New(impl)
		},
	})
	singleton.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &ConfigClient{}
		},
		ParamFactory: func() interface{} {
			var _ paramInterface = &Param{}
			return &Param{}
		},
		ConstructFunc: func(i interface{}, p interface{}) (interface{}, error) {
			param := p.(paramInterface)
			impl := i.(*ConfigClient)
			return param.New(impl)
		},
	})
}

type paramInterface interface {
	New(impl *ConfigClient) (*ConfigClient, error)
}
type configClient_ struct {
	GetConfig_          func(param vo.ConfigParam) (string, error)
	PublishConfig_      func(param vo.ConfigParam) (bool, error)
	DeleteConfig_       func(param vo.ConfigParam) (bool, error)
	ListenConfig_       func(params vo.ConfigParam) (err error)
	CancelListenConfig_ func(params vo.ConfigParam) (err error)
	SearchConfig_       func(param vo.SearchConfigParm) (*model.ConfigPage, error)
}

func (c *configClient_) GetConfig(param vo.ConfigParam) (string, error) {
	return c.GetConfig_(param)
}

func (c *configClient_) PublishConfig(param vo.ConfigParam) (bool, error) {
	return c.PublishConfig_(param)
}

func (c *configClient_) DeleteConfig(param vo.ConfigParam) (bool, error) {
	return c.DeleteConfig_(param)
}

func (c *configClient_) ListenConfig(params vo.ConfigParam) (err error) {
	return c.ListenConfig_(params)
}

func (c *configClient_) CancelListenConfig(params vo.ConfigParam) (err error) {
	return c.CancelListenConfig_(params)
}

func (c *configClient_) SearchConfig(param vo.SearchConfigParm) (*model.ConfigPage, error) {
	return c.SearchConfig_(param)
}

type ConfigClientIOCInterface interface {
	GetConfig(param vo.ConfigParam) (string, error)
	PublishConfig(param vo.ConfigParam) (bool, error)
	DeleteConfig(param vo.ConfigParam) (bool, error)
	ListenConfig(params vo.ConfigParam) (err error)
	CancelListenConfig(params vo.ConfigParam) (err error)
	SearchConfig(param vo.SearchConfigParm) (*model.ConfigPage, error)
}

func GetConfigClient(p *Param) (*ConfigClient, error) {
	i, err := normal.GetImpl(util.GetSDIDByStructPtr(new(ConfigClient)), p)
	if err != nil {
		return nil, err
	}
	impl := i.(*ConfigClient)
	return impl, nil
}

func GetConfigClientIOCInterface(p *Param) (ConfigClientIOCInterface, error) {
	i, err := normal.GetImplWithProxy(util.GetSDIDByStructPtr(new(ConfigClient)), p)
	if err != nil {
		return nil, err
	}
	impl := i.(ConfigClientIOCInterface)
	return impl, nil
}

func GetConfigClientSingleton(p *Param) (*ConfigClient, error) {
	i, err := singleton.GetImpl(util.GetSDIDByStructPtr(new(ConfigClient)), p)
	if err != nil {
		return nil, err
	}
	impl := i.(*ConfigClient)
	return impl, nil
}

func GetConfigClientIOCInterfaceSingleton(p *Param) (ConfigClientIOCInterface, error) {
	i, err := singleton.GetImplWithProxy(util.GetSDIDByStructPtr(new(ConfigClient)), p)
	if err != nil {
		return nil, err
	}
	impl := i.(ConfigClientIOCInterface)
	return impl, nil
}
