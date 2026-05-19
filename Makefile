.PHONY: web build run test docker

# Build the SvelteKit frontend into the Go embed directory.
web:
	cd web && bun install --frozen-lockfile
	cd web && bun run build

# Build the single binary (frontend must be built first).
build: web
	go build -trimpath -ldflags "-s -w" -o dms-ui ./cmd/dms-ui

test:
	go vet ./...
	go test ./...

# Run locally against a docker-mailserver container on the local Docker
# daemon. Override MAIL_CONTAINER / DOCKER_HOST as needed, e.g.
#   make run DOCKER_HOST=ssh://my-host MAIL_CONTAINER=mailserver
MAIL_CONTAINER ?= mailserver
run:
	MAIL_CONTAINER=$(MAIL_CONTAINER) LISTEN_ADDR=127.0.0.1:8099 ./dms-ui

docker:
	docker build -t dms-ui:dev .
