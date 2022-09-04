# syntax=docker/dockerfile:1

FROM golang:1.18

EXPOSE 8000

WORKDIR /go/src/app
# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

RUN go build -o /event
CMD ["/event", "--host", "0.0.0.0"]
