.PHONY: default local win mac install

# 版本号：优先取最近一次 git tag，否则回退到 dev
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null | sed 's/^v//' || echo dev)

LDFLAGS = -X github.com/Alexchent/goscan/cmd.version=$(VERSION)

default:
	go build -ldflags "$(LDFLAGS)" -o scan github.com/Alexchent/goscan

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o scan github.com/Alexchent/goscan

win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o scan github.com/Alexchent/goscan

mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o scan github.com/Alexchent/goscan

install:
	go install github.com/Alexchent/goscan@latest