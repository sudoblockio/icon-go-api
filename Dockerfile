FROM golang:1.16-buster AS builder

# GO ENV VARS
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# COPY SRC
WORKDIR /build
COPY ./src .

RUN go mod tidy

# CREATE SWAGGER DOCS
RUN go get github.com/swaggo/swag/cmd/swag
RUN go get github.com/arsmn/fiber-swagger/v2@v2.20.0
RUN go get github.com/alecthomas/template
RUN go get github.com/riferrei/srclient@v0.3.0
WORKDIR /build/
RUN swag init --parseDependency -g api/api.go

# BUILD
WORKDIR /build
RUN go build -o main ./

FROM ubuntu as prod

# For SSL certs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /build/main /
CMD ["/main"]
