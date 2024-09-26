FROM golang:1.22

WORKDIR /app

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./*.go ./

COPY ./docs ./docs

RUN go build -o /hr-api

COPY ./frontend ./frontend

EXPOSE 5000

CMD ["/hr-api"]