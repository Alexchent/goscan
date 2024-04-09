.PHONY: default local win mac

default local:
	CGO_ENABLED=0 go build -o gscan

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gscan

win:
   CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o gscan

mac:
   CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o gscan