SHELL=/bin/bash
name=kubemonitor
tag=latest
GOProxyIP=10.60.6.54
all: build_docker_image

build_docker_image:
	@echo Start build image
	docker build --build-arg GOLANG_PROXY_IP=$(GOProxyIP) -t=$(name):$(tag) -f build/package/Dockerfile ..
	#docker build -t=$(name):$(tag) -f build/package/Dockerfile ..
	docker image prune --filter label=stage=build -f
