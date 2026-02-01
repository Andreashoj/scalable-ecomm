FROM golang:1.25.2-alpine

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o user-service ./services/user/main.go

EXPOSE 8080

CMD ["./user-service"]