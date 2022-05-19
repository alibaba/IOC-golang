package triple

import (
	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/autowire/singleton"
)

func init() {
	autowire.RegisterAutowire(func() autowire.Autowire {
		grpcAutowire := &Autowire{}
		grpcAutowire.Autowire = singleton.NewSingletonAutowire(&sdIDParser{}, &paramLoader{}, grpcAutowire)
		return grpcAutowire
	}())
}

const Name = "triple"

type Autowire struct {
	autowire.Autowire
}

// TagKey re-write SingletonAutowire
func (a *Autowire) TagKey() string {
	return Name
}

func (a *Autowire) CanBeEntrance() bool {
	return false
}

// GetAllStructDescriptors re-write SingletonAutowire
func (a *Autowire) GetAllStructDescriptors() map[string]*autowire.StructDescriptor {
	return tripleStructDescriptorMap
}

var tripleStructDescriptorMap = make(map[string]*autowire.StructDescriptor)

func RegisterStructDescriptor(s *autowire.StructDescriptor) {
	s.SetAutowireType(Name)
	tripleStructDescriptorMap[s.ID()] = s
}

func GetImpl(extensionId string) (interface{}, error) {
	return autowire.Impl(Name, extensionId, nil)
}
