FROM golang:1.12-alpine AS go-builder

LABEL stage=build

ENV GO111MODULE on

ARG GOLANG_PROXY_IP
ENV GOPROXY=http://${GOLANG_PROXY_IP}:3000

WORKDIR /go/src/kubemonitor

COPY ./kubemonitor /go/src/kubemonitor
RUN apk add git && \
    cd /go/src/kubemonitor && \
    CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -o kubemonitor cmd/kubemonitor/main.go

FROM alpine AS base
WORKDIR /kubemonitor
COPY --from=go-builder /go/src/kubemonitor/kubemonitor /kubemonitor
RUN apk add --no-cache util-linux
RUN wget -O /kubectl https://storage.googleapis.com/kubernetes-release/release/v1.20.0/bin/linux/amd64/kubectl && \
    chmod +x /kubectl && \
    cp /kubectl /usr/local/bin/
RUN mkdir -p /kubemonitor/cache


ENTRYPOINT [ "/kubemonitor/kubemonitor" ]
