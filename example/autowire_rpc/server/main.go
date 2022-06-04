package main

import (
	"github.com/alibaba/ioc-golang"
	_ "github.com/alibaba/ioc-golang/example/autowire_rpc/server/pkg/service"
)

func main() {
	// start
	if err := ioc.Load(); err != nil {
		panic(err)
	}
	select {}
}
