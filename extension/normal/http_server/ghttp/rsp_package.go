package ghttp

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RspPackage interface {
	SetSuccessPkg(w http.ResponseWriter, msg interface{}, retCode httpCode)
	SetErrorPkg(w http.ResponseWriter, err error, retCode httpCode)
}

type RspPackageFactory func() RspPackage

type DefaultRspPackage struct {
}

func NewDefaultRspPackage() RspPackage {
	return &DefaultRspPackage{}
}

func (rpkg *DefaultRspPackage) SetSuccessPkg(w http.ResponseWriter, result interface{}, retCode httpCode) {
	if retCode == UnsetHttpCode {
		retCode = DefaultHttpSuccessCode
	}
	w.WriteHeader(int(retCode))
	rspPkgBody, err := json.Marshal(result)
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf(`{"retcode":%d, "retmsg":"fatel err: marshal rspPackage failed", "result": null}`, -1)))
		return
	}
	_, _ = w.Write(rspPkgBody)
}

func (rpkg *DefaultRspPackage) SetErrorPkg(w http.ResponseWriter, err error, retCode httpCode) {
	if retCode == UnsetHttpCode {
		retCode = http.StatusInternalServerError
	}
	w.WriteHeader(int(retCode))
	rspPackageBody, err := json.Marshal(err.Error())
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf(`{"retcode":%d, "retmsg":"fatel err: marshal rspPackage failed", "result": null}`, -1)))
		return
	}
	_, _ = w.Write(rspPackageBody)
}

// FomattedRspPackage
// 框架提供的格式化回包，包含三个字段,与下面的 DefaultRspPackage 选择使用
type FomattedRspPackage struct {
	Retcode int32       `json:"retcode"`
	Retmsg  string      `json:"retmsg"`
	Result  interface{} `json:"result"`
}

func NewFomattedRspPackage() RspPackage {
	return &FomattedRspPackage{}
}

func (rpkg *FomattedRspPackage) SetSuccessPkg(w http.ResponseWriter, result interface{}, retCode httpCode) {
	rpkg.Retmsg = "ok"
	rpkg.Retcode = 0
	rpkg.Result = result
	if retCode == UnsetHttpCode {
		retCode = DefaultHttpSuccessCode
	}
	w.WriteHeader(int(retCode))
	rspPackageBody, err := json.Marshal(*rpkg)
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf(`{"retcode":%d, "retmsg":"fatel err: marshal rspPackage failed", "result": null}`, -1)))
		return
	}
	_, _ = w.Write(rspPackageBody)
}

func (rpkg *FomattedRspPackage) SetErrorPkg(w http.ResponseWriter, err error, retCode httpCode) {
	rpkg.Retmsg = err.Error()
	rpkg.Retcode = -1
	rpkg.Result = nil
	if retCode == UnsetHttpCode {
		retCode = DefaultHttpSuccessCode
	}
	w.WriteHeader(int(retCode))
	rpkgBody, err := json.Marshal(*rpkg)
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf(`{"retcode":%d, "retmsg":"fatel err: marshal rspPackage failed", "result": null}`, -1)))
		return
	}
	_, _ = w.Write(rpkgBody)
}

//

// ResultAndOKRspPackage contains ok and result
type ResultAndOKRspPackage struct {
	Result interface{} `json:"result"`
	OK     bool        `json:"ok"`
}

func NewResultAndOKRspPackage() RspPackage {
	return &ResultAndOKRspPackage{}
}

func (rpkg *ResultAndOKRspPackage) SetSuccessPkg(w http.ResponseWriter, result interface{}, retCode httpCode) {
	rpkg.Result = result
	rpkg.OK = true
	if retCode == UnsetHttpCode {
		retCode = DefaultHttpSuccessCode
	}
	w.WriteHeader(int(retCode))
	rpkgBody, err := json.Marshal(*rpkg)
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf(`{"retcode":%d, "retmsg":"fatel err: marshal rspPackage failed", "result": null}`, -1)))
		return
	}
	_, _ = w.Write(rpkgBody)
}

func (rpkg *ResultAndOKRspPackage) SetErrorPkg(w http.ResponseWriter, err error, retCode httpCode) {
	rpkg.Result = err.Error()
	rpkg.OK = false
	if retCode == UnsetHttpCode {
		retCode = DefaultHttpErrorCode
	}
	w.WriteHeader(int(retCode))
	rpkgBody, err := json.Marshal(*rpkg)
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf(`{"retcode":%d, "retmsg":"fatel err: marshal rspPackage failed", "result": null}`, -1)))
		return
	}
	_, _ = w.Write(rpkgBody)
}
