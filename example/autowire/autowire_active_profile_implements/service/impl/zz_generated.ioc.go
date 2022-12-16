//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by iocli, run 'iocli gen' to re-generate

package impl

import (
	autowire "github.com/alibaba/ioc-golang/autowire"
	normal "github.com/alibaba/ioc-golang/autowire/normal"
	singleton "github.com/alibaba/ioc-golang/autowire/singleton"
	util "github.com/alibaba/ioc-golang/autowire/util"
	service "github.com/alibaba/ioc-golang/example/autowire/autowire_active_profile_implements/service"
)

func init() {
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &serviceDefaultImpl_{}
		},
	})
	serviceDefaultImplStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &serviceDefaultImpl{}
		},
		Metadata: map[string]interface{}{
			"aop": map[string]interface{}{},
			"autowire": map[string]interface{}{
				"common": map[string]interface{}{
					"implements": []interface{}{
						new(service.Service),
					},
				},
			},
		},
	}
	singleton.RegisterStructDescriptor(serviceDefaultImplStructDescriptor)
	var _ service.Service = &serviceDefaultImpl{}
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &serviceDevImpl_{}
		},
	})
	serviceDevImplStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &serviceDevImpl{}
		},
		Metadata: map[string]interface{}{
			"aop": map[string]interface{}{},
			"autowire": map[string]interface{}{
				"common": map[string]interface{}{
					"implements": []interface{}{
						new(service.Service),
					},
					"activeProfile": "dev",
				},
			},
		},
	}
	singleton.RegisterStructDescriptor(serviceDevImplStructDescriptor)
	var _ service.Service = &serviceDevImpl{}
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &serviceProImpl_{}
		},
	})
	serviceProImplStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &serviceProImpl{}
		},
		Metadata: map[string]interface{}{
			"aop": map[string]interface{}{},
			"autowire": map[string]interface{}{
				"common": map[string]interface{}{
					"implements": []interface{}{
						new(service.Service),
					},
					"activeProfile": "pro",
				},
			},
		},
	}
	singleton.RegisterStructDescriptor(serviceProImplStructDescriptor)
	var _ service.Service = &serviceProImpl{}
}

type serviceDefaultImpl_ struct {
	GetHelloString_ func(name string) string
}

func (s *serviceDefaultImpl_) GetHelloString(name string) string {
	return s.GetHelloString_(name)
}

type serviceDevImpl_ struct {
	GetHelloString_ func(name string) string
}

func (s *serviceDevImpl_) GetHelloString(name string) string {
	return s.GetHelloString_(name)
}

type serviceProImpl_ struct {
	GetHelloString_ func(name string) string
}

func (s *serviceProImpl_) GetHelloString(name string) string {
	return s.GetHelloString_(name)
}

type serviceDefaultImplIOCInterface interface {
	GetHelloString(name string) string
}

type serviceDevImplIOCInterface interface {
	GetHelloString(name string) string
}

type serviceProImplIOCInterface interface {
	GetHelloString(name string) string
}

var _serviceDefaultImplSDID string

func GetserviceDefaultImplSingleton() (*serviceDefaultImpl, error) {
	if _serviceDefaultImplSDID == "" {
		_serviceDefaultImplSDID = util.GetSDIDByStructPtr(new(serviceDefaultImpl))
	}
	i, err := singleton.GetImpl(_serviceDefaultImplSDID, nil)
	if err != nil {
		return nil, err
	}
	impl := i.(*serviceDefaultImpl)
	return impl, nil
}

func GetserviceDefaultImplIOCInterfaceSingleton() (serviceDefaultImplIOCInterface, error) {
	if _serviceDefaultImplSDID == "" {
		_serviceDefaultImplSDID = util.GetSDIDByStructPtr(new(serviceDefaultImpl))
	}
	i, err := singleton.GetImplWithProxy(_serviceDefaultImplSDID, nil)
	if err != nil {
		return nil, err
	}
	impl := i.(serviceDefaultImplIOCInterface)
	return impl, nil
}

type ThisserviceDefaultImpl struct {
}

func (t *ThisserviceDefaultImpl) This() serviceDefaultImplIOCInterface {
	thisPtr, _ := GetserviceDefaultImplIOCInterfaceSingleton()
	return thisPtr
}

var _serviceDevImplSDID string

func GetserviceDevImplSingleton() (*serviceDevImpl, error) {
	if _serviceDevImplSDID == "" {
		_serviceDevImplSDID = util.GetSDIDByStructPtr(new(serviceDevImpl))
	}
	i, err := singleton.GetImpl(_serviceDevImplSDID, nil)
	if err != nil {
		return nil, err
	}
	impl := i.(*serviceDevImpl)
	return impl, nil
}

func GetserviceDevImplIOCInterfaceSingleton() (serviceDevImplIOCInterface, error) {
	if _serviceDevImplSDID == "" {
		_serviceDevImplSDID = util.GetSDIDByStructPtr(new(serviceDevImpl))
	}
	i, err := singleton.GetImplWithProxy(_serviceDevImplSDID, nil)
	if err != nil {
		return nil, err
	}
	impl := i.(serviceDevImplIOCInterface)
	return impl, nil
}

type ThisserviceDevImpl struct {
}

func (t *ThisserviceDevImpl) This() serviceDevImplIOCInterface {
	thisPtr, _ := GetserviceDevImplIOCInterfaceSingleton()
	return thisPtr
}

var _serviceProImplSDID string

func GetserviceProImplSingleton() (*serviceProImpl, error) {
	if _serviceProImplSDID == "" {
		_serviceProImplSDID = util.GetSDIDByStructPtr(new(serviceProImpl))
	}
	i, err := singleton.GetImpl(_serviceProImplSDID, nil)
	if err != nil {
		return nil, err
	}
	impl := i.(*serviceProImpl)
	return impl, nil
}

func GetserviceProImplIOCInterfaceSingleton() (serviceProImplIOCInterface, error) {
	if _serviceProImplSDID == "" {
		_serviceProImplSDID = util.GetSDIDByStructPtr(new(serviceProImpl))
	}
	i, err := singleton.GetImplWithProxy(_serviceProImplSDID, nil)
	if err != nil {
		return nil, err
	}
	impl := i.(serviceProImplIOCInterface)
	return impl, nil
}

type ThisserviceProImpl struct {
}

func (t *ThisserviceProImpl) This() serviceProImplIOCInterface {
	thisPtr, _ := GetserviceProImplIOCInterfaceSingleton()
	return thisPtr
}
