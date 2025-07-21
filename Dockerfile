FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/marketplace
RUN go build -o /app/marketplace

FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/marketplace .

EXPOSE 8080
CMD ["./marketplace"]
