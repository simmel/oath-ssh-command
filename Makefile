ifndef GOPATH
export GOPATH=$(HOME)/code/go
endif
export PATH := $(PATH):$(GOPATH)/bin

oath-ssh-command: .deps main.go
	go build

$(GOPATH):
	@mkdir -p $(GOPATH)

.deps: Godeps/Godeps.json
	env
	go get github.com/tools/godep
	PATH=$(PATH) godep restore
	touch .deps

run: oath-ssh-command
	./oath-ssh-command

clean:
	rm -rf $(GOPATH)/{bin,pkg,src} \
	oath-ssh-command

test: oath-ssh-command
	go test -v
