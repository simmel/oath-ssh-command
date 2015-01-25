GOPATH=$(shell pwd)

totp-ssh-command: main.go
	go build
