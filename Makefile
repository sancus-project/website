.PHONY: all clean fmt build run

MOD = $(shell sed -n -e 's/^module \(.*\)/\1/p' go.mod)
PORT ?= 8080
DEV_PORT ?= 8081

all: build

# clean
#
clean:
	git ls-files -o static/ | xargs -rt rm

# fmt
#
fmt:
	find -name '*.go' | xargs -rt gofmt -l -w -s

# build
#
.PHONY: go-build npm-build

go-build:
	go get -v ./...

npm-build:
	npm run build

build: npm-build go-build

# run
#
run:
	go run -v ./cmd/server -p $(PORT) -t 0
