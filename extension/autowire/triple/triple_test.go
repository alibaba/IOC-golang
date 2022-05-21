package triple

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang/autowire"
)

type mockInterface interface {
}

type mockImpl struct {
}

const mockInterfaceName = "mockInterface"
const mockImplName = "mockImpl"

func TestAutowire_RegisterAndGetAllStructDescriptors(t *testing.T) {
	t.Run("test config autowire register and get all struct descriptors", func(t *testing.T) {
		sd := &autowire.StructDescriptor{
			Interface: new(mockInterface),
			Factory: func() interface{} {
				return &mockImpl{}
			},
		}
		RegisterStructDescriptor(sd)
		a := &Autowire{}
		allStructDesc := a.GetAllStructDescriptors()
		assert.NotNil(t, allStructDesc)
		structDesc, ok := allStructDesc[mockInterfaceName+"-"+mockImplName]
		assert.True(t, ok)
		assert.Equal(t, mockInterfaceName+"-"+mockImplName, structDesc.ID())
	})
}

func TestAutowire_TagKey(t *testing.T) {
	t.Run("test grpc autowire tag", func(t *testing.T) {
		a := &Autowire{}
		assert.Equal(t, Name, a.TagKey())
	})
}
