# syntax=docker/dockerfile:1.3-labs

###########################
####     Base image    ####
###########################
FROM golang:1.18-stretch AS base
MAINTAINER Matija Martinic <matija@volume.finance>
WORKDIR /app

###########################
#### Local development ####
###########################
FROM base AS local-dev
RUN cd /tmp && go install github.com/cespare/reflex@latest

COPY <<EOF /live-reload.sh
#!/usr/bin/env bash
exec reflex \\
  -d none \\
  -s \\
  -g "\*.go" \\
  -g "\*.proto" \\
  -g "\*.yaml" \\
  -- \\
  \$@
EOF

RUN chmod +x /live-reload.sh

# air is not set to entrypoint because I want to override that behaviour
# when using docker-compose run.
CMD ["/live-reload.sh", "go", "run", "./cmd/sparrow", "-c", "config.local-dev.yaml", "start"]

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
COPY --from=builder /app/config.example.yaml /config.example.yaml
