FROM golang:latest

COPY . /server
RUN go mod download

EXPOSE 80

