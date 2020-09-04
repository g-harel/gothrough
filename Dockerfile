FROM golang:1.15-alpine AS build

# Required to fetch go modules.
RUN apk add git

# Add non-standard-library packages.
COPY packages.txt .
RUN xargs -a packages.txt go get -u

WORKDIR /gothrough
COPY . .

# Build server binary
RUN go build -o website .

# Build index.
RUN go run ./scripts/index/main.go -dest=".interface_index" "$GOPATH" "/usr/local/go"

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