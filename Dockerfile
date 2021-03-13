FROM golang:1.16.0-buster

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go build

ENTRYPOINT ["/app/pasuri"]

EXPOSE 8080/tcp
