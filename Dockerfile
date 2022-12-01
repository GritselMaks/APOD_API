FROM golang:1.18 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod tidy

COPY . /app
WORKDIR /app/cmd/apiserver

RUN go build -o /apiservice .

CMD [ "/apiservice" ]