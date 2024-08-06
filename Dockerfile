FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o main .

FROM debian:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8081

CMD ["./main"]
