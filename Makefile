APP_BIN=app/build/app

CURRENT_DIR=$(CURDIR)
GOLANGCI_LINT_VERSION ?= v2.12.2
GOLANGCI_LINT_VERSION_NUMBER := $(patsubst v%,%,$(GOLANGCI_LINT_VERSION))
TOOLS_DIR := .bin
GOLANGCI_LINT := $(TOOLS_DIR)/golangci-lint-$(GOLANGCI_LINT_VERSION)
GOLANGCI_LINT_RUNNER := ../$(GOLANGCI_LINT)

ifeq ($(OS),Windows_NT)
	GOLANGCI_LINT := $(TOOLS_DIR)/golangci-lint-$(GOLANGCI_LINT_VERSION).exe
	GOLANGCI_LINT_RUNNER := ../$(GOLANGCI_LINT)
endif

APP_DIR=${CURRENT_DIR}/app
INTERNAL_DIR=${APP_DIR}/internal
SQLC_DIR=${CURRENT_DIR}/sqlc
CMD_DIR=${APP_DIR}/cmd
MIGRATIONS_DIR=${APP_DIR}/migrations

DB_EXT=sql

create-migration:
	@read -p "Enter migration name: " name; \
	migrate create -ext $(DB_EXT) -dir $(MIGRATIONS_DIR) -seq $$name

run-app-watch:
	cd ${APP_DIR} && nodemon --watch . --ext go --signal SIGINT --exec 'go run ${CMD_DIR}/main.go'

run-app:
	cd ${APP_DIR} && go run ${CMD_DIR}/main.go

run:
	cd ${APP_DIR} && go run ${CMD_DIR}/main.go

swaggen:
	cd ${APP_DIR} && swag init -g internal/app/app.go -o internal/api/docs

sqlc-gen:
	cd ${SQLC_DIR}/ && sqlc generate

local-infra-up:
	docker compose -f docker-compose.local-infra.yml up -d

local-infra-down:
	docker compose -f docker-compose.local-infra.yml down

build:
	cd ${APP_DIR} && go build -o ../build/app.exe ./cmd/main.go

test:
	cd ${APP_DIR} && go test ./...

test-verbose:
	cd ${APP_DIR} && go test -v ./...

$(GOLANGCI_LINT):
ifeq ($(OS),Windows_NT)
	powershell -NoProfile -ExecutionPolicy Bypass -Command "$$ErrorActionPreference='Stop'; $$ProgressPreference='SilentlyContinue'; New-Item -ItemType Directory -Force -Path '$(TOOLS_DIR)' | Out-Null; Remove-Item -Recurse -Force '$(TOOLS_DIR)/golangci-lint-dist' -ErrorAction SilentlyContinue; Invoke-WebRequest -Uri 'https://github.com/golangci/golangci-lint/releases/download/$(GOLANGCI_LINT_VERSION)/golangci-lint-$(GOLANGCI_LINT_VERSION_NUMBER)-windows-amd64.zip' -OutFile '$(TOOLS_DIR)/golangci-lint.zip'; Expand-Archive -Force '$(TOOLS_DIR)/golangci-lint.zip' '$(TOOLS_DIR)/golangci-lint-dist'; Copy-Item '$(TOOLS_DIR)/golangci-lint-dist/golangci-lint-$(GOLANGCI_LINT_VERSION_NUMBER)-windows-amd64/golangci-lint.exe' '$(GOLANGCI_LINT)' -Force"
else
	mkdir -p "$(TOOLS_DIR)"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(TOOLS_DIR)" "$(GOLANGCI_LINT_VERSION)"
endif

lint: $(GOLANGCI_LINT)
	cd ${APP_DIR} && "$(GOLANGCI_LINT_RUNNER)" run ./...

migrate-up:
	cd ${APP_DIR} && migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/drivers_db?sslmode=disable" up

migrate-down:
	cd ${APP_DIR} && migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/drivers_db?sslmode=disable" down

clean:
	rm -rf ${APP_DIR}/build

help:
	@echo "Available commands:"
	@echo "  make run              - Run the application"
	@echo "  make build            - Build the application"
	@echo "  make test             - Run tests"
	@echo "  make test-verbose     - Run tests with verbose output"
	@echo "  make sqlc-gen         - Generate sqlc code"
	@echo "  make swaggen          - Generate swagger docs"
	@echo "  make migrate-up       - Run database migrations"
	@echo "  make migrate-down     - Rollback database migrations"
	@echo "  make local-infra-up   - Start local infrastructure (Postgres, MinIO)"
	@echo "  make local-infra-down - Stop local infrastructure"
	@echo "  make lint             - Run golangci-lint"
	@echo "  make clean            - Clean build artifacts"
