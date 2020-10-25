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
RUN go run ./scripts/index/main.go -dest=".index" "$GOPATH"

#

FROM alpine:3.12
RUN apk add ca-certificates

WORKDIR /gothrough
# Copy server binary from first stage.
COPY --from=build /gothrough/website .
# Copy index from first stage.
COPY --from=build /gothrough/.index .
# Copy static files from project source.
COPY pages pages

ENV INDEX=.index
CMD ./website