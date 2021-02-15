.PHONY: all clean fmt build run

MOD = $(shell sed -n -e 's/^module \(.*\)/\1/p' go.mod)
PORT ?= 8080
DEV_PORT ?= 8081

all: build

# clean
#
clean:
	git ls-files -o static/ | xargs -rt rm
	go mod tidy

# fmt
#
fmt:
	find -name '*.go' | xargs -rt gofmt -l -w -s

# build
#
.PHONY: go-build go-generate npm-build

go-build:
	go get -v ./...

go-generate: assets/files.go

assets/files.go: FORCE
	[ -x "`which file2go`" ] || go get -v github.com/amery/file2go/cmd/file2go
	cd assets; find * -type f ! -name .gitignore -a ! -name '*.go' \
		| sort -V | xargs -t file2go -p assets -o files.go

npm-build:
	npm run build

build: npm-build go-generate go-build

# run
#
run:
	go run -v ./cmd/server -p $(PORT) -t 0

.PHONY: FORCE
FORCE:
