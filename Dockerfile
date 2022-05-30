# Builder
ARG GITHUB_PATH=gitlab.com/g6834/team17/analytics-service

FROM golang:1.18-alpine as builder

WORKDIR /app/${GITHUB_PATH}

RUN apk add --update make git curl
COPY Makefile Makefile
COPY . .
RUN make build

# Mail server
FROM alpine:latest as server
LABEL org.opencontainers.image.source https://${GITHUB_PATH}
WORKDIR /root/

COPY --from=builder /app/${GITHUB_PATH}/bin/analytics-service .

RUN chown root:root task-service

EXPOSE 3000

CMD ["./analytics-service"]
