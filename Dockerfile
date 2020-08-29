FROM golang:1.15-alpine AS build

WORKDIR /gothrough

COPY . .

# Build server binary
RUN go build -o website .

# TODO install extra packages

# Build index.
RUN go run ./scripts/gen_index.go -dest=".interface_index" "$GOPATH" "/usr/local/go"

#

FROM alpine:3.12

RUN apk add ca-certificates

WORKDIR /gothrough

# Copy server binary from first stage.
COPY --from=build /gothrough/website .

# Copy index from first stage.
COPY --from=build /gothrough/.interface_index .

# Copy static files from project source.
COPY pages pages

ENV INDEX=.interface_index
CMD ./website