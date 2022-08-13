package aop

import (
	"fmt"

	"github.com/alibaba/ioc-golang/aop"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=initAfterIsCalledAfterPanicTestInterceptor

type afterIsCalledAfterPanicTestInterceptor struct {
	AfterIsCalledNum int
}

func initAfterIsCalledAfterPanicTestInterceptor(a *afterIsCalledAfterPanicTestInterceptor) (*afterIsCalledAfterPanicTestInterceptor, error) {
	a.AfterIsCalledNum = 0
	return a, nil
}

func (p *afterIsCalledAfterPanicTestInterceptor) GetAfterIsCalledNum() int {
	return p.AfterIsCalledNum
}

func (p *afterIsCalledAfterPanicTestInterceptor) BeforeInvoke(ctx *aop.InvocationContext) {

}

func (p *afterIsCalledAfterPanicTestInterceptor) AfterInvoke(ctx *aop.InvocationContext) {
	p.AfterIsCalledNum++
}

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type PanicAfterCalledTestSubApp struct {
}

func (p *PanicAfterCalledTestSubApp) RunWithPanic(panicMsg string) {
	panic(panicMsg)
}

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type PanicAfterCalledTestApp struct {
	PanicAfterCalledTestSubApp PanicAfterCalledTestSubAppIOCInterface `singleton:""`
}

func (p *PanicAfterCalledTestApp) RunWithPanic(panicMst string) (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = fmt.Sprintf("%+v", r)
		}
	}()
	p.PanicAfterCalledTestSubApp.RunWithPanic(panicMst)
	return ""
}
