package app

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type App struct {
	// inject main.ServiceImpl1 pointer to Service interface with proxy wrapper
	ServiceImpl1 ServiceImpl1IOCInterface `singleton:""`

	// inject main.ServiceImpl2 pointer to Service interface with proxy wrapper
	ServiceImpl2 ServiceImpl2IOCInterface `singleton:""`
}

func (a *App) Run() {
	for {
		time.Sleep(time.Second * 3)
		/*
			ServiceImpl1.GetHelloString() calls
				ServiceImpl2.GetHelloString()
		*/
		logrus.Println(a.ServiceImpl1.GetHelloString("laurence"))

		/*
			ServiceImpl2.GetHelloString()
		*/
		logrus.Println(a.ServiceImpl2.GetHelloString("laurence"))
	}
}

type Service interface {
	GetHelloString(string) string
}

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type ServiceImpl1 struct {
	ServiceImpl2 ServiceImpl2IOCInterface `singleton:""`
}

func (s *ServiceImpl1) GetHelloString(name string) string {
	str := fmt.Sprintf("This is ServiceImpl1, hello %s", name)
	logrus.Debugf(str)
	logrus.Info(str)
	logrus.Warnf(str)
	logrus.Errorf(str)

	// call service2
	s.ServiceImpl2.GetHelloString(name)
	return fmt.Sprintf("This is ServiceImpl1, hello %s", name)
}

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type ServiceImpl2 struct {
}

func (s *ServiceImpl2) GetHelloString(name string) string {
	str := fmt.Sprintf("This is ServiceImpl2, hello %s", name)
	logrus.Debugf(str)
	logrus.Info(str)
	logrus.Warnf(str)
	logrus.Errorf(str)
	return str
}
