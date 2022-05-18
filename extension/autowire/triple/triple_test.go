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

func TestAutowire_RegisterAndGetAllStructDescribers(t *testing.T) {
	t.Run("test config autowire register and get all struct describers", func(t *testing.T) {
		sd := &autowire.StructDescriber{
			Interface: new(mockInterface),
			Factory: func() interface{} {
				return &mockImpl{}
			},
		}
		RegisterStructDescriber(sd)
		a := &Autowire{}
		allStructDesc := a.GetAllStructDescribers()
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
