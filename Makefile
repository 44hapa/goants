LOCAL_BIN=$(CURDIR)/bin
PROJECT_NAME=ants

#include $(GOROOT)/src/Make.inc

export GO111MODULE=off
GOENV:=GO111MODULE=off

TARG=MyBot
GOFILES=\
	ants.go\
	map.go\
	main.go\
	debugging.go\
	MyBot.go\

#include $(GOROOT)/src/Make.cmd

.PHONY: build
build:
	$(GOENV) CGO_ENABLED=0 go build -v -o /Users/pavel.sukhorukov/Downloads/tools/sample_bots/go/ants1 ./cmd
#	$(GOENV) CGO_ENABLED=0 go build -v -o $(LOCAL_BIN)/$(PROJECT_NAME) ./cmd

.PHONY: run
run:
	$(GOENV) go run cmd/main.go
