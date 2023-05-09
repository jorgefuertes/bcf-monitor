# BcfMonitor Makefile

SHELL = /usr/bin/env bash
APIHOST = api.blockchainfue.com
OUT = tmp/bcf-monitor
CONF = conf/dev.yaml
BUILD_NUM_FILE ?= .build
FLAGS = -s -w

clean:
	@echo Cleaning...
	@rm -Rf tmp
	@go clean -testcache

tunnels:
	scripts/tunnel-up.sh

test:
	pushd pkg/mail;	go test -v .; popd
	pushd pkg/monitor/mongo; go test -v .; popd
	pushd pkg/monitor/redis; go test -v .; popd

build:
	@echo Building $(OUT)...
	go build -ldflags "$(FLAGS)" -o $(OUT)

run:
	go run .
