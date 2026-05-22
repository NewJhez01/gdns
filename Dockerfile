# syntax=docker/dockerfile:1

FROM golang:1.26.3

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gdns ./cmd/main.go

EXPOSE 8080

CMD ["/docker-gdns"]
