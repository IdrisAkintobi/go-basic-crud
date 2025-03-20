# Load environment variables from .env
include .env
export $(shell sed 's/=.*//' .env)

guard-%:
	@ test -n "${$*}" || (echo "FATAL: Environment variable $* is not set!"; exit 1)

db.migrate.up: guard-DB_NAME guard-GOOSE_DBSTRING
	@ createdb ${DB_NAME} 2>/dev/null || true
	@ env GOOSE_DBSTRING="${GOOSE_DBSTRING}" goose up

db.migrate.down: guard-DB_NAME guard-GOOSE_DBSTRING
	@ createdb ${DB_NAME} 2>/dev/null || true
	@ env GOOSE_DBSTRING="${GOOSE_DBSTRING}" goose down

db.test.prepare: guard-TEST_DATABASE_NAME guard-TEST_DATABASE_URL
	@ createdb ${TEST_DATABASE_NAME} 2>/dev/null || true
	@ env GOOSE_DBSTRING="${TEST_DATABASE_URL}" goose up

test: db.test.prepare
	go test -v ./tests/...
