package ghttp

type httpCode uint32

const (
	UnsetHttpCode          = httpCode(0)
	DefaultHttpSuccessCode = httpCode(200)
	DefaultHttpErrorCode   = httpCode(500)
)
