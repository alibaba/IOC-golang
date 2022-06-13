package http_server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"github.com/alibaba/ioc-golang/extension/normal/http_server/ghttp"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
// +ioc:autowire:paramType=HTTPServerConfig
// +ioc:autowire:constructFunc=Create

type Impl struct {
	router       *mux.Router
	mws          []negroni.Handler
	iocGolangMWs []ghttp.Filter
	config       *HTTPServerConfig
}

func (hs *Impl) UseMW(filters ...negroni.Handler) {
	hs.mws = append(hs.mws, filters...)
}

func (hs *Impl) UseIOCGolangMW(filters ...ghttp.Filter) {
	hs.iocGolangMWs = append(hs.iocGolangMWs, filters...)
}

func (hs *Impl) Run(ctx context.Context) {
	s := negroni.Classic()
	for _, handler := range hs.mws {
		s.Use(handler)
	}
	s.UseHandler(hs.router)
	s.Run(":" + hs.config.Port)
}

// RegisterRouterWithRawHttpHandler user API
func (hs *Impl) RegisterRouterWithRawHttpHandler(path string, handler func(w http.ResponseWriter, r *http.Request), method string) {
	hs.router.HandleFunc(path, handler).Methods(method)
}

// RegisterRouter user API
func (hs *Impl) RegisterRouter(path string, handler func(*ghttp.GRegisterController) error, req interface{}, rsp interface{}, method string, filters ...ghttp.Filter) {
	filters = append(hs.iocGolangMWs, filters...)
	ghttp.RegisterRouter(path, hs.router, handler, req, rsp, method, filters)
}

// RegisterWSRouter user API
func (hs *Impl) RegisterWSRouter(path string, handler func(*ghttp.GRegisterWSController)) {
	ghttp.RegisterWSRouter(path, hs.router, handler)
}
