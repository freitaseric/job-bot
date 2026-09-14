# syntax=docker/dockerfile:1

FROM golang:1.27-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/job-bot \
    ./cmd/job-bot


FROM alpine:3.22

RUN apk add --no-cache ca-certificates \
    && addgroup -S jobbot \
    && adduser -S -G jobbot jobbot

WORKDIR /app

COPY --from=builder /out/job-bot /app/job-bot

USER jobbot

ENTRYPOINT ["/app/job-bot"]