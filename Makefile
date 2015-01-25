GOPATH=$(shell pwd)

totp-ssh-command: .depman.cache main.go
	go build

.depman.cache: deps.json
	go get github.com/vube/depman
	./bin/depman
