FROM golang:1.25.2-alpine

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o shopping-service ./services/shopping/main.go

EXPOSE 8080

CMD ["./shopping-service"]