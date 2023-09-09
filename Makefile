-include .env
export $(shell [ -f ".env" ] && sed 's/=.*//' .env)

export DATE := $(shell date +%Y.%m.%d-%H%M)

export LATEST_COMMIT := $(shell git rev-parse HEAD 2> /dev/null)
export COMMIT_CNT := $(shell git rev-list --all 2> /dev/null | wc -l | sed 's/ //g' )
export BRANCH := $(shell git branch  2> /dev/null |grep -v "no branch"| grep \*|cut -d ' ' -f2)
export GIT_REPO := $(shell git config --get remote.origin.url  2> /dev/null)
export COMMIT_DATE := $(shell git log -1 --format=%cd  2> /dev/null)


ifeq ($(BRANCH),)
BRANCH := master
endif

export VERSION_FILE   := version.txt
export TAG     := $(shell [ -f "$(VERSION_FILE)" ] && cat "$(VERSION_FILE)" || echo '0.2.0')
export VERMAJMIN      := $(subst ., ,$(TAG))
export VERSION        := $(word 1,$(VERMAJMIN))
export MAJOR          := $(word 2,$(VERMAJMIN))
export MINOR          := $(word 3,$(VERMAJMIN))
export NEW_MINOR      := $(shell expr "$(MINOR)" + 1)
export NEW_TAG := $(VERSION).$(MAJOR).$(NEW_MINOR)

export COMMIT_CNT := $(shell git rev-list HEAD | wc -l | sed 's/ //g' )
export BUILD_NUMBER := $(BRANCH)-$(COMMIT_CNT)
export COMPILE_LDFLAGS=-s -X "main.BuildDate=${DATE}" \
                          -X "main.LatestCommit=${LATEST_COMMIT}" \
						  -X "main.Version=${NEW_TAG}"\
						  -X "main.GitRepo=${GIT_REPO}" \
                          -X "main.GitBranch=${BRANCH}"



build_info: ## Build info
	@echo ''
	@echo '---------------------------------------------------------'
	@echo 'DATE              $(DATE)'
	@echo 'LATEST_COMMIT     $(LATEST_COMMIT)'
	@echo 'BRANCH            $(BRANCH)'
	@echo 'COMPILE_LDFLAGS   $(COMPILE_LDFLAGS)'
	@echo 'TAG               $(TAG)'
	@echo 'NEW_TAG           $(NEW_TAG)'	
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
	rm -rf .history
	gofmt -s -d -w -l .

####################################################################################################################
##
## Build of binaries
##
####################################################################################################################
all: build_app ## build and run tests

binaries: build_app ## build binaries in bin dir

create_dir:
	@mkdir -p ./bin

build_app: create_dir
		go build -o ./bin/$(BIN_NAME) -a -ldflags '$(COMPILE_LDFLAGS)' $(APP_PATH)


install:  build_app ## intall binary into ~/bin/ directory
	@mkdir -p ~/bin/
	cp ./bin/tldr ~/bin/


upgrade: #upgrade go libs to the latest version
	go get -u ./...
	go mod tidy



publish: #publish to github
	@echo "\n\n\n\n\n\nRunning git add\n"
	echo "$(NEW_TAG)" > "$(VERSION_FILE)"
	git add -A
	@echo "\n\n\n\n\n\nRunning git commit v$(NEW_TAG)\n"
	git commit -m "latest version: v$(NEW_TAG)"

	@echo "\n\n\n\n\n\nRunning git tag\n"
	git tag  "v$(NEW_TAG)"

	@echo "\n\n\n\n\n\nRunning git push\n"
	git push -f origin "v$(NEW_TAG)"

	git push -f
