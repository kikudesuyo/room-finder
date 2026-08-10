FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /room-finder-api .

FROM alpine:3.21

WORKDIR /app
COPY --from=builder /room-finder-api ./room-finder-api
EXPOSE 8081
ENTRYPOINT ["./room-finder-api"]
