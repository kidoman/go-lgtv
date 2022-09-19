SHELL := /bin/bash

.EXPORT_ALL_VARIABLES:
SRC_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BIN_DIR := $(SRC_DIR)/bin
GOPROXY ?= http://proxy.golang.org
GO111MODULE := on

$(@info $(shell mkdir -p $(BIN_DIR)))

bin/lgtv: $(SRC_DIR)/cmd/main.go
	go build -o $(BIN_DIR)/lgtv $(SRC_DIR)/cmd/main.go

.PHONY: install
install: bin/lgtv
	sudo cp bin/lgtv /usr/local/bin/lgtv
