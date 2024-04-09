FROM golang:1.21.5-alpine

RUN mkdir /app

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o Back .

CMD ["/app/Back"]