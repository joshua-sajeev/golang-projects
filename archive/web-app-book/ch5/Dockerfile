# syntax=docker/dockerfile:1
FROM golang:latest
WORKDIR /app
ADD . /app/
RUN go build -o ./out/main .
EXPOSE 8080
ENTRYPOINT ["./out/main"]
