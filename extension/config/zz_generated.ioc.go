//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by iocli, run 'iocli gen' to re-generate

package config

import (
	autowire "github.com/alibaba/ioc-golang/autowire"
	normal "github.com/alibaba/ioc-golang/autowire/normal"
	autowireconfig "github.com/alibaba/ioc-golang/extension/autowire/config"
)

func init() {
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &configFloat64_{}
		},
	})
	configFloat64StructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return new(ConfigFloat64)
		},
		ParamFactory: func() interface{} {
			return new(ConfigFloat64)
		},
		ConstructFunc: func(i interface{}, p interface{}) (interface{}, error) {
			param := p.(configFloat64Interface)
			impl := i.(*ConfigFloat64)
			return param.New(impl)
		},
	}
	autowireconfig.RegisterStructDescriptor(configFloat64StructDescriptor)
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &configInt64_{}
		},
	})
	configInt64StructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return new(ConfigInt64)
		},
		ParamFactory: func() interface{} {
			return new(ConfigInt64)
		},
		ConstructFunc: func(i interface{}, p interface{}) (interface{}, error) {
			param := p.(configInt64Interface)
			impl := i.(*ConfigInt64)
			return param.New(impl)
		},
	}
	autowireconfig.RegisterStructDescriptor(configInt64StructDescriptor)
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &configInt_{}
		},
	})
	configIntStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return new(ConfigInt)
		},
		ParamFactory: func() interface{} {
			return new(ConfigInt)
		},
		ConstructFunc: func(i interface{}, p interface{}) (interface{}, error) {
			param := p.(configIntInterface)
			impl := i.(*ConfigInt)
			return param.New(impl)
		},
	}
	autowireconfig.RegisterStructDescriptor(configIntStructDescriptor)
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &configMap_{}
		},
	})
	configMapStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return new(ConfigMap)
		},
		ParamFactory: func() interface{} {
			return new(ConfigMap)
		},
		ConstructFunc: func(i interface{}, p interface{}) (interface{}, error) {
			param := p.(configMapInterface)
			impl := i.(*ConfigMap)
			return param.New(impl)
		},
	}
	autowireconfig.RegisterStructDescriptor(configMapStructDescriptor)
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &configSlice_{}
		},
	})
	configSliceStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return new(ConfigSlice)
		},
		ParamFactory: func() interface{} {
			return new(ConfigSlice)
		},
		ConstructFunc: func(i interface{}, p interface{}) (interface{}, error) {
			param := p.(configSliceInterface)
			impl := i.(*ConfigSlice)
			return param.New(impl)
		},
	}
	autowireconfig.RegisterStructDescriptor(configSliceStructDescriptor)
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &configString_{}
		},
	})
	configStringStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return new(ConfigString)
		},
		ParamFactory: func() interface{} {
			return new(ConfigString)
		},
		ConstructFunc: func(i interface{}, p interface{}) (interface{}, error) {
			param := p.(configStringInterface)
			impl := i.(*ConfigString)
			return param.New(impl)
		},
	}
	autowireconfig.RegisterStructDescriptor(configStringStructDescriptor)
}

type configFloat64Interface interface {
	New(impl *ConfigFloat64) (*ConfigFloat64, error)
}
type configInt64Interface interface {
	New(impl *ConfigInt64) (*ConfigInt64, error)
}
type configIntInterface interface {
	New(impl *ConfigInt) (*ConfigInt, error)
}
type configMapInterface interface {
	New(impl *ConfigMap) (*ConfigMap, error)
}
type configSliceInterface interface {
	New(impl *ConfigSlice) (*ConfigSlice, error)
}
type configStringInterface interface {
	New(impl *ConfigString) (*ConfigString, error)
}
type configFloat64_ struct {
	Value_ func() float64
	New_   func(impl *ConfigFloat64) (*ConfigFloat64, error)
}

func (c *configFloat64_) Value() float64 {
	return c.Value_()
}

func (c *configFloat64_) New(impl *ConfigFloat64) (*ConfigFloat64, error) {
	return c.New_(impl)
}

type configInt64_ struct {
	Value_ func() int64
	New_   func(impl *ConfigInt64) (*ConfigInt64, error)
}

func (c *configInt64_) Value() int64 {
	return c.Value_()
}

func (c *configInt64_) New(impl *ConfigInt64) (*ConfigInt64, error) {
	return c.New_(impl)
}

type configInt_ struct {
	Value_ func() int
	New_   func(impl *ConfigInt) (*ConfigInt, error)
}

func (c *configInt_) Value() int {
	return c.Value_()
}

func (c *configInt_) New(impl *ConfigInt) (*ConfigInt, error) {
	return c.New_(impl)
}

type configMap_ struct {
	Value_ func() map[string]interface{}
	New_   func(impl *ConfigMap) (*ConfigMap, error)
}

func (c *configMap_) Value() map[string]interface{} {
	return c.Value_()
}

func (c *configMap_) New(impl *ConfigMap) (*ConfigMap, error) {
	return c.New_(impl)
}

type configSlice_ struct {
	Value_ func() []interface{}
	New_   func(impl *ConfigSlice) (*ConfigSlice, error)
}

func (c *configSlice_) Value() []interface{} {
	return c.Value_()
}

func (c *configSlice_) New(impl *ConfigSlice) (*ConfigSlice, error) {
	return c.New_(impl)
}

type configString_ struct {
	Value_ func() string
	New_   func(impl *ConfigString) (*ConfigString, error)
}

func (c *configString_) Value() string {
	return c.Value_()
}

func (c *configString_) New(impl *ConfigString) (*ConfigString, error) {
	return c.New_(impl)
}

type ConfigFloat64IOCInterface interface {
	Value() float64
	New(impl *ConfigFloat64) (*ConfigFloat64, error)
}

type ConfigInt64IOCInterface interface {
	Value() int64
	New(impl *ConfigInt64) (*ConfigInt64, error)
}

type ConfigIntIOCInterface interface {
	Value() int
	New(impl *ConfigInt) (*ConfigInt, error)
}

type ConfigMapIOCInterface interface {
	Value() map[string]interface{}
	New(impl *ConfigMap) (*ConfigMap, error)
}

type ConfigSliceIOCInterface interface {
	Value() []interface{}
	New(impl *ConfigSlice) (*ConfigSlice, error)
}

type ConfigStringIOCInterface interface {
	Value() string
	New(impl *ConfigString) (*ConfigString, error)
}
