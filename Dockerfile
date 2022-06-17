FROM golang:latest

COPY . /go/src/app

WORKDIR /go/src/app

RUN go build -o apiserver cmd/apiserver/main.go

CMD ["./apiserver"]