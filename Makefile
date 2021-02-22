XC_OS="linux"
XC_ARCH="386 amd64 arm arm64"
XC_PARALLEL="4"
BIN="../bin"
SRC=$(shell find . -name "*.go")

ifeq (, $(shell which gox))
$(warning "could not find gox in $(PATH), run: go get github.com/mitchellh/gox")
endif

.PHONY: all build

default: all

all: build

build:
	gox -os=$(XC_OS) -arch=$(XC_ARCH) -parallel=$(XC_PARALLEL) -output=$(BIN)/{{.Dir}}_{{.OS}}_{{.Arch}};
