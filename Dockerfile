# syntax=docker/dockerfile:1.3-labs

###########################
####     Base image    ####
###########################
FROM golang:1.20-bullseye AS base
WORKDIR /app

###########################
#### Local development ####
###########################
FROM base AS local-dev
RUN cd /tmp && go install github.com/cespare/reflex@latest

COPY <<EOF /hack-start-because-cosmos-always-wants-to-read-pass-from-stdin.sh
#!/usr/bin/env bash
go run ./cmd/pigeon -c config.local-dev.yaml start < /dev/null
EOF
RUN chmod +x /hack-start-because-cosmos-always-wants-to-read-pass-from-stdin.sh

CMD ["/app/scripts/live-reload.sh", "/hack-start-because-cosmos-always-wants-to-read-pass-from-stdin.sh"]

###########################
####     Builder       ####
###########################
FROM base AS builder
COPY . /app
RUN \
	--mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	cd /app && go build -o /sparrow ./cmd/sparrow

###########################
####  Local testnet    ####
###########################
FROM ubuntu AS local-testnet
ENTRYPOINT ["/sparrow"]
COPY --from=builder /sparrow /sparrow


###########################
####     Release       ####
###########################
FROM base AS release
RUN go install github.com/goreleaser/goreleaser@latest
COPY . /app

CMD ["goreleaser", "release", "--rm-dist"]
