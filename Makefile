.PHONY: default local win mac

default:
	CGO_ENABLED=0 go build -o scan main.go
   
linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o scan main.go

win:
   CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o scan main.go

mac:
   CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o scan main.go