FROM golang:1.18

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . /app
WORKDIR /app/cmd/apiserver

RUN go build -o /apiservice .

EXPOSE 8080

CMD [ "/apiservice" ]