package protocol_impl

import (
	"context"

	"github.com/alibaba/ioc-golang/autowire/normal"
	"github.com/alibaba/ioc-golang/extension/normal/http_server"
)

var httpServerSingleton http_server.HttpServer

func getHTTPServerSingleton(exportPort string) http_server.HttpServer {
	if httpServerSingleton == nil {
		serverImpl, err := normal.GetImpl(`github.com/alibaba/ioc-golang/extension/normal/http_server.Impl`, &http_server.HTTPServerConfig{
			Port: exportPort,
		})
		if err != nil {
			panic(err)
		}
		httpServerSingleton = serverImpl.(http_server.HttpServer)
		go func() {
			httpServerSingleton.Run(context.Background())
		}()
	}
	return httpServerSingleton
}
