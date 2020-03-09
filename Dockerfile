FROM golang:1.13-alpine

RUN mkdir /app
ADD . /app

WORKDIR /app
RUN go mod download

WORKDIR /app/cmd/app
RUN go build -o main *.go

EXPOSE 8080
CMD ["./main"]
