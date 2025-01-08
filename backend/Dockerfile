FROM golang:1.20-alpine
WORKDIR /app
COPY . .
RUN go mod init github.com/jw81/moody-weather/backend || true
RUN go mod tidy
RUN go build -o main .
EXPOSE 8080
CMD ["./main"]