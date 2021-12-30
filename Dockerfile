FROM golang:1.16 as base

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go build -o bin/main .

CMD ["/app/bin/main"]