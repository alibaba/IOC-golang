package allimpls

import (
	"fmt"
	"reflect"

	"github.com/alibaba/ioc-golang/autowire"
)

type sdIDParser struct {
	sdidFieldInterfaceReflectTypeMap map[string]reflect.Type
}

/*
Parse support parse field like:
ResourceServiceClient resources.ResourceServiceClient `grpc:"resource-service"`
to struct descriptor ID 'ResourceServiceClient-ResourceServiceClient'
*/
func (p *sdIDParser) Parse(fi *autowire.FieldInfo) (string, error) {
	if fi.FieldReflectType.Kind() != reflect.Slice {
		return "", fmt.Errorf("[Autowire allimpls] invalid field %s, field should be interface slice", fi)
	}
	if fi.FieldReflectType.Elem().Kind() != reflect.Interface {
		return "", fmt.Errorf("[Autowire allimpls] invalid field %s, field should be interface slice", fi)
	}
	return fi.FieldType, nil
}

func (p *sdIDParser) newFieldInterfaceSliceValue(sdid string) interface{} {
	fieldInterfaceReflectType, ok := p.sdidFieldInterfaceReflectTypeMap[sdid]
	if !ok {
		// never occurs
		return nil
	}

	return reflect.MakeSlice(reflect.SliceOf(fieldInterfaceReflectType), 0, 0)
}

var sdidParserSingleton *sdIDParser

func getSDIDParserSingleton() *sdIDParser {
	if sdidParserSingleton == nil {
		sdidParserSingleton = &sdIDParser{
			sdidFieldInterfaceReflectTypeMap: impledInterfaceSDIDTypeMap,
		}
	}
	return sdidParserSingleton
}
