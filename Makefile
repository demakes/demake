DEMAKE_TEST_SETTINGS_SQLITE ?= `readlink -f settings/test/sqlite.json`
DEMAKE_TEST_SETTINGS_SQLITE_IN_MEMORY ?= `readlink -f settings/test/sqlite-in-memory.json`
DEMAKE_TEST_SETTINGS_POSTGRES ?= `readlink -f settings/test/postgres.json`

TESTARGS ?= ""
BENCHARGS ?= "."
all: build install

.PHONY: build install

build:
	@go build ./...

test-postgres:
	echo "Testing Postgres"
	@DEMAKE_SETTINGS=$(DEMAKE_TEST_SETTINGS_POSTGRES) go test ./... -count=1 -parallel=1 $(TESTARGS)

test-sqlite:
	echo "Testing SQLite"
	@DEMAKE_SETTINGS=$(DEMAKE_TEST_SETTINGS_SQLITE) go test ./... -count=1 -parallel=1 $(TESTARGS)

test-sqlite-in-memory:
	echo "Testing SQLite (in-memory)"
	@DEMAKE_SETTINGS=$(DEMAKE_TEST_SETTINGS_SQLITE_IN_MEMORY) go test ./... -count=1 -parallel=1 $(TESTARGS)

bench-sqlite:
	echo "Benchmarking SQLite"
	@DEMAKE_SETTINGS=$(DEMAKE_TEST_SETTINGS_SQLITE) go test -bench $(BENCHARGS) -count=1 -parallel=1 ./... $(TESTARGS)

bench-postgres:
	echo "Benchmarking Postgres"
	@DEMAKE_SETTINGS=$(DEMAKE_TEST_SETTINGS_POSTGRES) go test -bench $(BENCHARGS) -count=1 -parallel=1 ./... $(TESTARGS)

test: test-sqlite test-postgres

bench: bench-sqlite bench-postgres

install:
	@go install ./...

watch: install
	.scripts/watch_and_run.sh
