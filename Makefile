.PHONY: default local
default:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gscan

local:
	CGO_ENABLED=0 go build -o gscan