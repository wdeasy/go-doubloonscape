FROM golang:1.17 as base

FROM base as dev

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

WORKDIR /opt/app/doubloonscape
CMD ["air"]

FROM base as built

WORKDIR /go/app/app/doubloonscape
COPY . .

ENV CGO_ENABLED=0

RUN go get -d -v ./...
RUN go build -o /tmp/app/doubloonscape ./*.go

FROM busybox

COPY --from=built /tmp/app/doubloonscape /usr/bin/app/doubloonscape
CMD ["app/doubloonscape"]