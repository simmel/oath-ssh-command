GOPATH=$(shell pwd)

oath-ssh-command: .depman.cache main.go
	go build

.depman.cache: deps.json
	go get github.com/vube/depman
	./bin/depman

run: oath-ssh-command
	./oath-ssh-command

clean:
	rm -rf .depman.cache \
	bin \
	pkg \
	src

test: oath-ssh-command
	go test
