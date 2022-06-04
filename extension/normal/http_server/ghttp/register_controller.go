package ghttp

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/gorilla/websocket"
)

var v *validator.Validate
var defaultSchemaDecoder *schema.Decoder

func init() {
	defaultSchemaDecoder = schema.NewDecoder()
	defaultSchemaDecoder.IgnoreUnknownKeys(true)
	v = validator.New()
}

type GRegisterController struct {
	Req         interface{}
	Rsp         interface{}
	R           *http.Request
	W           http.ResponseWriter
	VarsMap     map[string]string
	RspCode     httpCode
	IfNeedWrite bool
}

type GRegisterWSController struct {
	WSConn *websocket.Conn
	R      *http.Request
}

func (trc *GRegisterController) GetReqData(r *http.Request) error {
	var err error

	if err = r.ParseForm(); err != nil {
		return err
	}
	if r.Header == nil {
		return errors.New("r.Header nil ptr error")
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("read body err:", err)
		return err
	}
	if len(data) != 0 {
		if err := json.Unmarshal(data, trc.Req); err != nil {
			log.Println("request unmarshal err:", err)
			return err
		}
	}

	if err := defaultSchemaDecoder.Decode(trc.Req, r.Form); err != nil {
		log.Println("decode r.form err:", err)
		return err
	}

	if err = v.Struct(trc.Req); err != nil {
		log.Println("validator check failed:", err)
		return err
	}

	return nil
}

// Key is made of $(path)_$(method)
func (trc *GRegisterController) Key() string {
	return trc.R.URL.Path + "_" + trc.R.Method
}
