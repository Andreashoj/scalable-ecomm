FROM golang:1.25.2-alpine

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o product-catalog-service ./services/product-catalog/main.go

EXPOSE 8080

CMD ["./product-catalog-service"]