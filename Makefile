GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod
GOFMT = $(GOCMD) fmt
GOVET = $(GOCMD) vet

build: darwin web

all: darwin linux web

linux-full: linux web

darwin-full: darwin web

darwin:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -a -o bin/aries.darwin backend/*.go

web:
	mkdir -p bin/
	cd frontend/; yarn install; yarn build
	rm -rf bin/public
	mv frontend/dist bin/public

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -a -installsuffix cgo -o bin/aries.linux backend/*.go

clean:
	$(GOCLEAN)
	rm -rf bin

dep:
	cd backend; $(GOGET) -u
	$(GOMOD) tidy
	$(GOMOD) verify

fmt:
	cd backend; $(GOFMT)

vet:
	cd backend; $(GOVET)

