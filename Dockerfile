##########################################
#               Base image               #
##########################################
FROM golang:1.24-alpine AS base

ARG SRC_PATH

WORKDIR /go/src/github.com/${SRC_PATH}

COPY go.mod go.sum ./
RUN go mod download
COPY . .

##########################################
#               Linter step              #
##########################################
FROM base AS linter

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.0.2 && \
    make lint

##########################################
#               Build Step               #
##########################################
FROM base AS builder

RUN apk add --no-cache upx=4.2.4-r0

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o builds/app-linux-amd64 ./cmd && \
    upx --best --lzma builds/app-linux-amd64

##########################################
#                Artifact                #
##########################################
FROM alpine:3.20

ARG SRC_PATH

RUN apk --no-cache add ca-certificates

WORKDIR /app/
COPY --from=builder /go/src/github.com/${SRC_PATH}/builds/app-linux-amd64 ./app

EXPOSE 8080

CMD ["./app"]