.PHONY: all clean fmt build run

MOD = $(shell sed -n -e 's/^module \(.*\)/\1/p' go.mod)
PORT ?= 8080
DEV_PORT ?= 8081

B = $(CURDIR)/build

ASSETS_FILES_FILE = $(B)/assets.files
ASSETS_FILES_MK = $(B)/assets.mk

GO_FILES_FILE = $(B)/go.files
GO_FILES_MK = $(B)/go.mk

NPM_FILES_FILE = $(B)/npm.files
NPM_FILES_MK = $(B)/npm.mk
NPM_BUILT_MARK = $(B)/npm.built

ALL_FILE_FILES = $(ASSETS_FILES_FILE) \
		 $(GO_FILES_FILE) \
		 $(NPM_FILES_FILE)

ALL_FILE_MK_FILES = $(ASSETS_FILES_MK) \
		    $(GO_FILES_MK) \
		    $(NPM_FILES_MK)

ASSETS_GO_FILE = assets/files.go

GENERATED_GO_FILES = $(ASSETS_GO_FILE)

all: build

# helpers
#
$(GOPATH)/bin/file2go:
	go get -v github.com/amery/file2go/cmd/file2go

# clean
#
clean:
	go clean -x -r -modcache
	git ls-files -o assets/ | xargs -rt rm
	rm -rf $(B)
	rm -f $(GENERATED_GO_FILES)

# fmt
#
fmt:
	find -name '*.go' | xargs -rt gofmt -l -w -s

#
# npm.files -> npm.mk -> npm.built -> assets.files -> assets.mk -> assets/file.go
#

# file listings
#
$(ASSETS_FILES_FILE): $(NPM_BUILT_MARK)
$(ASSETS_FILES_FILE): BASE=assets/

$(NPM_FILES_FILE): BASE=src/

$(ASSETS_FILES_FILE) $(NPM_FILES_FILE): FORCE
	@mkdir -p $(dir $@)
	@find $(BASE) ! -type d -a ! -name '*.go' -a ! -name '.gitignore' | sort -V > $@~
	@if ! cmp -s $@ $@~; then \
		mv $@~ $@; \
		echo $@ updated; \
	else \
		rm $@~; \
	fi

$(GO_FILES_FILE): FORCE
	@mkdir -p $(dir $@)
	@(find * -name '*.go'; echo $(GENERATED_GO_FILES) | tr ' ' '\n') | sed -e '/^[ \t]*$$/d;' | sort -uV > $@~
	@if ! cmp -s $@ $@~; then \
		mv $@~ $@; \
		echo $@ updated; \
	else \
		rm $@~; \
	fi

$(ASSETS_FILES_MK): $(ASSETS_FILES_FILE)
$(ASSETS_FILES_MK): PREFIX=ASSETS

$(GO_FILES_MK): $(GO_FILES_FILE)
$(GO_FILES_MK): PREFIX=GO

$(NPM_FILES_MK): $(NPM_FILES_FILE)
$(NPM_FILES_MK): PREFIX=NPM

$(ALL_FILE_MK_FILES): Makefile
$(ALL_FILE_MK_FILES):
	echo "$(PREFIX)_FILES = $$(cat $< | tr '\n' ' ')" > $@~
	mv $@~ $@

# npm-build
#
.PHONY: npm-build

npm-build:
	npm run build
	touch $(NPM_BUILT_MARK)

include $(NPM_FILES_MK)

$(NPM_BUILT_MARK): $(NPM_FILES_MK) $(NPM_FILES)
	npm run build
	touch $(NPM_BUILT_MARK)

# go-build
#
.PHONY: go-build

include $(ASSETS_FILES_MK)
include $(GO_FILES_MK)

$(ASSETS_GO_FILE): $(ASSETS_FILES_FILE) $(GOPATH)/bin/file2go $(ASSETS_FILES)
	cut -d/ -f2- < $< | (cd $(@D); xargs -t file2go -p assets -o $(@F))

go-build: $(GO_FILES)
	go get -v ./...

# build
#
build: go-build

# run
#
run:
	go run -v ./cmd/server -p $(PORT) -t 0

.PHONY: FORCE
FORCE:
