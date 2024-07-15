FROM golang:1.22.3

WORKDIR /market-history-storage

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN GOOS=linux go build -a -o server cmd/main.go

EXPOSE ${HTTP_SERVER_PORT}

CMD ["./server"]