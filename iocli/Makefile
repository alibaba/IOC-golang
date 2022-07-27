tidy:
	go mod tidy -compat=1.17

test:
	go test ./...

lint: tidy
	golangci-lint run

build-all-platform:
	mkdir -p ./release/linux-amd64 && GOOS=linux GOARCH=amd64 go build -o ./release/linux-amd64
	mkdir -p ./release/linux-arm64 && GOOS=linux GOARCH=arm64 go build -o ./release/linux-arm64
	mkdir -p ./release/darwin-amd64 && GOOS=darwin GOARCH=amd64 go build -o ./release/darwin-amd64
	mkdir -p ./release/darwin-arm64 && GOOS=darwin GOARCH=arm64 go build -o ./release/darwin-arm64
	mkdir -p ./release/windows-amd64 && GOOS=windows GOARCH=amd64 go build -o ./release/windows-amd64
	mkdir -p ./release/windows-arm64 && GOOS=windows GOARCH=arm64 go build -o ./release/windows-arm64
