.PHONY: test clean build

ARCH := $(shell uname -m)
ifeq ($(ARCH),x86_64)
	ARCH := amd64
endif
OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
VERSION := $(shell curl -s https://api.github.com/repos/Schumann-IT/go-ieftool/releases/latest | grep "tag_name" | awk '{print $$2}' | sed 's|[\"\,]*||g')

clean:
	@rm -Rf build
	@rm -f ./ieftool

test:
	@go test -v ./internal ./cmd

build:
	@go build -o ieftool

install: ieftool
	@sudo mv ieftool /usr/local/bin/ieftool

ieftool:
	@curl -s -L -o ieftool https://github.com/Schumann-IT/go-ieftool/releases/download/$(VERSION)/ieftool-$(OS)-$(ARCH)
	@chmod +x ieftool
	@if [ "$(OS)" = "darwin" ]; then\
        xattr -d com.apple.quarantine ./ieftool /dev/null 2>&1 | true; \
    fi