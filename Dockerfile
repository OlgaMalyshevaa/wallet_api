FROM golang:1.23

WORKDIR /app

RUN apt-get update && apt-get install -y postgresql-client

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o wallet ./cmd/wallet_app

CMD ["/app/wallet"]

EXPOSE 8080