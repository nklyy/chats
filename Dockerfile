FROM golang:1.17 AS base
ENV CGO_ENABLED=0

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ENV APP_ENV=production
RUN go build -o ./cmd/main ./cmd/main.go

CMD ["/app/cmd/main"]