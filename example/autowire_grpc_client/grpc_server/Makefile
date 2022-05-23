IMAGE = laurencelizhixin/grpc-server-sayhello
TAG = 1.0.0

buildx-publish:
	GOARCH=amd64 GOOS=linux go build -o app_amd64 .
	GOARCH=arm64 GOOS=linux go build -o app_arm64 .
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		-t ${IMAGE}:${TAG} \
		. --push
	rm ./app*