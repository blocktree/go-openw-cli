.PHONY: all clean
.PHONY: openw-cli
.PHONY: deps

# Check for required command tools to build or stop immediately
EXECUTABLES = git go find pwd
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH)))

GO ?= latest

# openw-cli
OPENWCLIVERSION = $(shell git describe --tags `git rev-list --tags --max-count=1`)
OPENWCLIBINARY = openw-cli
OPENWCLIMAIN = main.go

BUILDDIR = build
GITREV = $(shell git rev-parse --short HEAD)
BUILDTIME = $(shell date +'%Y-%m-%d_%T')

OPENWCLILDFLAGS="-X github.com/blocktree/go-openw-cli/openwcli.Version=${OPENWCLIVERSION} \
	-X github.com/blocktree/go-openw-cli/openwcli.GitRev=${GITREV} \
	-X github.com/blocktree/go-openw-cli/openwcli.BuildTime=${BUILDTIME} \
	-X github.com/blocktree/go-openw-cli/openwcli.FixAppID=${APPID} \
	-X github.com/blocktree/go-openw-cli/openwcli.FixAppKey=${APPKEY}"

# OS platfom
# options: windows-6.0/*,darwin-10.10/amd64,linux/amd64,linux/386,linux/arm64,linux/mips64, linux/mips64le
TARGETS="darwin-10.10/amd64,linux/amd64,windows-6.0/*"

deps:
	go get -u github.com/gythialy/xgo

build:
	GO111MODULE=on go build -ldflags $(OPENWCLILDFLAGS) -i -o $(shell pwd)/$(BUILDDIR)/$(OPENWCLIBINARY) $(shell pwd)/$(OPENWCLIMAIN)
	@echo "Build $(OPENWCLIBINARY) done."

all: openw-cli

clean:
	rm -rf $(shell pwd)/$(BUILDDIR)/

openw-cli:
	xgo --dest=$(BUILDDIR) --ldflags=$(OPENWCLILDFLAGS) --out=$(OPENWCLIBINARY)-$(OPENWCLIVERSION)-$(GITREV) --targets=$(TARGETS) \
	--pkg=$(OPENWCLIMAIN) .
