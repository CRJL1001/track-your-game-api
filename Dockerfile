FROM golang:1.26.6-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . . 
RUN CGO_ENABLED=0 GOOS=linux go build -o /track-your-game-api

FROM alpine:latest
WORKDIR /app
COPY --from=builder /track-your-game-api . 
RUN chmod +x ./track-your-game-api
EXPOSE 8080
CMD ["./track-your-game-api"]
