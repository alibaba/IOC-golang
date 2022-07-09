tidy:
	go mod tidy -compat=1.17

test:
	go test ./...

lint: tidy
	golangci-lint run