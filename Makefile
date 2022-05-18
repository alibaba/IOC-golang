proto-gen:
	protoc --go_out=./debug/api --go-grpc_out=./debug/api ./debug/api/ioc_golang/boot/debug.proto

tidy:
	go mod tidy -compat=1.17

test:
	go test ./... -cover

imports:
	goimports -local github.com/alibaba/ioc-golang -w .

lint: tidy
	golangci-lint run