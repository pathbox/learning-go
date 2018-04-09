FROM golang:1.10.0

RUN adduser --disabled-password --gecos '' api
USER api

WORKDIR /go/src/app
COPY . .

RUN go install -v ./...

CMD [ "app" ]