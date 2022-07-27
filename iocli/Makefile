tidy:
	go mod tidy -compat=1.17

test:
	go test ./...

lint: tidy
	golangci-lint run

build-all-platform:
	mkdir -p ./.release/linux-amd64 && GOOS=linux GOARCH=amd64 go build -o ./.release/linux-amd64 && cd .release && tar -czvf ./iocli-linux-amd64.tar.gz ./linux-amd64
	mkdir -p ./.release/linux-arm64 && GOOS=linux GOARCH=arm64 go build -o ./.release/linux-arm64 && cd .release  && tar -czvf ./iocli-linux-arm64.tar.gz ./linux-arm64
	mkdir -p ./.release/darwin-amd64 && GOOS=darwin GOARCH=amd64 go build -o ./.release/darwin-amd64 && cd .release  && tar -czvf ./iocli-darwin-amd64.tar.gz ./darwin-amd64
	mkdir -p ./.release/darwin-arm64 && GOOS=darwin GOARCH=arm64 go build -o ./.release/darwin-arm64 && cd .release  && tar -czvf ./iocli-darwin-arm64.tar.gz ./darwin-arm64
	mkdir -p ./.release/windows-amd64 && GOOS=windows GOARCH=amd64 go build -o ./.release/windows-amd64 && cd .release  && tar -czvf ./iocli-windows-amd64.tar.gz ./windows-amd64
	mkdir -p ./.release/windows-arm64 && GOOS=windows GOARCH=arm64 go build -o ./.release/windows-arm64 && cd .release  && tar -czvf ./iocli-windows-arm64.tar.gz ./windows-arm64
