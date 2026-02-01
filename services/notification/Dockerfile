FROM golang:1.25.2-alpine

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o notification-service ./services/notification/main.go

EXPOSE 8080

CMD ["./notification-service"]