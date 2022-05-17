package grpc

import (
	"github.com/alibaba/IOC-Golang/autowire"
	"github.com/alibaba/IOC-Golang/autowire/util"
)

type sdIDParser struct {
}

/*
Parse support parse field like:
ResourceServiceClient resources.ResourceServiceClient `grpc:"resource-service"`
to struct describer ID 'ResourceServiceClient-ResourceServiceClient'
*/
func (p *sdIDParser) Parse(fi *autowire.FieldInfo) (string, error) {
	grpcInterfaceName := fi.FieldType
	return util.GetIdByNamePair(grpcInterfaceName, grpcInterfaceName), nil
}
