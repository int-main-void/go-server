#
# go service
#
# depends on having go and glide installed in the environment
#

.DEFAULT_GOAL := all


IMAGE_DIR=image
PROJECT_NAME=mwn-go-server
VERSION=dev-latest

DOCKER_OS=linux
DOCKER_ARCH=386

#IMAGE_ARTIFACTS=bin pkg conf tools test
IMAGE_ARTIFACTS=bin conf 

export GOPATH=$(PWD)

$(info GOPATH is $(GOPATH))

clean:
	rm -rf $(IMAGE_DIR)
	rm -rf bin/
	rm -rf pkg/
	rm -rf src/vendor

src/vendor:
	cd src && glide install

build: src/vendor
	GOOS=$(DOCKER_OS) GOARCH=$(DOCKER_ARCH) go install greatsagemonkey.com/goservice/go-service

image: build
	mkdir -p $(IMAGE_DIR)/root
	cp conf/Dockerfile $(IMAGE_DIR)
	cp -a $(IMAGE_ARTIFACTS) $(IMAGE_DIR)/root
	docker build -t $(PROJECT_NAME):$(VERSION) $(IMAGE_DIR)

unit-test:
	go test -cover greatsagemonkey.com/...

godoc:
	godoc

lint:
	go get -u github.com/golang/lint/golint
	./bin/golint src/greatsagemonkey.com/...

vet:
	go tool vet -all src/greatsagemonkey.com

validate: vet lint

all: image

run: image
	docker-compose -f conf/docker-compose.yml up

bin-local:
	go install greatsagemonkey.com/goservice/go-service

run-local: bin-local
	CONFIG_FILENAME=conf/go-service-conf.json STAGE=dev ./bin/go-service
