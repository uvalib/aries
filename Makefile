GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

build: build-darwin build-linux

all: deps build-darwin build-linux build-web deploy-web

build-darwin:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -a -o bin/aries.darwin cmd/aries/*.go

build-web:
	mkdir -p bin/
	cd frontend/; yarn build

deploy-web:
	mv frontend/dist bin/public

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -a -installsuffix cgo -o bin/aries.linux cmd/aries/*.go
	cp services.csv bin/services.csv

clean:
	$(GOCLEAN)
	rm -rf bin

deps:
	dep ensure
	dep status
