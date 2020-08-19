.PHONY: build

APPNAME ?= pi-wifi
DOCKER_IMAGE ?= opny/${APPNAME}

BUILD_PATH ?= ./build

build: build/amd64 build/arm64 build/arm

build/amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BUILD_PATH}/amd64 .

build/arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ${BUILD_PATH}/arm64 .

build/arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o ${BUILD_PATH}/arm7 .
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o ${BUILD_PATH}/arm6 .

docker/build: build docker/build/amd64 docker/build/arm64 docker/build/arm

docker/build/manifest:
	docker manifest create \
		${DOCKER_IMAGE} \
		--amend ${DOCKER_IMAGE}-amd64 \
		--amend ${DOCKER_IMAGE}-arm64 \
		--amend ${DOCKER_IMAGE}-arm6 \
		--amend ${DOCKER_IMAGE}-arm7
	
	docker manifest annotate ${DOCKER_IMAGE} ${DOCKER_IMAGE}-amd64 --arch amd64
	docker manifest annotate ${DOCKER_IMAGE} ${DOCKER_IMAGE}-arm64 --arch arm64
	docker manifest annotate ${DOCKER_IMAGE} ${DOCKER_IMAGE}-arm6 --arch arm --variant v6
	docker manifest annotate ${DOCKER_IMAGE} ${DOCKER_IMAGE}-arm7 --arch arm --variant v7

	docker manifest push ${DOCKER_IMAGE}

docker/build/amd64:
	docker build . -t ${DOCKER_IMAGE}-amd64 --build-arg ARCH=amd64

docker/build/arm64:
	docker build . -t ${DOCKER_IMAGE}-arm64 --build-arg ARCH=arm64

docker/build/arm:
	docker build . -t ${DOCKER_IMAGE}-arm6 --build-arg ARCH=arm7
	docker build . -t ${DOCKER_IMAGE}-arm7 --build-arg ARCH=arm6


docker/push: docker/build docker/push/amd64 docker/push/arm64 docker/push/arm docker/build/manifest
	docker manifest push ${DOCKER_IMAGE}

docker/push/amd64:
	docker push ${DOCKER_IMAGE}-amd64

docker/push/arm64:
	docker push ${DOCKER_IMAGE}-arm64

docker/push/arm:
	docker push ${DOCKER_IMAGE}-arm6
	docker push ${DOCKER_IMAGE}-arm7
