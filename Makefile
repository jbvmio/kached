.DEFAULT_GOAL := test

export GO111MODULE=on

all: test

vars:
YTIME="1546300800"
BT=$(shell date +%s)
GCT=$(shell git rev-list -1 HEAD --timestamp | awk '{print $$1}')
GC=$(shell git rev-list -1 HEAD --abbrev-commit)
REV=$(shell echo $(GCT) - $(YTIME) | bc)

test: vars
	@echo "build time       :" $(shell echo $(BT))
	@echo "git commit time  :" $(shell echo $(GCT))
	@echo "git commit hash  :" $(shell echo $(GC))
	@echo "revision         :" $(shell echo $(REV))
	mkdir -p ./dbDir
	touch ./dbDir/tmpfile
	rm ./dbDir/*
	go test .
