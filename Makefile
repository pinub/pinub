#!/usr/bin/env bash

GOCMD	:= go
GOBUILD	:= $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST	:= $(GOCMD) test
GORUN   := $(GOCMD) run

BINARY_NAME := server

HANDLER ?= handler
PACKAGE ?= ${HANDLER}

TEST_DB := "user=dan dbname=pinub_test sslmode=disable"

.PHONY: all build clean lambda run test todo

all: test build
build:	; @$(GOBUILD) -o ${BINARY_NAME} -v cmd/server/server.go
clean:	; @$(GOCLEAN) && rm -rf ${BINARY_NAME} ${HANDLER} ${PACKAGE}.zip .cache/
run:	; @$(GORUN) -v cmd/server/server.go
test:	; @DATABASE_URL=${TEST_DB} $(GOTEST) -cover ./...

lambda: clean
	GOOS=linux $(GOBUILD) -o ${HANDLER} -v cmd/lambda/lambda.go
	zip -r ${PACKAGE}.zip ${HANDLER} views

todo:
	@grep \
		--exclude-dir=./vendor \
		--exclude-dir=./client/node_modules \
		--exclude=./Makefile \
		--text \
		--color \
		-nRo ' TODO:.*' .

