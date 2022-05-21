package config

import (
	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/autowire/normal"
)

func init() {
	autowire.RegisterAutowire(func() autowire.Autowire {
		configAutowire := &Autowire{}
		configAutowire.Autowire = normal.NewNormalAutowire(nil, &paramLoader{}, configAutowire)
		return configAutowire
	}())
}

const Name = "config"

type Autowire struct {
	autowire.Autowire
}

// TagKey re-write NormalAutowire
func (a *Autowire) TagKey() string {
	return Name
}

// GetAllStructDescriptors re-write NormalAutowire
func (a *Autowire) GetAllStructDescriptors() map[string]*autowire.StructDescriptor {
	return configStructDescriptorMap
}

var configStructDescriptorMap = make(map[string]*autowire.StructDescriptor)

func RegisterStructDescriptor(s *autowire.StructDescriptor) {
	s.SetAutowireType(Name)
	configStructDescriptorMap[s.ID()] = s
}

func GetImpl(extensionId string, configPrefix string) (interface{}, error) {
	return autowire.Impl(Name, extensionId, configPrefix)
}
