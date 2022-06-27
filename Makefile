proto-gen:
	protoc --go_out=./debug/api --go-grpc_out=./debug/api ./debug/api/ioc_golang/debug/debug.proto

tidy-all:
	cd extension && go mod tidy -compat=1.17
	cd example && go mod tidy -compat=1.17
	cd iocli && go mod tidy -compat=1.17
	go mod tidy -compat=1.17

.PHONY:test
test:
	go test ./... -cover -p 1

imports:
	goimports -local github.com/alibaba/ioc-golang -w .

lint: tidy-all
	golangci-lint run