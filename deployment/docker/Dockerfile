FROM golang:1.26-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o localbox ./cmd

FROM alpine:latest
WORKDIR /app

COPY --from=builder /build/localbox .
COPY --from=builder /build/website ./website

RUN mkdir -p uploads

EXPOSE 8080
CMD ["./localbox"]
