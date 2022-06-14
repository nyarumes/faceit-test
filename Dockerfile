FROM golang:1.18.3-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN make build

EXPOSE 80 8082 8083

CMD make run