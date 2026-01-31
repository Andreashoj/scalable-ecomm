FROM golang:1.25.2-alpine

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o order-service ./services/order-service/main.go

EXPOSE 8080

CMD ["./order-service"]