/*
 * Copyright (c) 2022, Alibaba Group;
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package grpc

import (
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/alibaba/IOC-Golang/autowire"
	"github.com/alibaba/IOC-Golang/config"
)

type paramLoader struct {
}

/*
Load support load grpc field:
```go
ResourceServiceClient resources.ResourceServiceClient `grpc:"resource-service"`
```go

from:

```yaml
autowire:
  grpc:
    resource-service:
      address: 127.0.0.1:8080
```

Make Dial and generate *grpc.ClientConn as param
*/
func (p *paramLoader) Load(_ *autowire.StructDescriber, fi *autowire.FieldInfo) (interface{}, error) {
	if fi == nil {
		return nil, errors.New("not supported")
	}
	grpcConfig := &Config{}
	if err := config.LoadConfigByPrefix("autowire.grpc."+fi.TagValue, grpcConfig); err != nil {
		return nil, err
	}
	return grpc.Dial(grpcConfig.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

type Config struct {
	Address string
}
