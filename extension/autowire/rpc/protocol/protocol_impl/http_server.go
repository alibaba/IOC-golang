package protocol_impl

import (
	"context"

	"github.com/alibaba/ioc-golang/extension/normal/http_server"
)

var httpServerSingleton http_server.ImplIOCInterface

func getHTTPServerSingleton(exportPort string) http_server.ImplIOCInterface {
	if httpServerSingleton == nil {
		serverImpl, err := http_server.GetImplIOCInterface(&http_server.HTTPServerConfig{
			Port: exportPort,
		})
		if err != nil {
			panic(err)
		}
		httpServerSingleton = serverImpl
		go func() {
			httpServerSingleton.Run(context.Background())
		}()
	}
	return httpServerSingleton
}
