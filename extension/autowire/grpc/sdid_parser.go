package grpc

import (
	"github.com/alibaba/ioc-golang/autowire"
)

type sdIDParser struct {
}

/*
Parse support parse field like:
ResourceServiceClient resources.ResourceServiceClient `grpc:"resource-service"`
to struct descriptor ID 'ResourceServiceClient-ResourceServiceClient'
*/
func (p *sdIDParser) Parse(fi *autowire.FieldInfo) (string, error) {
	return fi.FieldType, nil
}
