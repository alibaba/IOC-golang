proto-gen:
	go get -u github.com/gogo/protobuf@v1.3.2
	protoc --go_out=./extension/aop/trace/api --go-grpc_out=./extension/aop/trace/api -I${GOPATH}/pkg/mod/github.com/gogo/protobuf@v1.3.2 -I./extension/aop/trace/api/jaeger/jaeger_model.proto -I. ./extension/aop/trace/api/ioc_golang/aop/trace/trace.proto
	protoc --go_out=./extension/aop/watch/api --go-grpc_out=./extension/aop/watch/api ./extension/aop/watch/api/ioc_golang/aop/watch/watch.proto
	protoc --go_out=./extension/aop/list/api --go-grpc_out=./extension/aop/list/api ./extension/aop/list/api/ioc_golang/aop/list/list.proto
	protoc --go_out=./extension/aop/monitor/api --go-grpc_out=./extension/aop/monitor/api ./extension/aop/monitor/api/ioc_golang/aop/monitor/monitor.proto

mockery-gen:
	cd extension/aop/monitor && sudo mockery --name=interceptor --inpackage  --filename=interceptor_mock.go --structname=mockInterceptor

tidy-all:
	cd extension && go mod tidy -compat=1.17
	cd example && go mod tidy -compat=1.17
	cd iocli && go mod tidy -compat=1.17
	go mod tidy -compat=1.17

imports:
	goimports -local github.com/alibaba/ioc-golang -w .

lint: tidy-all
	golangci-lint run

test-all:
	go test ./... -cover -p 1
	cd extension && go test ./... -cover -p 1
	cd example && go test ./... -cover -p 1
	cd iocli && go test ./... -cover -p 1