SHELL:=/bin/bash

.ONESHELL:

APP_VERSION ?= "0.1.0"
APP_EXECUTABLE ?= "termin-checker"

BUILD_COMMAND = CGO_ENABLED=0 go build -ldflags="-s -w"

all: build

build: FORCE setup
	CGO_ENABLED=0 go build -ldflags="-s -w" .

setup: FORCE
	go mod vendor
	go mod tidy

test: FORCE

dev: FORCE
	DEBUG=1 go run .

run: FORCE dev

.PHONY: FORCE
FORCE:
