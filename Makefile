BINARY := yatter-backend-go

# MAKEFILE_LIST: makeがパースするファイルリスト
# lastword: この時点ではincludeしてないので'Makefile'
MAKEFILE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

PATH := $(PATH):${MAKEFILE_DIR}bin
SHELL := env PATH="$(PATH)" /bin/bash
# for go
export CGO_ENABLED = 0
GOARCH = amd64

# HEADが指すオブジェクトのハッシュ値
COMMIT=$(shell git rev-parse HEAD)
# HEADが指すオブジェクトのブランチ名
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
GIT_URL=local-git://

# go tool link に渡すフラグ, リンク時に指定パッケージ内の指定変数書き換え
LDFLAGS := -ldflags "-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"

.PHONY: build
build: build-linux

.PHONY: build-default
build-default:
	go build ${LDFLAGS} -o build/${BINARY}

.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o build/${BINARY}-linux-${GOARCH} .

.PHONY: prepare
prepare: mod

.PHONY: mod
mod:
	go mod download

# daoパッケージのテストはコンテナ内で実行するので除く
.PHONY: test
test:
	go test $(filter-out %/dao, $(shell go list ${MAKEFILE_DIR}/...))
	@echo exclude dao package

.PHONY: alltest
alltest:
	go test $(shell go list ${MAKEFILE_DIR}/...)

.PHONY: wtest
wtest:
	docker-compose exec web make alltest

.PHONY: lint
lint:
	if ! [ -x $(GOPATH)/bin/golangci-lint ]; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.38.0 ; \
	fi
	golangci-lint run --concurrency 2

.PHONY: vet
vet:
	go vet ./...

.PHONY:	clean
clean:
	git clean -f -X app bin build
