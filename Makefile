.PHONY: all clean deps build run dev

MOD = $(shell sed -n -e 's/^module \(.*\)/\1/p' go.mod)
PORT ?= 8080
DEV_PORT ?= 8081
SERVER ?= server

B = $(CURDIR)/build

ASSETS_GO_FILE = assets/files.go

GENERATED_GO_FILES = $(ASSETS_GO_FILE)

all: build

# deps
#
.PHONY: go-deps npm-deps

deps: go-deps npm-deps

GOBIN = $(GOPATH)/bin
NPM_BIN_PREFIX = $(CURDIR)/node_modules/.bin

FILE2GO = $(GOBIN)/file2go
WEBPACK = $(NPM_BIN_PREFIX)/webpack

$(FILE2GO):
	go get -v github.com/amery/file2go/cmd/file2go

$(NPM_BIN_PREFIX)/%:
	npm i
	npm shrinkwrap

GO_DEPS = $(FILE2GO)
NPM_DEPS = $(WEBPACK)

go-deps: $(GO_DEPS)
npm-deps: $(NPM_DEPS)

# clean
#
clean:
	go clean -x -r -modcache
	git ls-files -o assets/ | xargs -rt rm
	rm -rf $(B) node_modules/
	rm -f $(GENERATED_GO_FILES) npm-shrinkwrap.json package-lock.json

# fmt
#
.PHONY: fmt lint npm-lint go-fmt

fmt: go-fmt npm-lint
lint: go-fmt npm-lint

go-fmt: go-deps
	find -name '*.go' | xargs -rt gofmt -l -w -s

npm-lint: npm-deps
	npm run lint

#
# npm.files -> npm.mk -> npm.built -> assets.files -> assets.mk -> assets/file.go
#

# file listings
#
ASSETS_FILES_FILE = $(B)/assets.files
GO_FILES_FILE     = $(B)/go.files
NPM_FILES_FILE    = $(B)/npm.files

ALL_FILE_FILES    = $(ASSETS_FILES_FILE) \
                    $(GO_FILES_FILE) \
                    $(NPM_FILES_FILE)

ASSETS_FILES_MK   = $(B)/assets.mk
GO_FILES_MK       = $(B)/go.mk
NPM_FILES_MK      = $(B)/npm.mk

ALL_FILE_MK_FILES = $(ASSETS_FILES_MK) \
                    $(GO_FILES_MK) \
                    $(NPM_FILES_MK)

# extra markers
NPM_BUILT_MARK    = $(B)/npm.built

$(ASSETS_FILES_FILE): $(NPM_BUILT_MARK)
$(ASSETS_FILES_FILE): FIND = assets/ ! -type d -a ! -name '*.go' -a ! -name '.gitignore'
$(GO_FILES_FILE):     FIND = * -name '*.go'
$(GO_FILES_FILE):     EXTRA = $(GENERATED_GO_FILES)
$(NPM_FILES_FILE):    FIND = src/ ! -type d

$(ASSETS_FILES_MK):   $(ASSETS_FILES_FILE)
$(GO_FILES_MK):       $(GO_FILES_FILE)
$(NPM_FILES_MK):      $(NPM_FILES_FILE)

$(ASSETS_FILES_MK):   PREFIX=ASSETS
$(GO_FILES_MK):       PREFIX=GO
$(NPM_FILES_MK):      PREFIX=NPM

$(ALL_FILE_FILES): FORCE
	@mkdir -p $(dir $@)
	@(for x in $(EXTRA); do echo $$x; done; find $(FIND)) | sed -e '/^[ \t]*$$/d;' | sort -uV > $@~
	@if ! cmp -s $@ $@~; then \
		mv $@~ $@; \
		echo $(@:$(CURDIR)/%=%) updated.; \
	else \
		rm $@~; \
	fi

$(ALL_FILE_MK_FILES): Makefile
$(ALL_FILE_MK_FILES):
	@echo "$(PREFIX)_FILES = $$(cat $< | tr '\n' ' ')" > $@~
	@mv $@~ $@
	@echo $(@:$(CURDIR)/%=%) updated.;

# npm-build
#
.PHONY: npm-build

npm-build: npm-deps
	npm run build
	@touch $(NPM_BUILT_MARK)

include $(NPM_FILES_MK)

$(NPM_BUILT_MARK): $(NPM_FILES_MK) $(NPM_FILES) npm-deps
	npm run build
	@touch $(NPM_BUILT_MARK)

# go-build
#
.PHONY: go-build

include $(ASSETS_FILES_MK)
include $(GO_FILES_MK)

$(ASSETS_GO_FILE): $(ASSETS_FILES_FILE) $(FILE2GO) $(ASSETS_FILES)
	@cut -d/ -f2- < $< | (cd $(@D); xargs -t file2go -p assets -o $(@F))
	@echo $(@:$(CURDIR)/%=%) updated.;

go-build: $(GO_FILES) go-deps
	go get -v ./...

# build
#
build: go-build

# run
#
run: go-deps
	go run -v ./cmd/$(SERVER) -p $(PORT) -t 0

dev: go-build npm-deps
	set -x; $(GOBIN)/$(SERVER) -p $(DEV_PORT) --dev & trap "kill $$!" EXIT; env HOST=0.0.0.0 PORT=$(PORT) BACKEND=$(DEV_PORT) npm start

#
#
.PHONY: FORCE
FORCE:
