FROM golang:1.22-alpine
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go mod verify

COPY . .
RUN go build -o service ./src/cmd/service
