package common

import (
	"fmt"
	"golang.org/x/net/context"
)

type InterceptorContext struct {
	context.Context
	sdid    string
	method  string
	isParam bool
}

func NewInterceptorContext(ctx context.Context, sdid, method string, isParam bool) *InterceptorContext {
	return &InterceptorContext{
		Context: ctx,
		sdid:    sdid,
		method:  method,
		isParam: isParam,
	}
}

// GetSDID get struct description id
func (i *InterceptorContext) GetSDID() string {
	return i.sdid
}

// GetMethod get current invoking method name
func (i *InterceptorContext) GetMethod() string {
	return i.method
}

// IsParam get if current values is param, or response value.
func (i *InterceptorContext) IsParam() bool {
	return i.isParam
}

func (i *InterceptorContext) String() string {
	return fmt.Sprintf("SDID: %s, Method: %s, IsParam: %t", i.sdid, i.method, i.isParam)
}
