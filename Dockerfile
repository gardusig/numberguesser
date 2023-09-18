FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go test ./test -v
