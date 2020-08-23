pkgs= $(shell go list ./... | grep -v /tests | grep -v /vendor/ | grep -v /common/)
curdir = $(shell pwd)
project_name = tauria-test
build_dir ?= bin/

build:
	@echo "building..."
	@go build -o $(curdir)/$(build_dir)$(project_name) $(curdir)/cmd/$(project_name)/app.go

run:
	@echo "running..."
	@go run $(curdir)/cmd/$(project_name)/app.go
