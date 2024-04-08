KLARO_TEST_SETTINGS_SQLITE ?= `readlink -f settings/test/sqlite.json`
KLARO_TEST_SETTINGS_SQLITE_IN_MEMORY ?= `readlink -f settings/test/sqlite-in-memory.json`
KLARO_TEST_SETTINGS_POSTGRES ?= `readlink -f settings/test/postgres.json`

all: build install

.PHONY: build install

build:
	@go build ./...

test-postgres:
	echo "Testing Postgres"
	@KLARO_SETTINGS=$(KLARO_TEST_SETTINGS_POSTGRES) go test ./... -count=1 -parallel=1

test-sqlite:
	echo "Testing SQLite"
	@KLARO_SETTINGS=$(KLARO_TEST_SETTINGS_SQLITE) go test ./... -count=1 -parallel=1

test-sqlite-in-memory:
	echo "Testing SQLite (in-memory)"
	@KLARO_SETTINGS=$(KLARO_TEST_SETTINGS_SQLITE_IN_MEMORY) go test ./... -count=1 -parallel=1

bench-sqlite:
	echo "Benchmarking SQLite"
	@KLARO_SETTINGS=$(KLARO_TEST_SETTINGS_SQLITE) go test ./... -bench=. -run=NONE -count=1 -parallel=1

bench-postgres:
	echo "Benchmarking Postgres"
	@KLARO_SETTINGS=$(KLARO_TEST_SETTINGS_POSTGRES) go test ./... -bench=. -run=NONE -count=1 -parallel=1

test: test-sqlite test-postgres

bench: bench-sqlite bench-postgres

install:
	@go install ./...

watch: install
	.scripts/watch_and_run.sh
