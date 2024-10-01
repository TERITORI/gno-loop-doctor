# syntax=docker/dockerfile:1

FROM ghcr.io/gnolang/tx-indexer:0.5.1 AS tx-indexer

# Build the application from source
FROM golang:1.22 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY **/*.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./indexer-runner

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/indexer-runner /indexer-runner

COPY --from=tx-indexer /tx-indexer /gno-tx-indexer

EXPOSE 8546

ENTRYPOINT ["/indexer-runner"]