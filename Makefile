all: build install

.PHONY: build install

build:
	@go build ./...

install:
	@go install ./...

watch: install
	.scripts/watch_and_run.sh
