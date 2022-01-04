FROM golang:1.17 as base

FROM base as dev

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

WORKDIR /opt/app/doubloonscape
CMD ["air"]

FROM base as built

WORKDIR /go/app/doubloonscape
COPY . .

ENV CGO_ENABLED=0

RUN go get -d -v ./...
RUN go build -o /tmp/doubloonscape ./*.go

FROM busybox

COPY --from=built /tmp/doubloonscape /usr/bin/doubloonscape
CMD ["doubloonscape"]