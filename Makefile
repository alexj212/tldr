VERSION=v0.1.56
export PROJ_PATH=github.com/alexj212/tldr




export DATE := $(shell date +%Y.%m.%d-%H%M)
export LATEST_COMMIT := $(shell git log --pretty=format:'%h' -n 1)
export BRANCH := $(shell git branch |grep -v "no branch"| grep \*|cut -d ' ' -f2)
export BUILT_ON_IP := $(shell [ $$(uname) = Linux ] && hostname -i || hostname )
export BIN_DIR=./bin

export BUILT_ON_OS=$(shell uname -a)
ifeq ($(BRANCH),)
BRANCH := master
endif

export COMMIT_CNT := $(shell git rev-list HEAD | wc -l | sed 's/ //g' )
export GIT_VERSION := ${BRANCH}-${COMMIT_CNT}
export COMPILE_LDFLAGS=-s -X "main.BuildDate=${DATE}" \
                          -X "main.GitCommit=${LATEST_COMMIT}" \
                          -X "main.GitVersion=${GIT_VERSION}" \
						  -X "main.Version=${VERSION}" \
                          -X "main.BuiltOnIP=${BUILT_ON_IP}" \
                          -X "main.BuiltOnOs=${BUILT_ON_OS}"



build_info: ## Display build info
	@echo ''
	@echo '---------------------------------------------------------'
	@echo 'BUILT_ON_IP       $(BUILT_ON_IP)'
	@echo 'BUILT_ON_OS       $(BUILT_ON_OS)'
	@echo 'DATE              $(DATE)'
	@echo 'LATEST_COMMIT     $(LATEST_COMMIT)'
	@echo 'BRANCH            $(BRANCH)'
	@echo 'COMMIT_CNT        $(COMMIT_CNT)'
	@echo 'COMPILE_LDFLAGS   $(COMPILE_LDFLAGS)'
	@echo 'PATH              $(PATH)'
	@echo '---------------------------------------------------------'
	@echo ''


####################################################################################################################
##
## help for each task - https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
##
####################################################################################################################
.PHONY: help

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help




####################################################################################################################
##
## Code vetting tools
##
####################################################################################################################

vet: ## run go vet on the project
	go vet .

tools: ## install dependent tools
	go install github.com/gordonklaus/ineffassign
	go install github.com/fzipp/gocyclo
	go install golang.org/x/lint/golint

lint: ## run golint on the project
	golint ./...


gocyclo: ## run gocyclo on the project
	@ gocyclo -avg -over 35 $(shell find . -name "*.go" |egrep -v "pb\.go|_test\.go")

ineffassign: ## run ineffassign on the project
	ineffassign ./...

check: lint vet ineffassign gocyclo ## run code checks on the project

doc: ## run godoc
	godoc -http=:6060

deps:## analyze project deps
	go list -f '{{ join .Deps  "\n"}}' . |grep "/" | grep -v "$(PROJ_PATH)"| grep "\." | sort |uniq

fmt: ## run fmt on the project
	## go fmt .
	gofmt -s -d -w -l .

####################################################################################################################
##
## Build of binaries
##
####################################################################################################################
all: build_app ## build and run tests

binaries: build_app ## build binaries in bin dir

create_dir:
	@mkdir -p $(BIN_DIR)

build_app: create_dir
		go build -o $(BIN_DIR)/$(BIN_NAME) -a -ldflags '$(COMPILE_LDFLAGS)' $(APP_PATH)



install:  build_app ## intall binary into ~/bin/ directory
	@mkdir -p ~/bin/
	cp ./bin/tldr ~/bin/
