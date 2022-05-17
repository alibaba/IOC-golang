tidy:
	go mod tidy -compat=1.17

test:
	go test ./...

imports:
	imports-formatter .

lint: tidy
	golangci-lint run