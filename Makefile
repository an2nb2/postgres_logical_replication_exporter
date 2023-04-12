VERSION ?= 0.1.0

BINARY_NAME ?= exporter

CGO_ENABLED = 0
LDFLAGS     = -ldflags "-X main.version=${VERSION} -X main.commit=$$(git rev-parse --short HEAD 2>/dev/null || echo \"none\")"
MOD         = -mod=readonly

COMPOSECLI = $(shell command -v nerdctl || command -v docker || command -v podman-compose)

build-prerequisites:
	mkdir -p bin

### BUILD ###################################################################

build:
	cd cmd && CGO_ENABLED=$(CGO_ENABLED) go build $(LDFLAGS) -o ../bin/$(OUTPUT_DIR)$(BINARY_NAME) $(MOD)

build-linux_amd64: build-prerequisites
	$(MAKE) GOOS=linux GOARCH=amd64 OUTPUT_DIR=linux_amd64/ build
build-darwin_amd64: build-prerequisites
	$(MAKE) GOOS=darwin GOARCH=amd64 OUTPUT_DIR=darwin_amd64/ build
build-windows_amd64: build-prerequisites
	$(MAKE) GOOS=windows GOARCH=amd64 OUTPUT_DIR=windows_amd64/ build

build-all: build-linux_amd64 build-darwin_amd64 build-windows_amd64

### TEST ####################################################################

test-exporter: build
	TEST_PRIMARY_URI=postgres://primary:primary@localhost:9432/primary?sslmode=disable \
	TEST_STANDBY_URI=postgres://standby:standby@localhost:9442/standby?sslmode=disable \
	go test ./...

.PHONY: test
test: test-exporter

clean:
	rm -r bin/*
	go clean -testcache

### IMAGE ###################################################################

image:
	docker build -t postgres_logical_replication_exporter:${VERSION} .


### LOCAL DEBUG #############################################################

start-db:
	$(COMPOSECLI) up -d

stop-db:
	$(COMPOSECLI) down

seed-primary:
	PGUSER=primary PGPASSWORD=primary psql -h localhost -d primary -p 9432 < db/schema.sql
	PGUSER=primary PGPASSWORD=primary psql -h localhost -d primary -p 9432 < db/data.sql
	PGUSER=primary PGPASSWORD=primary psql -h localhost -d primary -p 9432 < db/publications.sql

seed-standby:
	PGUSER=standby PGPASSWORD=standby psql -h localhost -d standby -p 9442 < db/schema.sql
	PGUSER=standby PGPASSWORD=standby psql -h localhost -d standby -p 9442 < db/subscriptions.sql

seed-db: seed-primary seed-standby
