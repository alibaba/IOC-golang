package triple

import (
	"errors"

	dubboConfig "dubbo.apache.org/dubbo-go/v3/config"

	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/config"
)

type paramLoader struct {
}

/*
Load support load grpc field:
```go
ResourceServiceClient resources.ResourceServiceClient `triple:"resource-service"`
```go

from:

```yaml
autowire:
  triple:
    resource-service:
      group:"myGroup"
      version:"myVersion"
```

Make Dial and generate *grpc.ClientConn as param
*/
func (p *paramLoader) Load(_ *autowire.StructDescriptor, fi *autowire.FieldInfo) (interface{}, error) {
	if fi == nil {
		return nil, errors.New("not supported")
	}
	tripleConfig := &Config{}
	if err := config.LoadConfigByPrefix("autowire.triple."+fi.TagValue, tripleConfig); err != nil {
		return nil, err
	}
	return tripleConfig, nil
}

type Config struct {
	dubboConfig.ReferenceConfig
}
