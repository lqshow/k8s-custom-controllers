SHELL := /bin/bash

SOURCES = $(shell find . -name '*.go') $(GENERATED_PROTOS)

all: build


.PHONY: clean-bin
clean-bin:
	rm -rf bin .go


.PHONY: build
build:
	$(MAKE) $(SOURCES)
	./script/build.sh


.PHONY: clean
clean:
	$(MAKE) clean-bin