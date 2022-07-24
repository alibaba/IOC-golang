package transaction

import (
	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/extension/aop/common"
)

func parseRollbackMethodNameFromSDMetadata(metadata autowire.Metadata, methodName string) (string, bool) {
	if aopMetadata := common.ParseAOPMetadataFromSDMetadata(metadata); aopMetadata != nil {
		if txAOPMetadataVal, ok := aopMetadata[Name]; ok {
			if txAOPMetadata, ok2 := txAOPMetadataVal.(map[string]string); ok2 {
				rollbackMethodName, found := txAOPMetadata[methodName]
				return rollbackMethodName, found
			}
		}
	}
	return "", false
}
