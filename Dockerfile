FROM golang:1.18-alpine as builder

RUN mkdir -p /app
COPY . /app
WORKDIR /app

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o application ./cmd/api

FROM scratch

COPY --from=builder /app /app

ENTRYPOINT ["/app/application"]
