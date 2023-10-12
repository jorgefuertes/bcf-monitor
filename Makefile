# BcfMonitor Makefile

include .env
include .secrets

SHELL = /usr/bin/env bash
build_num_create := $(shell ! test -f .build && echo 0 > .build)
build_num = $(shell echo $$(($$(cat .build) + 1)) > .build && cat .build)
FLAGS = -s -w
FLAGS += -X main.buildNumber=$(call build_num)
FLAGS += -X main.buildDate=$(shell date +'%Y%m%d')
FLAGS += -X main.version=$(VERSION)

clean:
	@echo Cleaning...
	@rm -Rf $(TMP_DIR)
	@mkdir $(TMP_DIR)
	@go clean -testcache

tunnels:
	scripts/tunnel-up.sh

test:
	go test ./... -run "^Test*"

every-release: DEST = $(TMP_DIR)/$(EXE_NAME)_$(OS)_$(ARCH)
every-release:
	@echo Building $(DEST)
	@mkdir -p $(DEST)
	@mkdir -p $(TMP_DIR)/release
	@cp $(CONF_DIST) $(DEST)/.
	@go mod tidy
	@GOOS=$(OS) CGO_ENABLED=0 GOARCH=$(ARCH) go build -ldflags "$(FLAGS)" -o $(DEST)/$(EXE_NAME)
	@echo Compressing to $(TMP_DIR)/release/$(EXE_NAME)_$(OS)_$(ARCH)_$(VERSION).tar.gz
	@tar -C $(TMP_DIR) -czf $(TMP_DIR)/release/$(EXE_NAME)_$(OS)_$(ARCH)_$(VERSION).tar.gz $(EXE_NAME)_$(OS)_$(ARCH)

build-linux-amd64:
	@OUT=tmp/bcf-monitor_linux_amd64 OS=linux ARCH=amd64 make every-release
build-linux-arm:
	@OUT=tmp/bcf-monitor_linux_arm OS=linux ARCH=arm make every-release
build-linux-arm64:
	@OUT=tmp/bcf-monitor_linux_arm64 OS=linux ARCH=arm64 make every-release
build-osx-amd64:
	@OUT=tmp/bcf-monitor_osx_amd64 OS=darwin ARCH=amd64 make every-release
build-osx-arm64:
	@OUT=tmp/bcf-monitor_osx_arm64 OS=darwin ARCH=arm64 make every-release

publish: test build-linux-amd64
	@echo "Service stop..."
	@ssh root@$(SERVER) "service bcf-monitor stop; true"
	@echo "Uploading..."
	@ssh root@$(SERVER) mkdir -p /root/bcf-monitor
	@scp tmp/bcf-monitor_linux_amd64/$(EXE_NAME) root@$(SERVER):/usr/local/bin/.
	@scp $(CONF_PROD) root@$(SERVER):/etc/bcfmonitor-conf.yaml
	@scp $(SERVICE) root@$(SERVER):/etc/systemd/system/.
	@echo "Starting service..."
	@ssh root@$(SERVER) "systemctl daemon-reload; service bcf-monitor restart"
	@echo "Ready"

release: build-linux-amd64 build-linux-arm build-linux-arm64 build-osx-amd64 build-osx-arm64

run:
	go run . -c $(CONF_DEV)
