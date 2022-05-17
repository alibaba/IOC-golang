package config

import (
	"github.com/alibaba/IOC-Golang/autowire"
	"github.com/alibaba/IOC-Golang/autowire/normal"
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

// GetAllStructDescribers re-write NormalAutowire
func (a *Autowire) GetAllStructDescribers() map[string]*autowire.StructDescriber {
	return configStructDescriberMap
}

var configStructDescriberMap = make(map[string]*autowire.StructDescriber)

func RegisterStructDescriber(s *autowire.StructDescriber) {
	s.SetAutowireType(Name)
	configStructDescriberMap[s.ID()] = s
}

func GetImpl(extensionId string, configPrefix string) (interface{}, error) {
	return autowire.Impl(Name, extensionId, configPrefix)
}
