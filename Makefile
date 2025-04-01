# Load environment variables from .env
include .env
export $(shell sed 's/=.*//' .env)

POSTGRES_CONTAINER=postgres-container

guard-%:
	@ test -n "${$*}" || (echo "FATAL: Environment variable $* is not set!"; exit 1)
# You'll have to create the application db and test db if you do not have createdb from Postgres client tools
db.migrate.up: guard-DB_NAME guard-GOOSE_DBSTRING
	@ createdb ${DB_NAME} 2>/dev/null || true
	@ env GOOSE_DBSTRING="${GOOSE_DBSTRING}" goose up
db.migrate.down: guard-DB_NAME guard-GOOSE_DBSTRING
	@ createdb ${DB_NAME} 2>/dev/null || true
	@ env GOOSE_DBSTRING="${GOOSE_DBSTRING}" goose down

# Create test database using Docker container command
db.test.prepare: guard-TEST_DATABASE_NAME guard-TEST_DATABASE_URL
	@ createdb ${TEST_DATABASE_NAME} 2>/dev/null || true
	@ docker exec $(POSTGRES_CONTAINER) psql -U ${DB_USER} -c "SELECT 1 FROM pg_database WHERE datname = '${TEST_DATABASE_NAME}'" | grep -q 1 || \
		docker exec $(POSTGRES_CONTAINER) createdb -U ${DB_USER} ${TEST_DATABASE_NAME}
	@ docker exec $(POSTGRES_CONTAINER) psql -U ${DB_USER} -d ${TEST_DATABASE_NAME} -c "DROP SCHEMA IF EXISTS public CASCADE; CREATE SCHEMA public;"
	@ env GOOSE_DBSTRING="${TEST_DATABASE_URL}" goose up

test: db.test.prepare
	go test -v ./tests/...

# Cleanup test database after tests
db.test.clean: 
	@ docker exec $(POSTGRES_CONTAINER) dropdb -U ${DB_USER} ${TEST_DATABASE_NAME} || true

.PHONY: test db.test.prepare db.test.clean