-include .env

VERSION := v0.1
GIT_REV := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell $(date +'%Y-%m-%d_%T'))
PROJECT_NAME := $(shell basename "$(PWD)")


# Go related variables.
GOBASE := $(shell pwd)
GOPATH := $(GOBASE)/vendor:$(GOBASE)
BUILD_PATH := $(GOBASE)/build
GOBIN := $(BUILD_PATH)/bin
GOFILES := transporter.go

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags '-w -s -extldflags "-static" -X act.buaa.edu.cn/jcspan/transporter/util.Version=$(VERSION) -X act.buaa.edu.cn/jcspan/transporter/util.GitRev=$(GIT_REV) -X act.buaa.edu.cn/jcspan/transporter/util.BuildTime=$(BUILD_TIME)'

objects = transporter

build: &(objects)
	@echo "  >  Building binary done"


&(objects):
	@echo "  >  Building binary start..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME) $(GOFILES)

clean:
	@-rm $(GOBIN)/$(PROJECTNAME) 2> /dev/null
	@-$(MAKE) go-clean

go-clean:
	@echo "  >  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean
