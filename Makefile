APP_NAME := snippetbox
APP_PATH := ./cmd/web
DSN ?= web:pass@/snippetbox?parseTime=true
ADDR ?= :4000

.PHONY: run build certs

run:
	@echo "Running $(APP_NAME)..."
	go run $(APP_PATH) -addr=$(ADDR) -dsn=$(DSN)

build:
	@echo "Building $(APP_NAME)..."
	go build -o ./bin/$(APP_NAME) $(APP_PATH)

certs:
	@echo "Generating self-signed certificate..."
	mkdir -p tls
	cd tls && go run $$(go env GOROOT)/src/crypto/tls/generate_cert.go --rsa-bits 2048 --host localhost