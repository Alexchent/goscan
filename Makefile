.PHONY: default local win mac

default local:
	CGO_ENABLED=0 go build -o gscan main.go
    CGO_ENABLED=0 go build -o sscan sscan.go

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gscan main.go

win:
   CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o gscan main.go

mac:
   CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o gscan main.go