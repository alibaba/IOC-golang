package common

import "fmt"

func GetJaegerCollectorEndpoint(jaegerCollectorAddress string) string {
	return fmt.Sprintf("http://%s/api/traces?format=jaeger.thrift", jaegerCollectorAddress)
}
