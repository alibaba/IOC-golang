package main

import (
	"github.com/sirupsen/logrus"

	"github.com/alibaba/ioc-golang"
)

func main() {
	ioc.Load()
	logrus.Info("I love IOCddd")
}
