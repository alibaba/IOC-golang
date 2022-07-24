package allimpls

import (
	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/extension/autowire/common"
)

func parseAllImpledIntefacesFromSDMetadata(metadata autowire.Metadata) []interface{} {
	autowireMetadata := common.ParseAutowireMetadataFromSDMetadata(metadata)
	if autowireMetadata == nil {
		return nil
	}
	result, ok := autowireMetadata[Name].([]interface{})
	if !ok {
		return nil
	}
	return result
}
