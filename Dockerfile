FROM golang:1.18-alpine as builder

RUN mkdir /app
ENV GO111MODULE=on CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY . /app

RUN go build -o apiservice /app/cmd/apiserver/main.go

RUN chmod +x /app/apiservice

#build tiny immage
FROM alpine:latest

RUN mkdir /app
COPY --from=builder /app/apiservice /app

#copy static components
WORKDIR /app
COPY .env /app/
RUN mkdir /migrations
COPY ./internal/store/migrations /app/migrations

CMD [ "/app/apiservice" ]