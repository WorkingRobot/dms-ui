# syntax=docker/dockerfile:1

# --- Stage 1: build the SvelteKit frontend into internal/web/dist ---
FROM node:22-alpine AS web
WORKDIR /src/web
COPY web/package.json web/package-lock.json* ./
RUN npm ci || npm install
COPY web/ ./
RUN npm run build

# --- Stage 2: build the Go binary with the frontend embedded ---
FROM golang:1.23-alpine AS go
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
# Bring in the freshly built static assets so //go:embed picks them up.
COPY --from=web /src/internal/web/dist ./internal/web/dist
# Quality gate runs inside the build (mirrors the org's "everything in the
# Dockerfile" CI convention) so a broken build/test fails the image.
RUN go vet ./... && go test ./...
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o /dms-ui ./cmd/dms-ui

# --- Stage 3: tiny runtime (needs the docker CLI to exec into the mail container) ---
FROM alpine:3.20
RUN apk add --no-cache docker-cli ca-certificates tzdata
COPY --from=go /dms-ui /usr/local/bin/dms-ui
EXPOSE 8080
ENV LISTEN_ADDR=:8080
USER root
ENTRYPOINT ["/usr/local/bin/dms-ui"]
