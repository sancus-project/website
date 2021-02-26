.PHONY: all clean deps build run dev

# locations
#
B = $(CURDIR)/build
GOBIN = $(GOPATH)/bin
NPX_BIN = $(CURDIR)/node_modules/.bin

# config files
#
MODD_RUN_CONF = $(B)/modd-run.conf
MODD_DEV_CONF = $(B)/modd-dev.conf

# tools
#
GO = go
GOFMT = gofmt
GOFMT_FLAGS = -w -l -s
GOGET = $(GO) get
GOGET_FLAGS = -v
NPM = npm

FILE2GO = $(GOBIN)/file2go
MODD = $(GOBIN)/modd
MODD_FLAGS = -b
WEBPACK = $(NPX_BIN)/webpack

# magic constants
#
MOD = $(shell sed -n -e 's/^module \(.*\)/\1/p' go.mod)
PORT ?= 8080
DEV_PORT ?= 8081
SERVER ?= server

# generated files
#
ASSETS_GO_FILE = assets/files.go

GENERATED_GO_FILES = $(ASSETS_GO_FILE)

# default target
all: build

# deps
#
.PHONY: go-deps npm-deps
deps: go-deps npm-deps

# go-deps
GO_DEPS = $(FILE2GO) $(MODD)

go-deps: $(GO_DEPS)

$(FILE2GO): URL=github.com/amery/file2go/cmd/file2go
$(MODD): URL=github.com/cortesi/modd/cmd/modd

$(GO_DEPS):
	$(GOGET) $(GOGET_FLAGS) $(URL)

# npm-deps
NPM_DEPS = $(WEBPACK)

$(NPXBIN)/%:
	$(NPM) i
	$(NPM) shrinkwrap

npm-deps: $(NPM_DEPS)

# clean
#
clean:
	$(GO) clean -x -r -modcache
	git ls-files -o assets/ | xargs -rt rm
	rm -rf $(B) node_modules/
	rm -f $(GENERATED_GO_FILES) npm-shrinkwrap.json package-lock.json

# fmt
#
.PHONY: fmt lint npm-lint go-fmt

fmt: go-fmt npm-lint
lint: go-fmt npm-lint

go-fmt: $(GO_DEPS)
	find -name '*.go' | xargs -rt $(GOFMT) $(GOFMT_FLAGS)

npm-lint: $(NPN_DEPS)
	$(NPM) run lint

# run
#
$(MODD_RUN_CONF): src/modd/run.conf
$(MODD_DEV_CONF): src/modd/dev.conf

$(MODD_RUN_CONF) $(MOD_DEV_CONF):
	sed \
		-e "s|@@PORT@@|$(PORT)|g" \
		-e "s|@@BACKEND@@|$(DEV_PORT)|g" \
		-e "s|@@NPM@@|$(NPM)|g" \
		-e "s|@@GO@@|$(GO)|g" \
		-e "s|@@GOFMT@@|$(GOFMT) $(GOFMT_FLAGS)|g" \
		-e "s|@@GOGET@@|$(GOGET)|g" \
		-e "s|@@FILE2GO@@|$(notdir $(FILE2GO))|g" \
		-e "s|@@SERVER@@|$(SERVER)|g" \
		$< > $@~
	mv $@~ $@

run: $(MODD_RUN_CONF)
dev: $(MODD_DEV_CONF)

run dev: $(MODD) go-deps npm-deps
run dev:
	$(MODD) $(MODD_FLAGS) -f $<

# build
#
ASSETS_FILES_FILTER = find $(dir $(ASSETS_GO_FILE)) -type f -a ! -name '*.go' -a ! -name '.*' -a ! -name '*~'
NPM_FILES_FILTER = find src/ -name '*.js' -o -name '*.scss'
GO_FILES_FILTER = find */ -name *.go

ASSETS_FILES = $(shell set -x; $(ASSETS_FILES_FILTER))
NPM_FILES = $(shell set -x; $(NPM_FILES_FILTER))
GO_FILES = $(shell set -x; $(GO_FILES_FILTER)) $(GENERATED_GO_FILES)

.PHONY: npm-build go-build

build: go-build

# npm-build
NPM_BUILT_MARK = $(B)/.npm-built

$(NPM_BUILT_MARK): $(NPM_FILES) npm-deps Makefile
	@$(NPM) run build
	@mkdir -p $(@D)
	@touch $@

npm-build: npm-deps
	@$(NPM) run build
	@mkdir -p $(dirname $(BPN_BUILT_MARK))
	@touch $(NPM_BUILT_MARK)

.INTERMEDIATE: $(NPM_BUILT_MARK)

# go-build
$(ASSETS_GO_FILE): $(NPM_BUILT_MARK) $(FILE2GO) $(ASSETS_FILES)
	$(ASSETS_FILES_FILTER) | sort -uV | sed -e 's|^$(@D)/||' | (cd $(@D) && xargs -t $(notdir $(FILE2GO)) -p assets -o $(@F))

go-build: $(GO_FILES) go-deps
	$(GOGET) $(GOGET_FLAGS) ./...

.SECONDARY: $(GOBIN)/$(SERVER)
