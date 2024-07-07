PKG=github.com/parkr/webhook-to-mqtt
REV:=$(shell git rev-parse HEAD)

all: fmt build test

fmt:
	go fmt $(PKG)/...

build:
	go build $(PKG)

test:
	go test -v ./...

docker-build:
	docker build -t parkr/webhook-to-mqtt:$(REV) .

docker-test: docker-build
	docker run --name=webhook-to-mqtt-test --rm -it --net=host parkr/webhook-to-mqtt:$(REV)

docker-release: docker-build
	docker push parkr/webhook-to-mqtt:$(REV)

dive: docker-build
	dive parkr/webhook-to-mqtt:$(REV)
