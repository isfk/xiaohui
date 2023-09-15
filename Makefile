AUTHOR=sfk@live.cn
DATE=$(shell date '+%Y-%m-%d_%H:%M')
VERSION=$(shell git describe --tags --always)
BRANCH=$(shell git symbolic-ref --short HEAD)

.PHONY: install
install:
	go install github.com/cosmtrek/air@latest
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/srikrsna/protoc-gen-gotag@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install git.isfk.cn/isfk/protoc-gen-echo@latest

.PHONY: build
build: clean
	go build -ldflags "-X main.Author=$(AUTHOR) -X main.Date=$(DATE) -X main.Version=$(VERSION) -X main.Branch=$(BRANCH)" -o ./bin .


.PHONY: build_linux
build_linux:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "-X main.Author=$(AUTHOR) -X main.Date=$(DATE) -X main.Version=$(VERSION) -X main.Branch=$(BRANCH)" -o ./bin .

.PHONY: run
run: build
	./bin/ares version

.PHONY: aircron
aircron:
	air cron

.PHONY: clean
clean:
	rm -rf ./bin/*

.PHONY: info
info:
	echo $(AUTHOR) > ./cmd/author.txt
	echo $(DATE) > ./cmd/date.txt
	echo $(VERSION) > ./cmd/version.txt
	echo $(BRANCH) > ./cmd/branch.txt