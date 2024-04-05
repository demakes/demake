all: build install

.PHONY: build install

build:
	@go build ./...

test:
	@KLARO_SETTINGS=`readlink -f settings/test/sqlite.json` go test ./... -count=1 -parallel=1

install:
	@go install ./...

watch: install
	.scripts/watch_and_run.sh
