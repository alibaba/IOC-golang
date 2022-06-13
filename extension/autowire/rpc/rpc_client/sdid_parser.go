package rpc_client

import (
	"strings"

	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/autowire/util"
)

type sdidParser struct {
}

func (s *sdidParser) Parse(fi *autowire.FieldInfo) (string, error) {
	injectStructName := fi.FieldType
	splitedTagValue := strings.Split(fi.TagValue, ",")
	if len(splitedTagValue) > 0 && splitedTagValue[0] != "" {
		injectStructName = splitedTagValue[0]
	} else {
		// no struct sdid in tag
		if !util.IsPointerField(fi.FieldReflectType) && strings.HasSuffix(injectStructName, "IOCRPCClient") {
			// is interface field without valid sdid from tag value, and with 'IOCInterface' suffix
			// load trim suffix as sdid
			var err error
			injectStructName, err = util.ToRPCClientStubSDID(injectStructName)
			if err != nil {
				return "", err
			}
		}
	}
	return autowire.GetSDIDByAliasIfNecessary(injectStructName), nil
}
