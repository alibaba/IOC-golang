FROM alpine:3.15

ARG TARGETARCH

WORKDIR /server
COPY "./app_${TARGETARCH}" /server/app

ENTRYPOINT /server/app
