FROM golang:1.15-alpine AS build

# Required to fetch go modules.
RUN apk add git

WORKDIR /gothrough

COPY . .

# Build server binary
RUN go build -o website .

# TODO install extra packages

# Build index.
RUN go run ./scripts/gen_index.go -dest=".interface_index" "$GOPATH" "$GOROOT"

#

FROM alpine:3.12

RUN apk add ca-certificates

WORKDIR /gothrough

# Copy server binary from first stage.
COPY --from=build /gothrough/website .

# Copy index from first stage.
COPY --from=build /gothrough/.interface_index .

# Copy static files from project source.
COPY templates templates

ENV INDEX=.interface_index
CMD ./website