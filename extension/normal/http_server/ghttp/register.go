package ghttp

import (
	"errors"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var rspImpPackage RspPackageFactory

func init() {
	RegisterRspPackage(NewDefaultRspPackage)
}

func writeDefaultHeader(rsp http.ResponseWriter, req *http.Request) {
	rsp.Header().Add("X-Content-Type-Options", "nosniff")
	ct := rsp.Header().Get("Content-Type")
	if ct == "" {
		ct = req.Header.Get("Content-Type")
		if req.Method == "GET" || ct == "" {
			ct = "application/json"
		}
		rsp.Header().Add("Content-Type", ct)
	}
}

func getIOCGolangHttpHandler(handler func(*GRegisterController) error, req, rsp interface{}, filters []Filter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		retPkg := rspImpPackage()

		// recovery
		defer func() {
			if e := recover(); e != nil {
				buf := make([]byte, 1024)
				buf = buf[:runtime.Stack(buf, false)]
				log.Panicf("%s\n%s\n", e, buf)
				retPkg.SetErrorPkg(w, errors.New("server panic"), DefaultHttpErrorCode)
			}
		}()

		writeDefaultHeader(w, r)

		chain := Chain{}
		chain.AddFilter(filters) // 注册过滤器

		tRegisterController := GRegisterController{
			R:           r,
			W:           w,
			RspCode:     UnsetHttpCode,
			IfNeedWrite: true,
			VarsMap:     mux.Vars(r),
		}

		if req != nil {
			requestType := reflect.TypeOf(req).Elem()
			tRegisterController.Req = reflect.New(requestType).Interface()
			if err := tRegisterController.GetReqData(r); err != nil {
				retPkg.SetErrorPkg(w, err, tRegisterController.RspCode) // go
				return
			}
		}

		if rsp != nil {
			rspType := reflect.TypeOf(rsp).Elem()
			tRegisterController.Rsp = reflect.New(rspType).Interface()
		}

		if err := chain.Handle(&tRegisterController, handler); err != nil {
			retPkg.SetErrorPkg(w, err, tRegisterController.RspCode)
			return
		}

		if !tRegisterController.IfNeedWrite {
			return
		}

		retPkg.SetSuccessPkg(w, tRegisterController.Rsp, tRegisterController.RspCode)
	}
}

func getIOCGolangWSHandler(handler func(*GRegisterWSController)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		retPkg := rspImpPackage()

		// recovery
		defer func() {
			if e := recover(); e != nil {
				buf := make([]byte, 1024)
				buf = buf[:runtime.Stack(buf, false)]
				log.Panicf("%s\n%s\n", e, buf)
				retPkg.SetErrorPkg(w, errors.New("server panic"), DefaultHttpErrorCode)
			}
		}()
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		tRegisterController := GRegisterWSController{
			WSConn: conn,
			R:      r,
		}
		handler(&tRegisterController)
	}
}

// 自定义回包函数
func RegisterRspPackage(rspUserImplPackageFactory RspPackageFactory) {
	rspImpPackage = rspUserImplPackageFactory
}

func checkMethod(method string) (string, bool) {
	if method == "GET" || method == "POST" || method == "DELETE" ||
		method == "PATCH" || method == "PUT" {
		return method, true
	}
	if method == "get" || method == "post" || method == "delete" ||
		method == "patch" || method == "put" {
		return strings.ToUpper(method), true
	}
	return "", false
}

func RegisterRouter(path string, r *mux.Router, handler func(*GRegisterController) error, req, rsp interface{}, method string, filters []Filter) {
	iocGolangHttpHandler := getIOCGolangHttpHandler(handler, req, rsp, filters)
	afterCheckedMethod, ok := checkMethod(method)
	if !ok {
		log.Panic("RegisterRouter: method unsupported")
		return
	}
	r.HandleFunc(path, iocGolangHttpHandler).Methods(afterCheckedMethod)
}

func RegisterWSRouter(path string, r *mux.Router, handler func(*GRegisterWSController)) {
	trpcWSHandler := getIOCGolangWSHandler(handler)
	r.HandleFunc(path, trpcWSHandler)
}

func NewHttpRegister() {

}
