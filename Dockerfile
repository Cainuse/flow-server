FROM golang:1.8

WORKDIR $GOPATH/src/github.com/flow_server

COPY . .

RUN go get -d -v ./...

EXPOSE 9090

CMD main

