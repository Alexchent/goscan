.PHONY: default local win mac install

# 版本号：优先取最近一次 git tag，否则回退到 dev
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null | sed 's/^v//' || echo dev)

LDFLAGS = -X github.com/Alexchent/goscan/cmd.version=$(VERSION)

default:
	go build -ldflags "$(LDFLAGS)" -o scan github.com/Alexchent/goscan/cmd/scan
	go build -ldflags "$(LDFLAGS)" -o sscan github.com/Alexchent/goscan/cmd/sscan

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o scan github.com/Alexchent/goscan/cmd/scan

win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o scan github.com/Alexchent/goscan/cmd/scan

mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o scan github.com/Alexchent/goscan/cmd/scan

install:
	@echo "先用 go install 安装两个二进制："
	go install github.com/Alexchent/goscan/cmd/scan@latest
	go install github.com/Alexchent/goscan/cmd/sscan@latest
