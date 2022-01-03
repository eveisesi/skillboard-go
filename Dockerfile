FROM golang:1.17.2 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o ./.build/skillboard-api ./cmd/skillz

FROM alpine:latest AS release
WORKDIR /app

RUN apk --no-cache add tzdata ca-certificates

COPY --from=builder /app/.build/skillboard-api .

LABEL maintainer="David Douglas <david@onetwentyseven.dev>"