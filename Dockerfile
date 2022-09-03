# syntax=docker/dockerfile:1

FROM golang:1.19

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /usr/local/bin/app
RUN go get -d -v ./...

RUN go run /usr/src/app
EXPOSE 8000
CMD ["app"]
