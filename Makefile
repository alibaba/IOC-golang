proto-gen:
	protoc --go_out=./aop/api --go-grpc_out=./aop/api ./aop/api/ioc_golang/aop/debug.proto

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