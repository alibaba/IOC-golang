package triple

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang/autowire"
)

type mockImpl struct {
}

const mockImplFullName = "github.com/alibaba/ioc-golang/extension/autowire/triple.mockImpl"

func TestAutowire_RegisterAndGetAllStructDescriptors(t *testing.T) {
	t.Run("test triple autowire register and get all struct descriptors", func(t *testing.T) {
		sd := &autowire.StructDescriptor{
			Factory: func() interface{} {
				return &mockImpl{}
			},
		}
		RegisterStructDescriptor(sd)
		a := &Autowire{}
		allStructDesc := a.GetAllStructDescriptors()
		assert.NotNil(t, allStructDesc)
		structDesc, ok := allStructDesc[mockImplFullName]
		assert.True(t, ok)
		assert.Equal(t, mockImplFullName, structDesc.ID())
	})
}

func TestAutowire_TagKey(t *testing.T) {
	t.Run("test triple autowire tag", func(t *testing.T) {
		a := &Autowire{}
		assert.Equal(t, Name, a.TagKey())
	})
}
